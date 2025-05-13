package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

// Error definitions
var (
	ErrMissingMorpheEntityField = func(entityName, fieldName string) error {
		return fmt.Errorf("missing entity field %s in entity %s", fieldName, entityName)
	}
)

// AllMorpheEntitiesToPSQLViews compiles all Morphe entities to PostgreSQL views
func AllMorpheEntitiesToPSQLViews(config MorpheCompileConfig, r *registry.Registry) (map[string]*psqldef.View, error) {
	allViews := map[string]*psqldef.View{}
	for entityName, entity := range r.GetAllEntities() {
		view, viewErr := MorpheEntityToPSQLView(config, r, entity)
		if viewErr != nil {
			return nil, viewErr
		}
		allViews[entityName] = view
	}
	return allViews, nil
}

// MorpheEntityToPSQLView compiles a single Morphe entity to a PostgreSQL view
func MorpheEntityToPSQLView(config MorpheCompileConfig, r *registry.Registry, entity yaml.Entity) (*psqldef.View, error) {
	if r == nil {
		return nil, triggerCompileMorpheEntityFailure(config.EntityHooks, config.MorpheConfig, entity, ErrNoRegistry)
	}

	morpheConfig, entity, compileStartErr := triggerCompileMorpheEntityStart(config.EntityHooks, config.MorpheConfig, entity)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheEntityFailure(config.EntityHooks, config.MorpheConfig, entity, compileStartErr)
	}
	config.MorpheConfig = morpheConfig

	view, viewErr := morpheEntityToPSQLView(config.MorpheConfig, r, entity)
	if viewErr != nil {
		return nil, triggerCompileMorpheEntityFailure(config.EntityHooks, config.MorpheConfig, entity, viewErr)
	}

	view, compileSuccessErr := triggerCompileMorpheEntitySuccess(config.EntityHooks, view)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheEntityFailure(config.EntityHooks, config.MorpheConfig, entity, compileSuccessErr)
	}

	return view, nil
}

func morpheEntityToPSQLView(config cfg.MorpheConfig, r *registry.Registry, entity yaml.Entity) (*psqldef.View, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}

	validateEntityErr := entity.Validate(r.GetAllModels(), r.GetAllEnums())
	if validateEntityErr != nil {
		return nil, validateEntityErr
	}

	viewName := strcase.ToSnakeCaseLower(entity.Name)
	if config.MorpheEntitiesConfig.ViewNameSuffix != "" {
		viewName += config.MorpheEntitiesConfig.ViewNameSuffix
	}

	// TODO: Extract all "root" models from the entity fields, and use the first one as the base table name
	tableName := Pluralize(strcase.ToSnakeCaseLower(entity.Name))

	view := &psqldef.View{
		Schema:     config.MorpheEntitiesConfig.Schema,
		Name:       viewName,
		FromSchema: config.MorpheEntitiesConfig.Schema,
		FromTable:  tableName,
		Columns:    []psqldef.ViewColumn{},
		Joins:      []psqldef.JoinClause{},
	}

	// Process entity fields to set up columns and joins
	joinTables := make(map[string]bool)
	// rootTableRelationships := make(map[string]string)
	joinTableRelationships := make(map[string]string)

	fieldNames := core.MapKeysSorted(entity.Fields)
	for _, fieldName := range fieldNames {
		field := entity.Fields[fieldName]
		// Convert field name to snake case for column name
		columnName := strcase.ToSnakeCaseLower(fieldName)

		// Parse the field type (e.g., "User.UUID" or "User.Child.AutoIncrement")
		fieldParts := strings.Split(string(field.Type), ".")
		if len(fieldParts) < 2 {
			return nil, fmt.Errorf("invalid field type format: %s", field.Type)
		}

		// Determine source reference for this column
		var sourceRef string
		if len(fieldParts) == 2 {
			// Direct field from main model
			sourceRef = fmt.Sprintf("%s.%s", tableName, strcase.ToSnakeCaseLower(fieldParts[1]))
		} else if len(fieldParts) == 3 {
			// Field from related model
			relatedModelName := fieldParts[1]
			relatedFieldName := fieldParts[2]
			relatedTableName := Pluralize(strcase.ToSnakeCaseLower(relatedModelName))
			sourceRef = fmt.Sprintf("%s.%s", relatedTableName, strcase.ToSnakeCaseLower(relatedFieldName))

			// Record that we need a join to this table
			joinTables[relatedTableName] = true

			// Record the relationship to set up join condition
			joinTableRelationships[relatedTableName] = relatedModelName
		}

		// Add column to the view
		column := psqldef.ViewColumn{
			Name:      columnName,
			SourceRef: sourceRef,
			Alias:     "", // No alias by default
		}
		view.Columns = append(view.Columns, column)
	}

	// Set up joins based on relationships
	for joinTable, _ := range joinTables {
		// Get related model name
		relatedModelName := joinTableRelationships[joinTable]
		if relatedModelName == "" {
			continue
		}

		// Find relationship in model
		modelName := entity.Name
		model, modelErr := r.GetModel(modelName)
		if modelErr != nil {
			return nil, modelErr
		}

		_, relationshipExists := model.Related[relatedModelName]
		if !relationshipExists {
			return nil, fmt.Errorf("relationship %s not found in model %s", relatedModelName, modelName)
		}

		joinType := "LEFT"
		// TODO: We can't just use uuid, we need to extract the primary identifier field from each.

		rootPrimaryId, rootPrimaryIdExists := model.Identifiers["primary"]
		if !rootPrimaryIdExists {
			return nil, fmt.Errorf("primary identifier not found in model '%s'", modelName)
		}

		relatedPrimaryId, relatedPrimaryIdExists := model.Identifiers["primary"]
		if !relatedPrimaryIdExists {
			return nil, fmt.Errorf("primary identifier not found in model '%s'", relatedModelName)
		}
		rootPrimaryIdName := strcase.ToSnakeCaseLower(rootPrimaryId.Fields[0])
		relatedPrimaryIdName := strcase.ToSnakeCaseLower(relatedPrimaryId.Fields[0])

		joinClause := psqldef.JoinClause{
			Type:   joinType,
			Schema: config.MorpheModelsConfig.Schema,
			Table:  joinTable,
			Alias:  joinTable,
			Conditions: []psqldef.JoinCondition{
				{
					LeftRef:  tableName + "." + rootPrimaryIdName,
					RightRef: joinTable + "." + relatedPrimaryIdName,
				},
			},
		}

		view.Joins = append(view.Joins, joinClause)
	}

	return view, nil
}

func triggerCompileMorpheEntityStart(hooks hook.CompileMorpheEntity, config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error) {
	if hooks.OnCompileMorpheEntityStart == nil {
		return config, entity, nil
	}

	return hooks.OnCompileMorpheEntityStart(config, entity)
}

func triggerCompileMorpheEntitySuccess(hooks hook.CompileMorpheEntity, view *psqldef.View) (*psqldef.View, error) {
	if hooks.OnCompileMorpheEntitySuccess == nil {
		return view, nil
	}

	return hooks.OnCompileMorpheEntitySuccess(view)
}

func triggerCompileMorpheEntityFailure(hooks hook.CompileMorpheEntity, config cfg.MorpheConfig, entity yaml.Entity, failureErr error) error {
	if hooks.OnCompileMorpheEntityFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheEntityFailure(config, entity, failureErr)
}

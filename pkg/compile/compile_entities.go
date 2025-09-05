package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
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

	validateEntityErr := entity.Validate(r.GetAllEntities(), r.GetAllModels(), r.GetAllEnums())
	if validateEntityErr != nil {
		return nil, validateEntityErr
	}

	viewName := strcase.ToSnakeCaseLower(entity.Name)
	if config.MorpheEntitiesConfig.ViewNameSuffix != "" {
		viewName += config.MorpheEntitiesConfig.ViewNameSuffix
	}

	tableName := Pluralize(strcase.ToSnakeCaseLower(entity.Name))

	view := &psqldef.View{
		Schema:     config.MorpheEntitiesConfig.Schema,
		Name:       viewName,
		FromSchema: config.MorpheEntitiesConfig.Schema,
		FromTable:  tableName,
		Columns:    []psqldef.ViewColumn{},
		Joins:      []psqldef.JoinClause{},
	}

	context := &entityCompileContext{
		config:                 config,
		registry:               r,
		entity:                 entity,
		view:                   view,
		tableName:              tableName,
		joinTables:             make(map[string]bool),
		joinTableRelationships: make(map[string]string),
	}

	if err := processEntityFields(context); err != nil {
		return nil, err
	}

	if err := setupJoinsForRegularRelationships(context); err != nil {
		return nil, err
	}

	return view, nil
}

// entityCompileContext holds all the context needed for entity compilation
type entityCompileContext struct {
	config                 cfg.MorpheConfig
	registry               *registry.Registry
	entity                 yaml.Entity
	view                   *psqldef.View
	tableName              string
	joinTables             map[string]bool
	joinTableRelationships map[string]string
}

// processEntityFields processes all entity fields and adds appropriate columns to the view
func processEntityFields(ctx *entityCompileContext) error {
	fieldNames := core.MapKeysSorted(ctx.entity.Fields)
	for _, fieldName := range fieldNames {
		field := ctx.entity.Fields[fieldName]
		columnName := strcase.ToSnakeCaseLower(fieldName)

		if err := processEntityField(ctx, fieldName, field, columnName); err != nil {
			return err
		}
	}
	return nil
}

// processEntityField processes a single entity field and handles all field type variations
func processEntityField(ctx *entityCompileContext, fieldName string, field yaml.EntityField, columnName string) error {
	fieldParts := strings.Split(string(field.Type), ".")
	if len(fieldParts) < 2 {
		return fmt.Errorf("invalid field type format: %s", field.Type)
	}

	// The last part is always the field name, everything before represents relationship chain
	relationshipChain := fieldParts[:len(fieldParts)-1]
	targetFieldName := fieldParts[len(fieldParts)-1]

	return processFieldPath(ctx, fieldName, relationshipChain, targetFieldName, columnName)
}

// processFieldPath processes a field path with arbitrary relationship chain depth
func processFieldPath(ctx *entityCompileContext, fieldName string, relationshipChain []string, targetFieldName string, columnName string) error {
	// For field paths like "User.UUID" or "User.Child.Name" or "User.Child.Grandchild.Field"
	// The first element should match the entity model name, then subsequent elements are relationships

	if len(relationshipChain) == 0 {
		return fmt.Errorf("invalid field path: no model specified")
	}

	rootModelName := relationshipChain[0]

	// If there's only one element, this is a direct field reference (e.g., "User.UUID")
	if len(relationshipChain) == 1 {
		return addRegularColumn(ctx, columnName, ctx.tableName, targetFieldName)
	}

	// More than one element means we have relationships to traverse
	currentModelName := rootModelName
	currentTableName := ctx.tableName

	// Start from index 1 since index 0 is the root model name
	for i := 1; i < len(relationshipChain); i++ {
		relationName := relationshipChain[i]

		currentModel, err := ctx.registry.GetModel(currentModelName)
		if err != nil {
			return err
		}

		relation, relationExists := currentModel.Related[relationName]
		if !relationExists {
			return fmt.Errorf("relationship %s not found in model %s", relationName, currentModelName)
		}

		// Handle polymorphic relationships
		if yamlops.IsRelationPoly(relation.Type) {
			// Pass the remaining chain starting from this relationship
			remainingChain := relationshipChain[i:]
			return handlePolymorphicFieldPath(ctx, fieldName, relationName, remainingChain, targetFieldName, columnName, currentModel)
		}

		// Handle regular relationships - set up join and continue traversal
		// Use relationName for table naming to maintain backward compatibility
		relatedTableName := Pluralize(strcase.ToSnakeCaseLower(relationName))

		// Record that we need a join to this table
		ctx.joinTables[relatedTableName] = true
		// Store the relationship name for join setup (keeping existing behavior)
		ctx.joinTableRelationships[relatedTableName] = relationName

		// Update current context for next iteration - use the actual target model
		targetModelName := yamlops.GetRelationTargetName(relationName, relation.Aliased)
		currentModelName = targetModelName
		currentTableName = relatedTableName
	}

	// We've traversed all relationships, now add the final field column
	return addRegularColumn(ctx, columnName, currentTableName, targetFieldName)
}

// handlePolymorphicFieldPath handles field paths that encounter polymorphic relationships
func handlePolymorphicFieldPath(ctx *entityCompileContext, fieldName string, relationName string, remainingChain []string, targetFieldName string, columnName string, currentModel yaml.Model) error {
	// Get the polymorphic relationship
	relation := currentModel.Related[relationName]

	// For polymorphic relationships, we need to determine what to do based on type and position in chain
	if len(remainingChain) == 1 {
		// This is a direct reference to a polymorphic relationship (e.g., "User.Commentable")
		// The relationName is the first (and only) element in remaining chain
		return handlePolymorphicRelationshipColumns(ctx, relationName, columnName, string(relation.Type))
	}

	// This is a nested field through a polymorphic relationship (e.g., "User.Commentable.SomeField.Name")
	// For nested polymorphic paths, we have different handling based on relationship type
	if yamlops.IsRelationPolyFor(relation.Type) && yamlops.IsRelationPolyOne(relation.Type) {
		// ForOnePoly: Include raw polymorphic columns (type + id)
		// We cannot traverse further through polymorphic relationships in entity views
		return addPolymorphicColumns(ctx, relationName, columnName)
	}

	if yamlops.IsRelationPolyFor(relation.Type) && yamlops.IsRelationPolyMany(relation.Type) {
		// ForManyPoly: Keep simple, no junction table materialization
		return nil
	}

	if yamlops.IsRelationPolyHas(relation.Type) {
		// HasOnePoly/HasManyPoly: Not materialized in entity views
		return nil
	}

	return nil
}

// handlePolymorphicRelationshipColumns creates polymorphic type and id columns
func handlePolymorphicRelationshipColumns(ctx *entityCompileContext, relationName, columnName string, relationType string) error {
	if yamlops.IsRelationPolyFor(relationType) && yamlops.IsRelationPolyOne(relationType) {
		// ForOnePoly: Include raw polymorphic columns (type + id)
		return addPolymorphicColumns(ctx, relationName, columnName)
	}

	if yamlops.IsRelationPolyFor(relationType) && yamlops.IsRelationPolyMany(relationType) {
		// ForManyPoly: Keep simple, no junction table materialization
		return nil
	}

	if yamlops.IsRelationPolyHas(relationType) {
		// HasOnePoly/HasManyPoly: Not materialized in entity views
		return nil
	}

	return nil
}

// addPolymorphicColumns adds both type and id columns for a polymorphic relationship
func addPolymorphicColumns(ctx *entityCompileContext, relationName, baseColumnName string) error {
	typeColumnName := baseColumnName + "_type"
	idColumnName := baseColumnName + "_id"

	dbRelationName := strcase.ToSnakeCaseLower(relationName)
	typeSourceRef := fmt.Sprintf("%s.%s_type", ctx.tableName, dbRelationName)
	idSourceRef := fmt.Sprintf("%s.%s_id", ctx.tableName, dbRelationName)

	// Add type column
	typeColumn := psqldef.ViewColumn{
		Name:      typeColumnName,
		SourceRef: typeSourceRef,
		Alias:     "",
	}
	ctx.view.Columns = append(ctx.view.Columns, typeColumn)

	// Add id column
	idColumn := psqldef.ViewColumn{
		Name:      idColumnName,
		SourceRef: idSourceRef,
		Alias:     "",
	}
	ctx.view.Columns = append(ctx.view.Columns, idColumn)

	return nil
}

// addRegularColumn adds a regular column to the view
func addRegularColumn(ctx *entityCompileContext, columnName, tableName, fieldName string) error {
	sourceRef := fmt.Sprintf("%s.%s", tableName, strcase.ToSnakeCaseLower(fieldName))

	column := psqldef.ViewColumn{
		Name:      columnName,
		SourceRef: sourceRef,
		Alias:     "",
	}
	ctx.view.Columns = append(ctx.view.Columns, column)

	return nil
}

// setupJoinsForRegularRelationships sets up joins for all recorded regular relationships
func setupJoinsForRegularRelationships(ctx *entityCompileContext) error {
	for joinTable := range ctx.joinTables {
		relatedModelName := ctx.joinTableRelationships[joinTable]
		if relatedModelName == "" {
			continue
		}

		if err := addJoinClause(ctx, joinTable, relatedModelName); err != nil {
			return err
		}
	}
	return nil
}

// addJoinClause adds a single join clause to the view
func addJoinClause(ctx *entityCompileContext, joinTable, relatedModelName string) error {
	model, err := ctx.registry.GetModel(ctx.entity.Name)
	if err != nil {
		return err
	}

	_, relationshipExists := model.Related[relatedModelName]
	if !relationshipExists {
		return fmt.Errorf("relationship %s not found in model %s", relatedModelName, ctx.entity.Name)
	}

	rootPrimaryId, rootExists := model.Identifiers["primary"]
	if !rootExists {
		return fmt.Errorf("primary identifier not found in model '%s'", ctx.entity.Name)
	}

	relatedPrimaryId, relatedExists := model.Identifiers["primary"]
	if !relatedExists {
		return fmt.Errorf("primary identifier not found in model '%s'", relatedModelName)
	}

	rootPrimaryIdName := strcase.ToSnakeCaseLower(rootPrimaryId.Fields[0])
	relatedPrimaryIdName := strcase.ToSnakeCaseLower(relatedPrimaryId.Fields[0])

	joinClause := psqldef.JoinClause{
		Type:   "LEFT",
		Schema: ctx.config.MorpheModelsConfig.Schema,
		Table:  joinTable,
		Alias:  joinTable,
		Conditions: []psqldef.JoinCondition{
			{
				LeftRef:  ctx.tableName + "." + rootPrimaryIdName,
				RightRef: joinTable + "." + relatedPrimaryIdName,
			},
		},
	}

	ctx.view.Joins = append(ctx.view.Joins, joinClause)
	return nil
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

// handleRelationshipReference handles when an entity field directly references a relationship
func handleRelationshipReference(ctx *entityCompileContext, fieldName, relationName string, relation yaml.ModelRelation, columnName string) error {
	if !yamlops.IsRelationPoly(relation.Type) {
		return fmt.Errorf("entity field '%s' cannot reference non-polymorphic relationship '%s' directly", fieldName, relationName)
	}

	return handlePolymorphicRelationshipColumns(ctx, relationName, columnName, string(relation.Type))
}

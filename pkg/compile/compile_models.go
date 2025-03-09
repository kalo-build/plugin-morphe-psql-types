package compile

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kaloseia/clone"
	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go-util/strcase"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/morphe-go/pkg/yamlops"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/typemap"
)

func AllMorpheModelsToPSQLTables(config MorpheCompileConfig, r *registry.Registry) (map[string][]*psqldef.Table, error) {
	allModelStructDefs := map[string][]*psqldef.Table{}
	for modelName, model := range r.GetAllModels() {
		modelStructs, modelErr := MorpheModelToPSQLTables(config, r, model)
		if modelErr != nil {
			return nil, modelErr
		}
		allModelStructDefs[modelName] = modelStructs
	}
	return allModelStructDefs, nil
}

func MorpheModelToPSQLTables(config MorpheCompileConfig, r *registry.Registry, model yaml.Model) ([]*psqldef.Table, error) {
	morpheConfig, model, compileStartErr := triggerCompileMorpheModelStart(config.ModelHooks, config.MorpheConfig, model)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, compileStartErr)
	}
	config.MorpheConfig = morpheConfig

	allModelStructs, structsErr := morpheModelToPSQLTables(config.MorpheConfig, r, model)
	if structsErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, structsErr)
	}

	allModelStructs, compileSuccessErr := triggerCompileMorpheModelSuccess(config.ModelHooks, allModelStructs)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, compileSuccessErr)
	}
	return allModelStructs, nil
}

func morpheModelToPSQLTables(config cfg.MorpheConfig, r *registry.Registry, model yaml.Model) ([]*psqldef.Table, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := model.Validate(r.GetAllEnums())
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	modelTable, modelTableErr := getModelTable(config, r, model)
	if modelTableErr != nil {
		return nil, modelTableErr
	}
	allModelTables := []*psqldef.Table{
		modelTable,
	}

	return allModelTables, nil
}

func getModelTable(config cfg.MorpheConfig, r *registry.Registry, model yaml.Model) (*psqldef.Table, error) {
	if r == nil {
		return nil, ErrNoRegistry
	}

	tableName := getTableNameFromModel(model.Name)
	modelTable := psqldef.Table{
		Schema:            config.MorpheModelsConfig.Schema,
		Name:              tableName,
		Columns:           []psqldef.Column{},
		Indices:           []psqldef.Index{},
		ForeignKeys:       []psqldef.ForeignKey{},
		UniqueConstraints: []psqldef.UniqueConstraint{},
	}
	primaryID, primaryIDExists := model.Identifiers["primary"]
	if !primaryIDExists {
		return nil, fmt.Errorf("no primary identifier set for model '%s'", model.Name)
	}
	typeMap := typemap.MorpheModelFieldToPSQLField
	if config.UseBigSerial {
		typeMap = typemap.MorpheModelFieldToPSQLFieldBigSerial
	}
	fieldColumns, fieldColumnsErr := getColumnsForModelFields(typeMap, primaryID, model.Fields)
	if fieldColumnsErr != nil {
		return nil, fieldColumnsErr
	}
	relatedTypeMap := typemap.MorpheModelFieldToPSQLFieldForeign
	if config.UseBigSerial {
		relatedTypeMap = typemap.MorpheModelFieldToPSQLFieldBigSerialForeign
	}

	relatedColumns, relatedColumnsErr := getColumnsForModelRelations(r, relatedTypeMap, model.Related)
	if relatedColumnsErr != nil {
		return nil, relatedColumnsErr
	}
	modelTable.Columns = append(fieldColumns, relatedColumns...)

	// Add foreign key constraints and indices for related models
	foreignKeys, foreignKeysErr := getForeignKeysForModelRelations(config.MorpheModelsConfig.Schema, model.Name, r, model.Related)
	if foreignKeysErr != nil {
		return nil, foreignKeysErr
	}
	modelTable.ForeignKeys = foreignKeys

	// Create corresponding indices for foreign keys
	indices := getIndicesForForeignKeys(tableName, foreignKeys)
	modelTable.Indices = indices

	return &modelTable, nil
}

func getColumnsForModelFields(typeMap map[yaml.ModelFieldType]psqldef.PSQLType, primaryID yaml.ModelIdentifier, modelFields map[string]yaml.ModelField) ([]psqldef.Column, error) {
	columns := []psqldef.Column{}

	modelFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range modelFieldNames {
		field := modelFields[fieldName]
		columnName := getColumnNameFromField(fieldName)

		columnType, supported := typeMap[field.Type]
		if !supported {
			return nil, fmt.Errorf("morphe model field '%s' has unsupported type '%s'", fieldName, field.Type)
		}

		column := psqldef.Column{
			Name:       columnName,
			Type:       columnType,
			NotNull:    false,
			PrimaryKey: slices.Index(primaryID.Fields, fieldName) != -1,
			Default:    "",
		}

		columns = append(columns, column)
	}

	return columns, nil
}

func getColumnsForModelRelations(r *registry.Registry, typeMap map[yaml.ModelFieldType]psqldef.PSQLType, relatedModels map[string]yaml.ModelRelation) ([]psqldef.Column, error) {
	columns := []psqldef.Column{}

	relatedModelNames := core.MapKeysSorted(relatedModels)
	for _, relatedModelName := range relatedModelNames {
		modelRelation := relatedModels[relatedModelName]
		relationType := modelRelation.Type
		relatedModel, modelErr := r.GetModel(relatedModelName)
		if modelErr != nil {
			return nil, modelErr
		}
		primaryID, hasPrimary := relatedModel.Identifiers["primary"]
		if !hasPrimary {
			return nil, fmt.Errorf("related model %s has no primary identifier", relatedModelName)
		}

		if len(primaryID.Fields) != 1 {
			return nil, fmt.Errorf("related entity %s primary identifier must have exactly one field", relatedModelName)
		}

		targetPrimaryIdName := primaryID.Fields[0]
		targetPrimaryIdField, primaryFieldExists := relatedModel.Fields[targetPrimaryIdName]
		if !primaryFieldExists {
			return nil, fmt.Errorf("related entity %s primary identifier field %s not found", relatedModelName, targetPrimaryIdName)
		}

		if yamlops.IsRelationFor(relationType) && yamlops.IsRelationOne(relationType) {
			columnName := fmt.Sprintf(
				`%s_%s`,
				strcase.ToSnakeCaseLower(relatedModelName),
				strcase.ToSnakeCaseLower(targetPrimaryIdName),
			)

			columnType, supported := typeMap[targetPrimaryIdField.Type]
			if !supported {
				return nil, fmt.Errorf("morphe related model field '%s' has unsupported type '%s'", targetPrimaryIdName, targetPrimaryIdField.Type)
			}

			column := psqldef.Column{
				Name:       columnName,
				Type:       columnType,
				NotNull:    false,
				PrimaryKey: false,
				Default:    "",
			}

			columns = append(columns, column)
		}
	}

	return columns, nil
}

func getForeignKeysForModelRelations(schema string, modelName string, r *registry.Registry, relatedModels map[string]yaml.ModelRelation) ([]psqldef.ForeignKey, error) {
	foreignKeys := []psqldef.ForeignKey{}

	relatedModelNames := core.MapKeysSorted(relatedModels)
	for _, relatedModelName := range relatedModelNames {
		modelRelation := relatedModels[relatedModelName]
		relationType := modelRelation.Type
		relatedModel, modelErr := r.GetModel(relatedModelName)
		if modelErr != nil {
			return nil, modelErr
		}
		primaryID, hasPrimary := relatedModel.Identifiers["primary"]
		if !hasPrimary {
			return nil, fmt.Errorf("related model %s has no primary identifier", relatedModelName)
		}

		if len(primaryID.Fields) != 1 {
			return nil, fmt.Errorf("related entity %s primary identifier must have exactly one field", relatedModelName)
		}

		targetPrimaryIdName := primaryID.Fields[0]
		_, primaryFieldExists := relatedModel.Fields[targetPrimaryIdName]
		if !primaryFieldExists {
			return nil, fmt.Errorf("related entity %s primary identifier field %s not found", relatedModelName, targetPrimaryIdName)
		}

		if yamlops.IsRelationFor(relationType) && yamlops.IsRelationOne(relationType) {
			columnName := fmt.Sprintf(
				`%s_%s`,
				strcase.ToSnakeCaseLower(relatedModelName),
				strcase.ToSnakeCaseLower(targetPrimaryIdName),
			)

			refTableName := getTableNameFromModel(relatedModelName)
			refColumnName := getColumnNameFromField(targetPrimaryIdName)

			foreignKey := psqldef.ForeignKey{
				Schema:         schema,
				Name:           fmt.Sprintf("fk_%s_%s", strcase.ToSnakeCaseLower(modelName), columnName),
				TableName:      getTableNameFromModel(modelName),
				ColumnNames:    []string{columnName},
				RefTableName:   refTableName,
				RefColumnNames: []string{refColumnName},
				OnDelete:       "CASCADE", // Using CASCADE as per the spec examples
				OnUpdate:       "",        // Default behavior
			}

			foreignKeys = append(foreignKeys, foreignKey)
		}
	}

	return foreignKeys, nil
}

func getIndicesForForeignKeys(tableName string, foreignKeys []psqldef.ForeignKey) []psqldef.Index {
	indices := []psqldef.Index{}

	for _, fk := range foreignKeys {
		for _, columnName := range fk.ColumnNames {
			indexName := fmt.Sprintf("idx_%s_%s", tableName, columnName)

			index := psqldef.Index{
				Name:      indexName,
				TableName: tableName,
				Columns:   []string{columnName},
				IsUnique:  false,
			}

			indices = append(indices, index)
		}
	}

	return indices
}

func getTableNameFromModel(modelName string) string {
	snakeCase := strcase.ToSnakeCaseLower(modelName)

	// Pluralize (simple pluralization)
	if strings.HasSuffix(snakeCase, "s") || strings.HasSuffix(snakeCase, "x") ||
		strings.HasSuffix(snakeCase, "z") || strings.HasSuffix(snakeCase, "ch") ||
		strings.HasSuffix(snakeCase, "sh") {
		return snakeCase + "es"
	} else if strings.HasSuffix(snakeCase, "y") {
		return snakeCase[:len(snakeCase)-1] + "ies"
	} else {
		return snakeCase + "s"
	}
}

func getColumnNameFromField(fieldName string) string {
	snakeCase := strcase.ToSnakeCaseLower(fieldName)
	return snakeCase
}

func triggerCompileMorpheModelStart(modelHooks hook.CompileMorpheModel, config cfg.MorpheConfig, model yaml.Model) (cfg.MorpheConfig, yaml.Model, error) {
	if modelHooks.OnCompileMorpheModelStart == nil {
		return config, model, nil
	}

	updatedConfig, updatedModel, startErr := modelHooks.OnCompileMorpheModelStart(config, model)
	if startErr != nil {
		return cfg.MorpheConfig{}, yaml.Model{}, startErr
	}

	return updatedConfig, updatedModel, nil
}

func triggerCompileMorpheModelSuccess(hooks hook.CompileMorpheModel, allModelStructs []*psqldef.Table) ([]*psqldef.Table, error) {
	if hooks.OnCompileMorpheModelSuccess == nil {
		return allModelStructs, nil
	}
	if allModelStructs == nil {
		return nil, ErrNoModelStructs
	}
	allModelStructsClone := clone.DeepCloneSlicePointers(allModelStructs)

	allModelStructs, successErr := hooks.OnCompileMorpheModelSuccess(allModelStructsClone)
	if successErr != nil {
		return nil, successErr
	}
	return allModelStructs, nil
}

func triggerCompileMorpheModelFailure(hooks hook.CompileMorpheModel, morpheConfig cfg.MorpheConfig, model yaml.Model, failureErr error) error {
	if hooks.OnCompileMorpheModelFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheModelFailure(morpheConfig, model.DeepClone(), failureErr)
}

package compile

import (
	"fmt"
	"slices"

	"github.com/kaloseia/clone"
	"github.com/kaloseia/go-util/core"
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

	// Name of the model
	schema := config.MorpheModelsConfig.Schema
	modelName := model.Name
	tableName := GetTableNameFromModel(modelName)

	// Get the appropriate type maps based on UseBigSerial config
	var typeMap map[yaml.ModelFieldType]psqldef.PSQLType
	var relatedTypeMap map[yaml.ModelFieldType]psqldef.PSQLType

	if config.MorpheModelsConfig.UseBigSerial {
		typeMap = typemap.MorpheModelFieldToPSQLFieldBigSerial
		relatedTypeMap = typemap.MorpheModelFieldToPSQLFieldBigSerialForeign
	} else {
		typeMap = typemap.MorpheModelFieldToPSQLField
		relatedTypeMap = typemap.MorpheModelFieldToPSQLFieldForeign
	}

	primaryID, primaryIDExists := model.Identifiers["primary"]
	if !primaryIDExists {
		return nil, fmt.Errorf("no primary identifier set for model '%s'", model.Name)
	}

	fieldColumns, fieldColumnsErr := getColumnsForModelFields(typeMap, primaryID, model.Fields)
	if fieldColumnsErr != nil {
		return nil, fieldColumnsErr
	}

	relatedColumns, relatedColumnsErr := getColumnsForModelRelations(r, relatedTypeMap, model.Related)
	if relatedColumnsErr != nil {
		return nil, relatedColumnsErr
	}

	modelTable := psqldef.Table{
		Schema:            schema,
		Name:              tableName,
		Columns:           append(fieldColumns, relatedColumns...),
		ForeignKeys:       []psqldef.ForeignKey{},
		Indices:           []psqldef.Index{},
		UniqueConstraints: []psqldef.UniqueConstraint{},
	}

	// Add foreign key constraints and indices for related models
	foreignKeys, foreignKeysErr := getForeignKeysForModelRelations(schema, tableName, r, model.Related)
	if foreignKeysErr != nil {
		return nil, foreignKeysErr
	}
	modelTable.ForeignKeys = foreignKeys

	// Create corresponding indices for foreign keys
	indices := getIndicesForForeignKeys(tableName, foreignKeys)
	modelTable.Indices = indices

	// Create junction tables for any ForMany relationships
	junctionTables, junctionTablesErr := getJunctionTablesForForManyRelations(schema, r, model)
	if junctionTablesErr != nil {
		return nil, junctionTablesErr
	}

	// Return the main model table and any junction tables
	tables := []*psqldef.Table{&modelTable}
	if len(junctionTables) > 0 {
		tables = append(tables, junctionTables...)
	}

	return tables, nil
}

func getColumnsForModelFields(typeMap map[yaml.ModelFieldType]psqldef.PSQLType, primaryID yaml.ModelIdentifier, modelFields map[string]yaml.ModelField) ([]psqldef.Column, error) {
	columns := []psqldef.Column{}

	modelFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range modelFieldNames {
		field := modelFields[fieldName]
		columnName := GetColumnNameFromField(fieldName)

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
			columnName := GetForeignKeyColumnName(relatedModelName, targetPrimaryIdName)

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

func getForeignKeysForModelRelations(schema string, tableName string, r *registry.Registry, relatedModels map[string]yaml.ModelRelation) ([]psqldef.ForeignKey, error) {
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
			columnName := GetForeignKeyColumnName(relatedModelName, targetPrimaryIdName)
			refTableName := GetTableNameFromModel(relatedModelName)
			refColumnName := GetColumnNameFromField(targetPrimaryIdName)

			foreignKey := psqldef.ForeignKey{
				Schema:         schema,
				Name:           GetForeignKeyConstraintName(tableName, columnName),
				TableName:      tableName,
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
			index := psqldef.Index{
				Name:      GetIndexName(tableName, columnName),
				TableName: tableName,
				Columns:   []string{columnName},
				IsUnique:  false,
			}

			indices = append(indices, index)
		}
	}

	return indices
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

// getJunctionTablesForForManyRelations creates junction tables for ForMany relationships
func getJunctionTablesForForManyRelations(schema string, r *registry.Registry, model yaml.Model) ([]*psqldef.Table, error) {
	junctionTables := []*psqldef.Table{}
	modelName := model.Name
	tableName := GetTableNameFromModel(modelName)

	// Get primary ID field for this model
	primaryID, hasPrimary := model.Identifiers["primary"]
	if !hasPrimary {
		return nil, fmt.Errorf("model %s has no primary identifier", modelName)
	}
	if len(primaryID.Fields) != 1 {
		return nil, fmt.Errorf("model %s primary identifier must have exactly one field", modelName)
	}
	primaryIdName := primaryID.Fields[0]

	relatedModelNames := core.MapKeysSorted(model.Related)
	for _, relatedModelName := range relatedModelNames {
		modelRelation := model.Related[relatedModelName]
		relationType := modelRelation.Type

		if yamlops.IsRelationFor(relationType) && yamlops.IsRelationMany(relationType) {
			relatedModel, modelErr := r.GetModel(relatedModelName)
			if modelErr != nil {
				return nil, modelErr
			}

			// Get primary ID field for related model
			relatedPrimaryID, hasRelatedPrimary := relatedModel.Identifiers["primary"]
			if !hasRelatedPrimary {
				return nil, fmt.Errorf("related model %s has no primary identifier", relatedModelName)
			}
			if len(relatedPrimaryID.Fields) != 1 {
				return nil, fmt.Errorf("related model %s primary identifier must have exactly one field", relatedModelName)
			}
			relatedPrimaryIdName := relatedPrimaryID.Fields[0]

			// Create junction table
			junctionTableName := GetJunctionTableName(modelName, relatedModelName)

			// Create column names
			sourceColumnName := GetForeignKeyColumnName(modelName, primaryIdName)
			targetColumnName := GetForeignKeyColumnName(relatedModelName, relatedPrimaryIdName)

			// Create columns
			columns := []psqldef.Column{
				{
					Name:       "id",
					Type:       psqldef.PSQLTypeSerial,
					PrimaryKey: true,
				},
				{
					Name: sourceColumnName,
					Type: psqldef.PSQLTypeInteger,
				},
				{
					Name: targetColumnName,
					Type: psqldef.PSQLTypeInteger,
				},
			}

			// Create foreign keys
			foreignKeys := []psqldef.ForeignKey{
				{
					Schema:       schema,
					Name:         GetJunctionTableForeignKeyConstraintName(junctionTableName, modelName, primaryIdName),
					TableName:    junctionTableName,
					ColumnNames:  []string{sourceColumnName},
					RefTableName: tableName,
					RefColumnNames: []string{
						GetColumnNameFromField(primaryIdName),
					},
					OnDelete: "CASCADE",
				},
				{
					Schema:       schema,
					Name:         GetJunctionTableForeignKeyConstraintName(junctionTableName, relatedModelName, relatedPrimaryIdName),
					TableName:    junctionTableName,
					ColumnNames:  []string{targetColumnName},
					RefTableName: GetTableNameFromModel(relatedModelName),
					RefColumnNames: []string{
						GetColumnNameFromField(relatedPrimaryIdName),
					},
					OnDelete: "CASCADE",
				},
			}

			// Create unique constraint
			uniqueConstraints := []psqldef.UniqueConstraint{
				{
					Name: GetJunctionTableUniqueConstraintName(
						junctionTableName,
						modelName, primaryIdName,
						relatedModelName, relatedPrimaryIdName,
					),
					TableName: junctionTableName,
					ColumnNames: []string{
						sourceColumnName,
						targetColumnName,
					},
				},
			}

			// Create indices for foreign keys
			indices := getIndicesForForeignKeys(junctionTableName, foreignKeys)

			// Create junction table
			junctionTable := &psqldef.Table{
				Schema:            schema,
				Name:              junctionTableName,
				Columns:           columns,
				ForeignKeys:       foreignKeys,
				Indices:           indices,
				UniqueConstraints: uniqueConstraints,
			}

			junctionTables = append(junctionTables, junctionTable)
		}
	}

	return junctionTables, nil
}

package compile

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kalo-build/clone"
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/typemap"
)

func AllMorpheModelsToPSQLTables(config MorpheCompileConfig, r *registry.Registry) (map[string][]*psqldef.Table, error) {
	allModelTableDefs := map[string][]*psqldef.Table{}
	for modelName, model := range r.GetAllModels() {
		modelTables, modelErr := MorpheModelToPSQLTables(config, r, model)
		if modelErr != nil {
			return nil, modelErr
		}
		allModelTableDefs[modelName] = modelTables
	}
	return allModelTableDefs, nil
}

func MorpheModelToPSQLTables(config MorpheCompileConfig, r *registry.Registry, model yaml.Model) ([]*psqldef.Table, error) {
	morpheConfig, model, compileStartErr := triggerCompileMorpheModelStart(config.ModelHooks, config.MorpheConfig, model)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, compileStartErr)
	}
	config.MorpheConfig = morpheConfig

	allModelTables, tablesErr := morpheModelToPSQLTables(config.MorpheConfig, r, model)
	if tablesErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, tablesErr)
	}

	allModelTables, compileSuccessErr := triggerCompileMorpheModelSuccess(config.ModelHooks, allModelTables)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, compileSuccessErr)
	}
	return allModelTables, nil
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

	validatePolyErr := validatePolymorphicRelationships(r, model)
	if validatePolyErr != nil {
		return nil, validatePolyErr
	}

	schema := config.MorpheModelsConfig.Schema
	modelName := model.Name
	tableName := GetTableNameFromModel(modelName)

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

	fieldColumns, enumForeignKeys, fieldColumnsErr := getColumnsForModelFields(config, r, typeMap, tableName, primaryID, model.Fields)
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
		ForeignKeys:       enumForeignKeys,
		Indices:           []psqldef.Index{},
		UniqueConstraints: []psqldef.UniqueConstraint{},
	}

	relationForeignKeys, foreignKeysErr := getForeignKeysForModelRelations(schema, tableName, r, model.Related)
	if foreignKeysErr != nil {
		return nil, foreignKeysErr
	}
	modelTable.ForeignKeys = append(modelTable.ForeignKeys, relationForeignKeys...)

	indices := getIndicesForForeignKeys(schema, tableName, modelTable.ForeignKeys)
	modelTable.Indices = indices

	// Apply spec-compliant processing to the model table
	addUniqueIndicesFromIdentifiers(&modelTable, model.Identifiers)
	quoteReservedColumnNames(&modelTable)
	ensureNamedForeignKeyConstraints(&modelTable)

	junctionTables, junctionTablesErr := getJunctionTablesForForManyRelations(schema, r, model)
	if junctionTablesErr != nil {
		return nil, junctionTablesErr
	}

	// Get polymorphic junction tables for ForManyPoly relationships
	polymorphicJunctionTables, polymorphicJunctionTablesErr := getJunctionTablesForForManyPolyRelations(schema, r, model)
	if polymorphicJunctionTablesErr != nil {
		return nil, polymorphicJunctionTablesErr
	}

	// Combine all junction tables
	allJunctionTables := append(junctionTables, polymorphicJunctionTables...)

	// Process junction tables as well
	for tableIdx := range allJunctionTables {
		quoteReservedColumnNames(allJunctionTables[tableIdx])
		ensureNamedForeignKeyConstraints(allJunctionTables[tableIdx])
	}

	tables := []*psqldef.Table{&modelTable}
	if len(allJunctionTables) > 0 {
		tables = append(tables, allJunctionTables...)
	}

	return tables, nil
}

func getColumnsForModelFields(config cfg.MorpheConfig, r *registry.Registry, typeMap map[yaml.ModelFieldType]psqldef.PSQLType, tableName string, primaryID yaml.ModelIdentifier, modelFields map[string]yaml.ModelField) ([]psqldef.TableColumn, []psqldef.ForeignKey, error) {
	columns := []psqldef.TableColumn{}
	enumForeignKeys := []psqldef.ForeignKey{}

	modelFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range modelFieldNames {
		field := modelFields[fieldName]
		columnName := GetColumnNameFromField(fieldName)

		columnType, supported := typeMap[field.Type]
		if supported {
			column := psqldef.TableColumn{
				Name:       columnName,
				Type:       columnType,
				NotNull:    false,
				PrimaryKey: slices.Index(primaryID.Fields, fieldName) != -1,
				Default:    "",
			}
			columns = append(columns, column)
			continue
		}

		enumType, enumErr := r.GetEnum(string(field.Type))
		if enumErr != nil {
			return nil, nil, fmt.Errorf("morphe model field '%s' has unsupported type '%s'", fieldName, field.Type)
		}

		columnName = columnName + "_id"
		enumTableName := Pluralize(strcase.ToSnakeCaseLower(enumType.Name))

		foreignKey := psqldef.ForeignKey{
			Schema:         config.MorpheModelsConfig.Schema,
			Name:           GetForeignKeyConstraintName(tableName, columnName),
			TableName:      tableName,
			ColumnNames:    []string{columnName},
			RefSchema:      config.MorpheEnumsConfig.Schema,
			RefTableName:   enumTableName,
			RefColumnNames: []string{"id"},
			OnDelete:       "CASCADE",
			OnUpdate:       "",
		}
		enumForeignKeys = append(enumForeignKeys, foreignKey)

		column := psqldef.TableColumn{
			Name:       columnName,
			Type:       psqldef.PSQLTypeInteger,
			NotNull:    true,
			PrimaryKey: slices.Index(primaryID.Fields, fieldName) != -1,
			Default:    "",
		}
		columns = append(columns, column)
	}

	return columns, enumForeignKeys, nil
}

func getColumnsForModelRelations(r *registry.Registry, typeMap map[yaml.ModelFieldType]psqldef.PSQLType, relatedModels map[string]yaml.ModelRelation) ([]psqldef.TableColumn, error) {
	columns := []psqldef.TableColumn{}

	relatedModelNames := core.MapKeysSorted(relatedModels)
	for _, relatedModelName := range relatedModelNames {
		modelRelation := relatedModels[relatedModelName]
		relationType := modelRelation.Type

		if yamlops.IsRelationPolyFor(relationType) && yamlops.IsRelationPolyOne(relationType) {
			typeColumnName := strcase.ToSnakeCaseLower(relatedModelName) + "_type"
			typeColumn := psqldef.TableColumn{
				Name:       typeColumnName,
				Type:       psqldef.PSQLTypeText,
				NotNull:    true,
				PrimaryKey: false,
				Default:    "",
			}
			columns = append(columns, typeColumn)

			idColumnName := strcase.ToSnakeCaseLower(relatedModelName) + "_id"
			idColumn := psqldef.TableColumn{
				Name:       idColumnName,
				Type:       psqldef.PSQLTypeText,
				NotNull:    true,
				PrimaryKey: false,
				Default:    "",
			}
			columns = append(columns, idColumn)
			continue
		}

		if yamlops.IsRelationPoly(relationType) {
			continue
		}

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

			column := psqldef.TableColumn{
				Name:       columnName,
				Type:       columnType,
				NotNull:    true,
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

		if yamlops.IsRelationPoly(relationType) {
			continue
		}

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
				RefSchema:      schema,
				RefTableName:   refTableName,
				RefColumnNames: []string{refColumnName},
				OnDelete:       "CASCADE",
				OnUpdate:       "",
			}

			foreignKeys = append(foreignKeys, foreignKey)
		}
	}

	return foreignKeys, nil
}

func getIndicesForForeignKeys(schema string, tableName string, foreignKeys []psqldef.ForeignKey) []psqldef.Index {
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

func triggerCompileMorpheModelSuccess(hooks hook.CompileMorpheModel, allModelTables []*psqldef.Table) ([]*psqldef.Table, error) {
	if hooks.OnCompileMorpheModelSuccess == nil {
		return allModelTables, nil
	}
	if allModelTables == nil {
		return nil, ErrNoModelTables
	}
	allModelTablesClone := clone.DeepCloneSlicePointers(allModelTables)

	allModelTables, successErr := hooks.OnCompileMorpheModelSuccess(allModelTablesClone)
	if successErr != nil {
		return nil, successErr
	}
	return allModelTables, nil
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

		if yamlops.IsRelationFor(relationType) && yamlops.IsRelationMany(relationType) && !yamlops.IsRelationPoly(relationType) {
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
			columns := []psqldef.TableColumn{
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
					RefSchema:    schema,
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
					RefSchema:    schema,
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
			indices := getIndicesForForeignKeys(schema, junctionTableName, foreignKeys)

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

// getJunctionTablesForForManyPolyRelations creates polymorphic junction tables for ForManyPoly relationships
func getJunctionTablesForForManyPolyRelations(schema string, r *registry.Registry, model yaml.Model) ([]*psqldef.Table, error) {
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
	for _, relationName := range relatedModelNames {
		modelRelation := model.Related[relationName]
		relationType := modelRelation.Type

		if yamlops.IsRelationPolyFor(relationType) && yamlops.IsRelationPolyMany(relationType) {
			// Create junction table name - use relation name instead of target model name
			junctionTableName := GetJunctionTableName(modelName, relationName)

			// Create column names
			sourceColumnName := GetForeignKeyColumnName(modelName, primaryIdName)
			typeColumnName := strcase.ToSnakeCaseLower(relationName) + "_type"
			idColumnName := strcase.ToSnakeCaseLower(relationName) + "_id"

			// Create columns
			columns := []psqldef.TableColumn{
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
					Name: typeColumnName,
					Type: psqldef.PSQLTypeText,
				},
				{
					Name: idColumnName,
					Type: psqldef.PSQLTypeText,
				},
			}

			// Create foreign key only for the source model (no FK for polymorphic columns)
			foreignKeys := []psqldef.ForeignKey{
				{
					Schema:       schema,
					Name:         GetJunctionTableForeignKeyConstraintName(junctionTableName, modelName, primaryIdName),
					TableName:    junctionTableName,
					ColumnNames:  []string{sourceColumnName},
					RefSchema:    schema,
					RefTableName: tableName,
					RefColumnNames: []string{
						GetColumnNameFromField(primaryIdName),
					},
					OnDelete: "CASCADE",
				},
			}

			// Create unique constraint on (source_id, target_type, target_id)
			uniqueConstraints := []psqldef.UniqueConstraint{
				{
					Name: GetPolymorphicJunctionTableUniqueConstraintName(
						junctionTableName,
						modelName, primaryIdName,
						relationName,
					),
					TableName: junctionTableName,
					ColumnNames: []string{
						sourceColumnName,
						typeColumnName,
						idColumnName,
					},
				},
			}

			// Create indices for foreign keys
			indices := getIndicesForForeignKeys(schema, junctionTableName, foreignKeys)

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

// addUniqueIndicesFromIdentifiers adds unique indices for model identifiers
func addUniqueIndicesFromIdentifiers(table *psqldef.Table, identifiers map[string]yaml.ModelIdentifier) {
	tableName := table.Name

	// Add unique indices for identifiers
	for idName, identifier := range identifiers {
		if idName == "primary" {
			continue
		}
		columnNames := make([]string, len(identifier.Fields))
		for fieldIdx, field := range identifier.Fields {
			columnNames[fieldIdx] = GetColumnNameFromField(field)
		}

		var indexName string
		if len(columnNames) == 1 {
			indexName = fmt.Sprintf("idx_%s_%s", tableName, columnNames[0])
		} else {
			indexName = fmt.Sprintf("idx_%s_%s", tableName, strings.Join(columnNames, "_"))
		}

		table.Indices = append(table.Indices, psqldef.Index{
			Name:      indexName,
			TableName: tableName,
			Columns:   columnNames,
			IsUnique:  true,
		})
	}
}

// quoteReservedColumnNames quotes column names that are PostgreSQL reserved words
func quoteReservedColumnNames(table *psqldef.Table) {
	reservedWords := map[string]bool{
		"name":    true,
		"type":    true,
		"user":    true,
		"case":    true,
		"when":    true,
		"then":    true,
		"else":    true,
		"end":     true,
		"null":    true,
		"true":    true,
		"false":   true,
		"select":  true,
		"insert":  true,
		"update":  true,
		"delete":  true,
		"from":    true,
		"where":   true,
		"group":   true,
		"order":   true,
		"limit":   true,
		"offset":  true,
		"join":    true,
		"on":      true,
		"using":   true,
		"and":     true,
		"or":      true,
		"not":     true,
		"between": true,
		"alter":   true,
		"table":   true,
		"index":   true,
		"unique":  true,
		"primary": true,
		"foreign": true,
		"key":     true,
	}

	for idxIdx, idx := range table.Indices {
		for colIdx, colName := range idx.Columns {
			if reservedWords[colName] {
				idx.Columns[colIdx] = fmt.Sprintf("\"%s\"", colName)
				table.Indices[idxIdx] = idx
			}
		}
	}
}

// ensureNamedForeignKeyConstraints ensures all foreign keys have proper names and CASCADE behavior
func ensureNamedForeignKeyConstraints(table *psqldef.Table) {
	for fkIdx, fk := range table.ForeignKeys {
		if fk.Name == "" {
			fk.Name = GetForeignKeyConstraintName(table.Name, fk.ColumnNames[0])
			table.ForeignKeys[fkIdx] = fk
		}

		// Ensure CASCADE behavior as per spec
		if fk.OnDelete == "" {
			fk.OnDelete = "CASCADE"
			table.ForeignKeys[fkIdx] = fk
		}
	}
}

// validatePolymorphicRelationships validates polymorphic relationships and their through properties
func validatePolymorphicRelationships(r *registry.Registry, model yaml.Model) error {
	for relationName, relation := range model.Related {
		relationType := relation.Type

		// Validate ForOnePoly and ForManyPoly relationships
		if yamlops.IsRelationPolyFor(relationType) {
			// Check if 'for' property is present and has at least one model
			if len(relation.For) == 0 {
				return fmt.Errorf("polymorphic relation '%s' must have at least one model in 'for' property", relationName)
			}

			// Validate that all models in 'for' property exist in registry
			for _, forModelName := range relation.For {
				_, modelErr := r.GetModel(forModelName)
				if modelErr != nil {
					return fmt.Errorf("model '%s' referenced in 'for' property not found in registry for relation '%s'", forModelName, relationName)
				}
			}

			// Check for circular polymorphic references
			circularErr := checkCircularPolymorphicReferences(r, model.Name, relationName, relation.For)
			if circularErr != nil {
				return circularErr
			}
		}

		// Validate HasOnePoly and HasManyPoly relationships
		if yamlops.IsRelationPolyHas(relationType) {
			if relation.Through == "" {
				return fmt.Errorf("polymorphic relation '%s' of type '%s' must have a 'through' property", relationName, relationType)
			}

			// Look for the through relationship in all models in the registry
			throughFound := false
			throughIsPolymorphic := false

			for _, registryModel := range r.GetAllModels() {
				if throughRelation, exists := registryModel.Related[relation.Through]; exists {
					throughFound = true
					if yamlops.IsRelationPoly(throughRelation.Type) {
						throughIsPolymorphic = true
						break
					}
				}
			}

			if !throughFound {
				return fmt.Errorf("through property '%s' not found in any model for relation '%s'", relation.Through, relationName)
			}

			if !throughIsPolymorphic {
				return fmt.Errorf("through property '%s' must reference a polymorphic relation", relation.Through)
			}
		}
	}

	return nil
}

// checkCircularPolymorphicReferences detects circular polymorphic references
func checkCircularPolymorphicReferences(r *registry.Registry, currentModelName string, relationName string, forModels []string) error {
	for _, targetModelName := range forModels {
		visited := make(map[string]bool)
		path := []string{currentModelName}

		if cyclePath := findCircularReference(r, currentModelName, targetModelName, visited, path); cyclePath != nil {
			return fmt.Errorf("circular polymorphic reference detected, break to prevent infinite loops: %s", formatCyclePath(cyclePath, relationName))
		}
	}

	return nil
}

// findCircularReference performs DFS to detect cycles and returns the full cycle path
func findCircularReference(r *registry.Registry, sourceModel, targetModel string, visited map[string]bool, path []string) []string {
	if targetModel == sourceModel && len(path) > 1 {
		return append(path, targetModel)
	}

	if visited[targetModel] {
		return nil
	}

	visited[targetModel] = true
	newPath := append(path, targetModel)

	targetModelData, err := r.GetModel(targetModel)
	if err != nil {
		return nil
	}

	for _, relation := range targetModelData.Related {
		if yamlops.IsRelationPolyFor(relation.Type) {
			for _, forModel := range relation.For {
				if cyclePath := findCircularReference(r, sourceModel, forModel, visited, newPath); cyclePath != nil {
					return cyclePath
				}
			}
		}
	}

	return nil
}

// formatCyclePath creates a readable representation of the circular reference path
func formatCyclePath(cyclePath []string, initialRelation string) string {
	if len(cyclePath) < 2 {
		return "unknown cycle"
	}

	result := fmt.Sprintf("%s -[%s]-> %s", cyclePath[0], initialRelation, cyclePath[1])
	for i := 1; i < len(cyclePath)-1; i++ {
		result += fmt.Sprintf(" -[polymorphic]-> %s", cyclePath[i+1])
	}

	return result
}

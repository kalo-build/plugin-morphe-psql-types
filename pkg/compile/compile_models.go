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
	columns, columnsErr := getColumnsForModelFields(primaryID, model.Fields)
	if columnsErr != nil {
		return nil, columnsErr
	}
	modelTable.Columns = columns

	return &modelTable, nil
}

func getColumnsForModelFields(primaryID yaml.ModelIdentifier, modelFields map[string]yaml.ModelField) ([]psqldef.Column, error) {
	columns := []psqldef.Column{}

	modelFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range modelFieldNames {
		field := modelFields[fieldName]
		columnName := getColumnNameFromField(fieldName)

		columnType, supported := typemap.MorpheModelFieldToGoField[field.Type]
		if !supported {
			return nil, fmt.Errorf("morphe model field %s has unsupported type: %s", fieldName, field.Type)
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

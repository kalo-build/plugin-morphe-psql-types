package compile

import (
	"fmt"

	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go-util/strcase"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
)

func AllMorpheEnumsToPSQLTables(config MorpheCompileConfig, r *registry.Registry) (map[string]*psqldef.Table, error) {
	allEnumTableDefs := map[string]*psqldef.Table{}
	for enumName, enum := range r.GetAllEnums() {
		enumTable, enumErr := MorpheEnumToPSQLTable(config, enum)
		if enumErr != nil {
			return nil, enumErr
		}
		allEnumTableDefs[enumName] = enumTable
	}
	return allEnumTableDefs, nil
}

// MorpheEnumToPSQLTable converts a Morphe enum to a PostgreSQL lookup table with seed data
func MorpheEnumToPSQLTable(config MorpheCompileConfig, enum yaml.Enum) (*psqldef.Table, error) {
	enumsConfig, enum, enumStartErr := triggerCompileMorpheEnumStart(config.EnumHooks, config.MorpheEnumsConfig, enum)
	if enumStartErr != nil {
		return nil, triggerCompileMorpheEnumFailure(config.EnumHooks, config.MorpheEnumsConfig, enum, enumStartErr)
	}

	table, createPSQLTableForEnumErr := createPSQLTableForEnum(enumsConfig, enum)
	if createPSQLTableForEnumErr != nil {
		return nil, triggerCompileMorpheEnumFailure(config.EnumHooks, enumsConfig, enum, createPSQLTableForEnumErr)
	}

	table, enumSuccessErr := triggerCompileMorpheEnumSuccess(config.EnumHooks, table)
	if enumSuccessErr != nil {
		return nil, triggerCompileMorpheEnumFailure(config.EnumHooks, enumsConfig, enum, enumSuccessErr)
	}

	return table, nil
}

// createPSQLTableForEnum creates a PostgreSQL table with seed data for a Morphe enum
func createPSQLTableForEnum(config cfg.MorpheEnumsConfig, enum yaml.Enum) (*psqldef.Table, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := enum.Validate()
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	tableName := strcase.ToSnakeCaseLower(enum.Name)
	tableName = Pluralize(tableName)

	serialType := psqldef.PSQLTypeSerial
	if config.UseBigSerial {
		serialType = psqldef.PSQLTypeBigSerial
	}

	seedData := psqldef.InsertStatement{
		Schema:    config.Schema,
		TableName: tableName,
		Columns:   []string{"key", "value", "value_type"},
		Values:    [][]any{},
	}

	entryNames := core.MapKeysSorted(enum.Entries)
	for _, key := range entryNames {
		value := enum.Entries[key]
		valueStr := fmt.Sprintf("%v", value)
		valueType := string(enum.Type)

		seedData.Values = append(seedData.Values, []any{
			key,
			valueStr,
			valueType,
		})
	}

	table := &psqldef.Table{
		Schema: config.Schema,
		Name:   tableName,
		Columns: []psqldef.Column{
			{
				Name:       "id",
				Type:       serialType,
				PrimaryKey: true,
			},
			{
				Name:    "key",
				Type:    psqldef.PSQLTypeText,
				NotNull: true,
			},
			{
				Name:    "value",
				Type:    psqldef.PSQLTypeText,
				NotNull: true,
			},
			{
				Name:    "value_type",
				Type:    psqldef.PSQLTypeText,
				NotNull: true,
			},
		},
		UniqueConstraints: []psqldef.UniqueConstraint{
			{
				Name:        "uk_" + tableName + "_key",
				TableName:   tableName,
				ColumnNames: []string{"key"},
			},
		},
		SeedData: []psqldef.InsertStatement{seedData},
	}

	return table, nil
}

// triggerCompileMorpheEnumStart triggers the start hook for enum compilation
func triggerCompileMorpheEnumStart(hooks hook.CompileMorpheEnum, config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error) {
	if hooks.OnCompileMorpheEnumStart == nil {
		return config, enum, nil
	}

	return hooks.OnCompileMorpheEnumStart(config, enum)
}

// triggerCompileMorpheEnumSuccess triggers the success hook for enum compilation
func triggerCompileMorpheEnumSuccess(hooks hook.CompileMorpheEnum, table *psqldef.Table) (*psqldef.Table, error) {
	if hooks.OnCompileMorpheEnumSuccess == nil {
		return table, nil
	}

	tableClone := table.DeepClone()

	updatedTable, err := hooks.OnCompileMorpheEnumSuccess(&tableClone)
	if err != nil {
		return nil, err
	}

	return updatedTable, nil
}

// triggerCompileMorpheEnumFailure triggers the failure hook for enum compilation
func triggerCompileMorpheEnumFailure(hooks hook.CompileMorpheEnum, config cfg.MorpheEnumsConfig, enum yaml.Enum, failureErr error) error {
	if hooks.OnCompileMorpheEnumFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheEnumFailure(config, enum, failureErr)
}

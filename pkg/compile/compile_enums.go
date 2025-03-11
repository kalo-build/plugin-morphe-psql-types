package compile

import (
	"fmt"

	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go-util/strcase"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
)

// MorpheEnumToPSQLTable converts a Morphe enum to a PostgreSQL lookup table and seed data
func MorpheEnumToPSQLTable(config MorpheCompileConfig, enum yaml.Enum) (*psqldef.Table, *psqldef.InsertStatement, error) {
	// Apply the start hook if configured
	var err error
	config.MorpheConfig, enum, err = triggerCompileMorpheEnumStart(config.EnumHooks, config.MorpheConfig, enum)
	if err != nil {
		return nil, nil, triggerCompileMorpheEnumFailure(config.EnumHooks, config.MorpheConfig, enum, err)
	}

	// Validate enum
	if err := validateEnum(config, enum); err != nil {
		return nil, nil, triggerCompileMorpheEnumFailure(config.EnumHooks, config.MorpheConfig, enum, err)
	}

	// Create table and seed data
	table, seedData := createPSQLTableForEnum(config, enum)

	// Apply the success hook
	table, seedData, err = triggerCompileMorpheEnumSuccess(config.EnumHooks, table, seedData)
	if err != nil {
		return nil, nil, triggerCompileMorpheEnumFailure(config.EnumHooks, config.MorpheConfig, enum, err)
	}

	return table, seedData, nil
}

// createPSQLTableForEnum creates a PostgreSQL table and seed data for a Morphe enum
func createPSQLTableForEnum(config MorpheCompileConfig, enum yaml.Enum) (*psqldef.Table, *psqldef.InsertStatement) {
	// Create the table name: singular enum name to plural snake_case
	// e.g., "UserRole" -> "user_roles"
	tableName := strcase.ToSnakeCaseLower(enum.Name)
	tableName = Pluralize(tableName)

	// Create the PostgreSQL table
	table := &psqldef.Table{
		Schema: config.MorpheConfig.MorpheEnumsConfig.Schema,
		Name:   tableName,
		Columns: []psqldef.Column{
			{
				Name:       "id",
				Type:       psqldef.PSQLTypeSerial,
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
		},
		// Add a unique constraint for the key column
		UniqueConstraints: []psqldef.UniqueConstraint{
			{
				Name:        "uk_" + tableName + "_key",
				TableName:   tableName,
				ColumnNames: []string{"key"},
			},
		},
	}

	// Create the INSERT statement for seed data
	seedData := &psqldef.InsertStatement{
		Schema:    config.MorpheConfig.MorpheEnumsConfig.Schema,
		TableName: tableName,
		Columns:   []string{"key", "value"},
		Values:    [][]any{},
	}

	entryNames := core.MapKeysSorted(enum.Entries)
	// Generate values for each enum entry
	for _, key := range entryNames {
		value := enum.Entries[key]
		valueStr := fmt.Sprintf("%v", value)

		seedData.Values = append(seedData.Values, []any{
			key,
			valueStr,
		})
	}

	return table, seedData
}

// triggerCompileMorpheEnumStart triggers the start hook for enum compilation
func triggerCompileMorpheEnumStart(hooks hook.CompileMorpheEnum, config cfg.MorpheConfig, enum yaml.Enum) (cfg.MorpheConfig, yaml.Enum, error) {
	if hooks.OnCompileMorpheEnumStart == nil {
		return config, enum, nil
	}

	return hooks.OnCompileMorpheEnumStart(config, enum)
}

// triggerCompileMorpheEnumSuccess triggers the success hook for enum compilation
func triggerCompileMorpheEnumSuccess(hooks hook.CompileMorpheEnum, table *psqldef.Table, seedData *psqldef.InsertStatement) (*psqldef.Table, *psqldef.InsertStatement, error) {
	if hooks.OnCompileMorpheEnumSuccess == nil {
		return table, seedData, nil
	}

	tableClone := table.DeepClone()
	seedDataClone := seedData.DeepClone()

	table, seedData, err := hooks.OnCompileMorpheEnumSuccess(&tableClone, &seedDataClone)
	if err != nil {
		return nil, nil, err
	}

	return table, seedData, nil
}

// triggerCompileMorpheEnumFailure triggers the failure hook for enum compilation
func triggerCompileMorpheEnumFailure(hooks hook.CompileMorpheEnum, config cfg.MorpheConfig, enum yaml.Enum, failureErr error) error {
	if hooks.OnCompileMorpheEnumFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheEnumFailure(config, enum, failureErr)
}

// validateEnum validates that an enum has all required fields
func validateEnum(config MorpheCompileConfig, enum yaml.Enum) error {
	if enum.Name == "" {
		return fmt.Errorf("enum name cannot be empty")
	}

	if len(enum.Entries) == 0 {
		return fmt.Errorf("enum must have at least one entry")
	}

	if config.MorpheConfig.MorpheEnumsConfig.Schema == "" {
		return fmt.Errorf("schema cannot be empty")
	}

	return nil
}

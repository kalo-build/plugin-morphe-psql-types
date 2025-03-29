package compile

import (
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/write"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

// MorpheStructureToPSQLTable creates a standard structures table according to the spec
func MorpheStructureToPSQLTable(config MorpheCompileConfig) (*psqldef.Table, error) {
	morpheConfig, configStartErr := triggerCompileMorpheStructureStart(config.StructureHooks, config.MorpheConfig)
	if configStartErr != nil {
		return nil, triggerCompileMorpheStructureFailure(config.StructureHooks, morpheConfig, configStartErr)
	}
	config.MorpheConfig = morpheConfig

	// If persistence is not enabled, return early
	if !morpheConfig.MorpheStructuresConfig.EnablePersistence {
		return nil, nil
	}

	// Create a fixed table definition based on the spec
	structureTable := createStandardStructureTable(morpheConfig.MorpheStructuresConfig)

	structureTable, structureTableErr := triggerCompileMorpheStructureSuccess(config.StructureHooks, structureTable)
	if structureTableErr != nil {
		return nil, triggerCompileMorpheStructureFailure(config.StructureHooks, morpheConfig, structureTableErr)
	}

	return structureTable, nil
}

// WriteStructureTableDefinition writes the structure table definition
func WriteStructureTableDefinition(hooks hook.WritePSQLTable, writer write.PSQLTableWriter, structureTable *psqldef.Table) (*psqldef.Table, []byte, error) {
	return WriteModelTableDefinition(hooks, writer, structureTable)
}

// createStandardStructureTable creates the standard structure table
func createStandardStructureTable(config cfg.MorpheStructuresConfig) *psqldef.Table {
	idType := psqldef.PSQLTypeSerial
	if config.UseBigSerial {
		idType = psqldef.PSQLTypeBigSerial
	}

	// Create columns
	columns := []psqldef.TableColumn{
		{
			Name:       "id",
			Type:       idType,
			NotNull:    false, // Changed to match ground truth format
			PrimaryKey: true,
		},
		{
			Name:    "\"type\"", // Quoted to match ground truth
			Type:    psqldef.PSQLTypeText,
			NotNull: true,
		},
		{
			Name:    "\"data\"", // Quoted to match ground truth
			Type:    psqldef.PSQLTypeJSONB,
			NotNull: true,
		},
		{
			Name:    "created_at",
			Type:    psqldef.PSQLTypeTimestampTZ,
			Default: "NOW()",
		},
		{
			Name:    "updated_at",
			Type:    psqldef.PSQLTypeTimestampTZ,
			Default: "NOW()",
		},
	}

	// Create indices
	indices := []psqldef.Index{
		{
			Name:     "idx_morphe_structures_type",
			Columns:  []string{"\"type\""}, // Quoted to match ground truth
			IsUnique: false,
		},
		{
			Name:     "idx_morphe_structures_data",
			Columns:  []string{"\"data\""}, // Quoted to match ground truth
			IsUnique: false,
			Using:    "GIN",
		},
	}

	// Create table
	return &psqldef.Table{
		Schema:  config.Schema,
		Name:    "morphe_structures",
		Columns: columns,
		Indices: indices,
	}
}

func triggerCompileMorpheStructureStart(hooks hook.CompileMorpheStructure, config cfg.MorpheConfig) (cfg.MorpheConfig, error) {
	if hooks.OnCompileMorpheStructureStart == nil {
		return config, nil
	}

	updatedConfig, startErr := hooks.OnCompileMorpheStructureStart(config)
	if startErr != nil {
		return cfg.MorpheConfig{}, startErr
	}

	return updatedConfig, nil
}

func triggerCompileMorpheStructureSuccess(hooks hook.CompileMorpheStructure, structureTable *psqldef.Table) (*psqldef.Table, error) {
	if hooks.OnCompileMorpheStructureSuccess == nil {
		return structureTable, nil
	}
	if structureTable == nil {
		return nil, ErrNoStructureTable
	}

	tableClone := structureTable.DeepClone()
	updatedTable, successErr := hooks.OnCompileMorpheStructureSuccess(&tableClone)
	if successErr != nil {
		return nil, successErr
	}
	return updatedTable, nil
}

func triggerCompileMorpheStructureFailure(hooks hook.CompileMorpheStructure, config cfg.MorpheConfig, failureErr error) error {
	if hooks.OnCompileMorpheStructureFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheStructureFailure(config, failureErr)
}

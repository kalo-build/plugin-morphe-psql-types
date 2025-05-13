package compile

import (
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/registry"
)

func MorpheToPSQL(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	if r.HasEnums() {
		allEnumTables, compileAllEnumsErr := AllMorpheEnumsToPSQLTables(config, r)
		if compileAllEnumsErr != nil {
			return compileAllEnumsErr
		}

		_, writeEnumTablesErr := WriteAllEnumTableDefinitions(config, allEnumTables)
		if writeEnumTablesErr != nil {
			return writeEnumTablesErr
		}
	}

	hasModels := r.HasModels()
	if hasModels {
		allModelTables, compileAllModelsErr := AllMorpheModelsToPSQLTables(config, r)
		if compileAllModelsErr != nil {
			return compileAllModelsErr
		}

		_, writeModelTablesErr := WriteAllModelTableDefinitions(config, allModelTables)
		if writeModelTablesErr != nil {
			return writeModelTablesErr
		}
	}

	if r.HasStructures() {
		if config.StructureWriter == nil {
			return ErrNoStructureWriter
		}

		structureTable, compileStructureErr := MorpheStructureToPSQLTable(config)
		if compileStructureErr != nil {
			return compileStructureErr
		}

		_, _, writeStructureErr := WriteStructureTableDefinition(config.WriteTableHooks, config.StructureWriter, structureTable)
		if writeStructureErr != nil {
			return writeStructureErr
		}
	}

	if r.HasEntities() {
		if !hasModels {
			return fmt.Errorf("entities compilation requires models to be compiled")
		}

		allEntityViews, compileAllEntityViewsErr := AllMorpheEntitiesToPSQLViews(config, r)
		if compileAllEntityViewsErr != nil {
			return compileAllEntityViewsErr
		}

		_, writeEntityViewsErr := WriteAllEntityViewDefinitions(config, allEntityViews)
		if writeEntityViewsErr != nil {
			return writeEntityViewsErr
		}
	}

	return nil
}

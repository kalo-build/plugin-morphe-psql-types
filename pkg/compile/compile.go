package compile

import "github.com/kaloseia/morphe-go/pkg/registry"

func MorpheToPSQL(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	allEnumTables, compileAllEnumsErr := AllMorpheEnumsToPSQLTables(config, r)
	if compileAllEnumsErr != nil {
		return compileAllEnumsErr
	}

	_, writeEnumTablesErr := WriteAllEnumTableDefinitions(config, allEnumTables)
	if writeEnumTablesErr != nil {
		return writeEnumTablesErr
	}

	allModelTables, compileAllModelsErr := AllMorpheModelsToPSQLTables(config, r)
	if compileAllModelsErr != nil {
		return compileAllModelsErr
	}

	_, writeModelTablesErr := WriteAllModelTableDefinitions(config, allModelTables)
	if writeModelTablesErr != nil {
		return writeModelTablesErr
	}

	return nil
}

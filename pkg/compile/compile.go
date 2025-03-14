package compile

import "github.com/kaloseia/morphe-go/pkg/registry"

func MorpheToPSQL(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	_, compileAllModelsErr := AllMorpheModelsToPSQLTables(config, r) // allModelTables
	if compileAllModelsErr != nil {
		return compileAllModelsErr
	}

	// _, writeModelStructsErr := WriteAllModelStructDefinitions(config, allModelStructDefs)
	// if writeModelStructsErr != nil {
	// 	return writeModelStructsErr
	// }

	return nil
}

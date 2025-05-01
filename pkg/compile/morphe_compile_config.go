package compile

import (
	"path"

	r "github.com/kalo-build/morphe-go/pkg/registry"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/write"
)

type MorpheCompileConfig struct {
	rcfg.MorpheLoadRegistryConfig
	cfg.MorpheConfig

	RegistryHooks r.LoadMorpheRegistryHooks

	ModelWriter write.PSQLTableWriter
	ModelHooks  hook.CompileMorpheModel

	EnumWriter write.PSQLTableWriter
	EnumHooks  hook.CompileMorpheEnum

	StructureWriter write.PSQLTableWriter
	StructureHooks  hook.CompileMorpheStructure

	EntityWriter write.PSQLViewWriter
	EntityHooks  hook.CompileMorpheEntity

	WriteTableHooks hook.WritePSQLTable
	WriteViewHooks  hook.WritePSQLView
}

func (config MorpheCompileConfig) Validate() error {
	loadRegistryErr := config.MorpheLoadRegistryConfig.Validate()
	if loadRegistryErr != nil {
		return loadRegistryErr
	}

	morpheCfgErr := config.MorpheConfig.Validate()
	if morpheCfgErr != nil {
		return morpheCfgErr
	}

	entitiesCfgErr := config.MorpheEntitiesConfig.Validate()
	if entitiesCfgErr != nil {
		return entitiesCfgErr
	}

	return nil
}

func DefaultMorpheCompileConfig(
	yamlRegistryPath string,
	baseOutputDirPath string,
) MorpheCompileConfig {
	return MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      path.Join(yamlRegistryPath, "enums"),
			RegistryModelsDirPath:     path.Join(yamlRegistryPath, "models"),
			RegistryStructuresDirPath: path.Join(yamlRegistryPath, "structures"),
			RegistryEntitiesDirPath:   path.Join(yamlRegistryPath, "entities"),
		},
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig: cfg.MorpheModelsConfig{
				Schema:       "public",
				UseBigSerial: false,
			},
			MorpheEnumsConfig: cfg.MorpheEnumsConfig{
				Schema:       "public",
				UseBigSerial: false,
			},
			MorpheStructuresConfig: cfg.MorpheStructuresConfig{
				Schema:            "public",
				UseBigSerial:      false,
				EnablePersistence: true,
			},
			MorpheEntitiesConfig: cfg.MorpheEntitiesConfig{
				Schema:         "public",
				ViewNameSuffix: "_entities",
			},
		},

		RegistryHooks: r.LoadMorpheRegistryHooks{},

		EnumWriter: &MorpheTableFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "enums"),
		},
		EnumHooks: hook.CompileMorpheEnum{},

		ModelWriter: &MorpheTableFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "models"),
		},
		ModelHooks: hook.CompileMorpheModel{},

		EntityWriter: &MorpheViewFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "entities"),
		},
		EntityHooks: hook.CompileMorpheEntity{},

		WriteTableHooks: hook.WritePSQLTable{},
		WriteViewHooks:  hook.WritePSQLView{},

		StructureWriter: &MorpheTableFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "structures"),
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}
}

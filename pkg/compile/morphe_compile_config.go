package compile

import (
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

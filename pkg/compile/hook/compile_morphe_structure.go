package hook

import (
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

type CompileMorpheStructure struct {
	OnCompileMorpheStructureStart   OnCompileMorpheStructureStartHook
	OnCompileMorpheStructureSuccess OnCompileMorpheStructureSuccessHook
	OnCompileMorpheStructureFailure OnCompileMorpheStructureFailureHook
}

type OnCompileMorpheStructureStartHook = func(config cfg.MorpheConfig) (cfg.MorpheConfig, error)
type OnCompileMorpheStructureSuccessHook = func(structureTable *psqldef.Table) (*psqldef.Table, error)
type OnCompileMorpheStructureFailureHook = func(config cfg.MorpheConfig, compileFailure error) error

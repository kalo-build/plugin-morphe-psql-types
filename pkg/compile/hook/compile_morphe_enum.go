package hook

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
)

type CompileMorpheEnum struct {
	OnCompileMorpheEnumStart   OnCompileMorpheEnumStartHook
	OnCompileMorpheEnumSuccess OnCompileMorpheEnumSuccessHook
	OnCompileMorpheEnumFailure OnCompileMorpheEnumFailureHook
}

type OnCompileMorpheEnumStartHook = func(config cfg.MorpheConfig, enum yaml.Enum) (cfg.MorpheConfig, yaml.Enum, error)
type OnCompileMorpheEnumSuccessHook = func(table *psqldef.Table, seedData *psqldef.InsertStatement) (*psqldef.Table, *psqldef.InsertStatement, error)
type OnCompileMorpheEnumFailureHook = func(config cfg.MorpheConfig, enum yaml.Enum, compileFailure error) error

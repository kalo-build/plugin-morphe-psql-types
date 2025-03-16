package hook

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
)

// CompileMorpheEnum contains hooks for compiling Morphe enums to PostgreSQL
type CompileMorpheEnum struct {
	// Called at the start of compilation for an enum
	OnCompileMorpheEnumStart OnCompileMorpheEnumStartHook

	// Called on successful compilation of an enum
	OnCompileMorpheEnumSuccess OnCompileMorpheEnumSuccessHook

	// Called when compilation of an enum fails
	OnCompileMorpheEnumFailure OnCompileMorpheEnumFailureHook
}

type OnCompileMorpheEnumStartHook = func(config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error)
type OnCompileMorpheEnumSuccessHook = func(table *psqldef.Table) (*psqldef.Table, error)
type OnCompileMorpheEnumFailureHook = func(config cfg.MorpheEnumsConfig, enum yaml.Enum, compileFailure error) error

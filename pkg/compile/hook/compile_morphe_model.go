package hook

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
)

type CompileMorpheModel struct {
	OnCompileMorpheModelStart   OnCompileMorpheModelStartHook
	OnCompileMorpheModelSuccess OnCompileMorpheModelSuccessHook
	OnCompileMorpheModelFailure OnCompileMorpheModelFailureHook
}

type OnCompileMorpheModelStartHook = func(config cfg.MorpheConfig, model yaml.Model) (cfg.MorpheConfig, yaml.Model, error)
type OnCompileMorpheModelSuccessHook = func(allModelTables []*psqldef.Table) ([]*psqldef.Table, error)
type OnCompileMorpheModelFailureHook = func(config cfg.MorpheConfig, model yaml.Model, compileFailure error) error

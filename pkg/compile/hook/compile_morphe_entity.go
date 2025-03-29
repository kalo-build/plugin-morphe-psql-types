package hook

import (
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

type CompileMorpheEntity struct {
	OnCompileMorpheEntityStart   OnCompileMorpheEntityStartHook
	OnCompileMorpheEntitySuccess OnCompileMorpheEntitySuccessHook
	OnCompileMorpheEntityFailure OnCompileMorpheEntityFailureHook
}

type OnCompileMorpheEntityStartHook = func(config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error)
type OnCompileMorpheEntitySuccessHook = func(view *psqldef.View) (*psqldef.View, error)
type OnCompileMorpheEntityFailureHook = func(config cfg.MorpheConfig, entity yaml.Entity, compileFailure error) error

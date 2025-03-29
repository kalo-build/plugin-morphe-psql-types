package compile

import (
	"fmt"

	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

// Error definitions
var (
	ErrMissingMorpheEntityField = func(entityName, fieldName string) error {
		return fmt.Errorf("missing entity field %s in entity %s", fieldName, entityName)
	}
)

// AllMorpheEntitiesToPSQLViews compiles all Morphe entities to PostgreSQL views
func AllMorpheEntitiesToPSQLViews(config MorpheCompileConfig, r *registry.Registry) (map[string]*psqldef.View, error) {
	allViews := map[string]*psqldef.View{}
	for entityName, entity := range r.GetAllEntities() {
		view, viewErr := MorpheEntityToPSQLView(config, r, entity)
		if viewErr != nil {
			return nil, viewErr
		}
		allViews[entityName] = view
	}
	return allViews, nil
}

// MorpheEntityToPSQLView compiles a single Morphe entity to a PostgreSQL view
func MorpheEntityToPSQLView(config MorpheCompileConfig, r *registry.Registry, entity yaml.Entity) (*psqldef.View, error) {
	if r == nil {
		return nil, triggerCompileMorpheEntityFailure(config.EntityHooks, config.MorpheConfig, entity, ErrNoRegistry)
	}

	morpheConfig, entity, compileStartErr := triggerCompileMorpheEntityStart(config.EntityHooks, config.MorpheConfig, entity)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheEntityFailure(config.EntityHooks, config.MorpheConfig, entity, compileStartErr)
	}
	config.MorpheConfig = morpheConfig

	view, viewErr := morpheEntityToPSQLView(config.MorpheConfig, r, entity)
	if viewErr != nil {
		return nil, triggerCompileMorpheEntityFailure(config.EntityHooks, config.MorpheConfig, entity, viewErr)
	}

	view, compileSuccessErr := triggerCompileMorpheEntitySuccess(config.EntityHooks, view)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheEntityFailure(config.EntityHooks, config.MorpheConfig, entity, compileSuccessErr)
	}

	return view, nil
}

func morpheEntityToPSQLView(config cfg.MorpheConfig, r *registry.Registry, entity yaml.Entity) (*psqldef.View, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}

	validateEntityErr := entity.Validate(r.GetAllModels(), r.GetAllEnums())
	if validateEntityErr != nil {
		return nil, validateEntityErr
	}

	viewName := strcase.ToSnakeCaseLower(entity.Name)
	if config.MorpheEntitiesConfig.ViewNameSuffix != "" {
		viewName += config.MorpheEntitiesConfig.ViewNameSuffix
	}

	view := &psqldef.View{
		Schema: config.MorpheEntitiesConfig.Schema,
		Name:   viewName,
	}

	return view, nil
}

func triggerCompileMorpheEntityStart(hooks hook.CompileMorpheEntity, config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error) {
	if hooks.OnCompileMorpheEntityStart == nil {
		return config, entity, nil
	}

	return hooks.OnCompileMorpheEntityStart(config, entity)
}

func triggerCompileMorpheEntitySuccess(hooks hook.CompileMorpheEntity, view *psqldef.View) (*psqldef.View, error) {
	if hooks.OnCompileMorpheEntitySuccess == nil {
		return view, nil
	}

	return hooks.OnCompileMorpheEntitySuccess(view)
}

func triggerCompileMorpheEntityFailure(hooks hook.CompileMorpheEntity, config cfg.MorpheConfig, entity yaml.Entity, failureErr error) error {
	if hooks.OnCompileMorpheEntityFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheEntityFailure(config, entity, failureErr)
}

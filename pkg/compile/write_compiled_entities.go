package compile

import (
	"github.com/kalo-build/clone"
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/write"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

func WriteAllEntityViewDefinitions(config MorpheCompileConfig, allEntityViewDefs map[string]*psqldef.View) (CompiledMorpheViews, error) {
	allWrittenEntities := CompiledMorpheViews{}

	sortedEntityNames := core.MapKeysSorted(allEntityViewDefs)
	for _, entityName := range sortedEntityNames {
		entityView := allEntityViewDefs[entityName]
		entityView, entityViewContents, writeErr := WriteEntityViewDefinition(config.WriteViewHooks, config.EntityWriter, entityView)
		if writeErr != nil {
			return nil, writeErr
		}
		allWrittenEntities.AddCompiledMorpheView(entityName, entityView, entityViewContents)
	}
	return allWrittenEntities, nil
}

func WriteEntityViewDefinition(hooks hook.WritePSQLView, writer write.PSQLViewWriter, entityView *psqldef.View) (*psqldef.View, []byte, error) {
	writer, entityView, writeStartErr := triggerWriteEntityViewStart(hooks, writer, entityView)
	if writeStartErr != nil {
		return nil, nil, triggerWriteEntityViewFailure(hooks, writer, entityView, writeStartErr)
	}

	entityViewContents, writeViewErr := writer.WriteView(entityView)
	if writeViewErr != nil {
		return nil, nil, triggerWriteEntityViewFailure(hooks, writer, entityView, writeViewErr)
	}

	entityView, entityViewContents, writeSuccessErr := triggerWriteEntityViewSuccess(hooks, entityView, entityViewContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteEntityViewFailure(hooks, writer, entityView, writeSuccessErr)
	}
	return entityView, entityViewContents, nil
}

func triggerWriteEntityViewStart(hooks hook.WritePSQLView, writer write.PSQLViewWriter, entityView *psqldef.View) (write.PSQLViewWriter, *psqldef.View, error) {
	if hooks.OnWritePSQLViewStart == nil {
		return writer, entityView, nil
	}
	if entityView == nil {
		return nil, nil, ErrNoEntityView
	}
	entityViewClone := entityView.DeepClone()

	updatedWriter, updatedEntityView, startErr := hooks.OnWritePSQLViewStart(writer, &entityViewClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedEntityView, nil
}

func triggerWriteEntityViewSuccess(hooks hook.WritePSQLView, entityView *psqldef.View, entityViewContents []byte) (*psqldef.View, []byte, error) {
	if hooks.OnWritePSQLViewSuccess == nil {
		return entityView, entityViewContents, nil
	}
	if entityView == nil {
		return nil, nil, ErrNoEntityView
	}
	entityViewClone := entityView.DeepClone()
	entityViewContentsClone := clone.Slice(entityViewContents)

	updatedEntityView, updatedEntityViewContents, successErr := hooks.OnWritePSQLViewSuccess(&entityViewClone, entityViewContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedEntityView, updatedEntityViewContents, nil
}

func triggerWriteEntityViewFailure(hooks hook.WritePSQLView, writer write.PSQLViewWriter, entityView *psqldef.View, failureErr error) error {
	if hooks.OnWritePSQLViewFailure == nil {
		return failureErr
	}

	entityViewClone := entityView.DeepClone()
	return hooks.OnWritePSQLViewFailure(writer, &entityViewClone, failureErr)
}

package compile

import (
	"github.com/kalo-build/clone"
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/write"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

func WriteAllModelTableDefinitions(config MorpheCompileConfig, allModelTableDefs map[string][]*psqldef.Table) (CompiledMorpheTables, error) {
	allWrittenModels := CompiledMorpheTables{}

	sortedModelNames := core.MapKeysSorted(allModelTableDefs)
	for _, modelName := range sortedModelNames {
		modelTables := allModelTableDefs[modelName]
		for _, modelTable := range modelTables {
			modelTable, modelTableContents, writeErr := WriteModelTableDefinition(config.WriteTableHooks, config.ModelWriter, modelTable)
			if writeErr != nil {
				return nil, writeErr
			}
			allWrittenModels.AddCompiledMorpheTable(modelName, modelTable, modelTableContents)
		}
	}
	return allWrittenModels, nil
}

func WriteModelTableDefinition(hooks hook.WritePSQLTable, writer write.PSQLTableWriter, modelTable *psqldef.Table) (*psqldef.Table, []byte, error) {
	writer, modelTable, writeStartErr := triggerWriteModelTableStart(hooks, writer, modelTable)
	if writeStartErr != nil {
		return nil, nil, triggerWriteModelTableFailure(hooks, writer, modelTable, writeStartErr)
	}

	modelTableContents, writeTableErr := writer.WriteTable(modelTable)
	if writeTableErr != nil {
		return nil, nil, triggerWriteModelTableFailure(hooks, writer, modelTable, writeTableErr)
	}

	modelTable, modelTableContents, writeSuccessErr := triggerWriteModelTableSuccess(hooks, modelTable, modelTableContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteModelTableFailure(hooks, writer, modelTable, writeSuccessErr)
	}
	return modelTable, modelTableContents, nil
}

func triggerWriteModelTableStart(hooks hook.WritePSQLTable, writer write.PSQLTableWriter, modelTable *psqldef.Table) (write.PSQLTableWriter, *psqldef.Table, error) {
	if hooks.OnWritePSQLTableStart == nil {
		return writer, modelTable, nil
	}
	if modelTable == nil {
		return nil, nil, ErrNoModelTable
	}
	modelTableClone := modelTable.DeepClone()

	updatedWriter, updatedModelTable, startErr := hooks.OnWritePSQLTableStart(writer, &modelTableClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedModelTable, nil
}

func triggerWriteModelTableSuccess(hooks hook.WritePSQLTable, modelTable *psqldef.Table, modelTableContents []byte) (*psqldef.Table, []byte, error) {
	if hooks.OnWritePSQLTableSuccess == nil {
		return modelTable, modelTableContents, nil
	}
	if modelTable == nil {
		return nil, nil, ErrNoModelTable
	}
	modelTableClone := modelTable.DeepClone()
	modelTableContentsClone := clone.Slice(modelTableContents)

	updatedModelTable, updatedModelTableContents, successErr := hooks.OnWritePSQLTableSuccess(&modelTableClone, modelTableContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedModelTable, updatedModelTableContents, nil
}

func triggerWriteModelTableFailure(hooks hook.WritePSQLTable, writer write.PSQLTableWriter, modelTable *psqldef.Table, failureErr error) error {
	if hooks.OnWritePSQLTableFailure == nil {
		return failureErr
	}

	modelTableClone := modelTable.DeepClone()
	return hooks.OnWritePSQLTableFailure(writer, &modelTableClone, failureErr)
}

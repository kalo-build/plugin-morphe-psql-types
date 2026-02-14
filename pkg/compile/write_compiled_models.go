package compile

import (
	"github.com/kalo-build/clone"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/write"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

// WriteAllModelTableDefinitions writes all model tables with dependency-based ordering but without order prefixes in filenames.
func WriteAllModelTableDefinitions(config MorpheCompileConfig, allModelTableDefs map[string][]*psqldef.Table) (CompiledMorpheTables, error) {
	allWrittenModels := CompiledMorpheTables{}

	// Flatten all tables for dependency sorting
	allTables := []*psqldef.Table{}
	tableToModel := make(map[string]string)
	for modelName, modelTables := range allModelTableDefs {
		for _, table := range modelTables {
			allTables = append(allTables, table)
			tableToModel[table.Name] = modelName
		}
	}

	// Sort tables by dependency order
	sortedTables, sortErr := SortTablesByDependency(allTables)
	if sortErr != nil {
		return nil, sortErr
	}

	// Write tables in dependency order without order prefix
	for _, modelTable := range sortedTables {
		modelName := tableToModel[modelTable.Name]

		modelTable, modelTableContents, writeErr := WriteModelTableDefinition(
			config.WriteTableHooks, config.ModelWriter, modelTable)
		if writeErr != nil {
			return nil, writeErr
		}
		allWrittenModels.AddCompiledMorpheTable(modelName, modelTable, modelTableContents)
	}

	return allWrittenModels, nil
}

// WriteAllModelTableDefinitionsWithOrder writes all model tables with dependency-based ordering.
// The startOrder parameter is the starting order number for file prefixes.
// Returns the compiled tables and the next order number to use.
func WriteAllModelTableDefinitionsWithOrder(config MorpheCompileConfig, allModelTableDefs map[string][]*psqldef.Table, startOrder int) (CompiledMorpheTables, error) {
	allWrittenModels := CompiledMorpheTables{}

	// Flatten all tables for dependency sorting
	allTables := []*psqldef.Table{}
	tableToModel := make(map[string]string)
	for modelName, modelTables := range allModelTableDefs {
		for _, table := range modelTables {
			allTables = append(allTables, table)
			tableToModel[table.Name] = modelName
		}
	}

	// Sort tables by dependency order
	sortedTables, sortErr := SortTablesByDependency(allTables)
	if sortErr != nil {
		return nil, sortErr
	}

	// Write tables in dependency order with incrementing order prefix
	currentOrder := startOrder
	for _, modelTable := range sortedTables {
		currentOrder++
		modelName := tableToModel[modelTable.Name]

		modelTable, modelTableContents, writeErr := WriteModelTableDefinitionWithOrder(
			config.WriteTableHooks, config.ModelWriter, modelTable, currentOrder)
		if writeErr != nil {
			return nil, writeErr
		}
		allWrittenModels.AddCompiledMorpheTable(modelName, modelTable, modelTableContents)
	}

	return allWrittenModels, nil
}

func WriteModelTableDefinition(hooks hook.WritePSQLTable, writer write.PSQLTableWriter, modelTable *psqldef.Table) (*psqldef.Table, []byte, error) {
	return WriteModelTableDefinitionWithOrder(hooks, writer, modelTable, 0)
}

func WriteModelTableDefinitionWithOrder(hooks hook.WritePSQLTable, writer write.PSQLTableWriter, modelTable *psqldef.Table, order int) (*psqldef.Table, []byte, error) {
	writer, modelTable, writeStartErr := triggerWriteModelTableStart(hooks, writer, modelTable)
	if writeStartErr != nil {
		return nil, nil, triggerWriteModelTableFailure(hooks, writer, modelTable, writeStartErr)
	}

	var modelTableContents []byte
	var writeTableErr error

	// Check if writer supports ordered writing
	if orderedWriter, ok := writer.(write.OrderedPSQLTableWriter); ok && order > 0 {
		modelTableContents, writeTableErr = orderedWriter.WriteTableWithOrder(modelTable, order)
	} else {
		modelTableContents, writeTableErr = writer.WriteTable(modelTable)
	}

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

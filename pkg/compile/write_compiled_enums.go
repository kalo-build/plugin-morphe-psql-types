package compile

import (
	"github.com/kalo-build/clone"
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/write"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

func WriteAllEnumTableDefinitions(config MorpheCompileConfig, allEnumTableDefs map[string]*psqldef.Table) (CompiledMorpheTables, error) {
	allWrittenEnums := CompiledMorpheTables{}

	sortedEnumNames := core.MapKeysSorted(allEnumTableDefs)
	for _, enumName := range sortedEnumNames {
		enumTable := allEnumTableDefs[enumName]
		enumTable, enumTableContents, writeErr := WriteEnumTableDefinition(config.WriteTableHooks, config.EnumWriter, enumTable)
		if writeErr != nil {
			return nil, writeErr
		}
		allWrittenEnums.AddCompiledMorpheTable(enumName, enumTable, enumTableContents)
	}
	return allWrittenEnums, nil
}

func WriteEnumTableDefinition(hooks hook.WritePSQLTable, writer write.PSQLTableWriter, enumTable *psqldef.Table) (*psqldef.Table, []byte, error) {
	writer, enumTable, writeStartErr := triggerWriteEnumTableStart(hooks, writer, enumTable)
	if writeStartErr != nil {
		return nil, nil, triggerWriteEnumTableFailure(hooks, writer, enumTable, writeStartErr)
	}

	enumTableContents, writeTableErr := writer.WriteTable(enumTable)
	if writeTableErr != nil {
		return nil, nil, triggerWriteEnumTableFailure(hooks, writer, enumTable, writeTableErr)
	}

	enumTable, enumTableContents, writeSuccessErr := triggerWriteEnumTableSuccess(hooks, enumTable, enumTableContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteEnumTableFailure(hooks, writer, enumTable, writeSuccessErr)
	}
	return enumTable, enumTableContents, nil
}

func triggerWriteEnumTableStart(hooks hook.WritePSQLTable, writer write.PSQLTableWriter, enumTable *psqldef.Table) (write.PSQLTableWriter, *psqldef.Table, error) {
	if hooks.OnWritePSQLTableStart == nil {
		return writer, enumTable, nil
	}
	if enumTable == nil {
		return nil, nil, ErrNoEnumTable
	}
	enumTableClone := enumTable.DeepClone()

	updatedWriter, updatedEnumTable, startErr := hooks.OnWritePSQLTableStart(writer, &enumTableClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedEnumTable, nil
}

func triggerWriteEnumTableSuccess(hooks hook.WritePSQLTable, enumTable *psqldef.Table, enumTableContents []byte) (*psqldef.Table, []byte, error) {
	if hooks.OnWritePSQLTableSuccess == nil {
		return enumTable, enumTableContents, nil
	}
	if enumTable == nil {
		return nil, nil, ErrNoEnumTable
	}
	enumTableClone := enumTable.DeepClone()
	enumTableContentsClone := clone.Slice(enumTableContents)

	updatedEnumTable, updatedEnumTableContents, successErr := hooks.OnWritePSQLTableSuccess(&enumTableClone, enumTableContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedEnumTable, updatedEnumTableContents, nil
}

func triggerWriteEnumTableFailure(hooks hook.WritePSQLTable, writer write.PSQLTableWriter, enumTable *psqldef.Table, failureErr error) error {
	if hooks.OnWritePSQLTableFailure == nil {
		return failureErr
	}

	enumTableClone := enumTable.DeepClone()
	return hooks.OnWritePSQLTableFailure(writer, &enumTableClone, failureErr)
}

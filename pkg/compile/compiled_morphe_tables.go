package compile

import "github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"

// CompiledMorpheTables maps Morphe.Name -> MorpheTable.Name -> CompiledTable
type CompiledMorpheTables map[string]map[string]CompiledTable

func (tables CompiledMorpheTables) AddCompiledMorpheTable(morpheName string, tableDef *psqldef.Table, tableContents []byte) {
	if tables[morpheName] == nil {
		tables[morpheName] = make(map[string]CompiledTable)
	}
	tables[morpheName][tableDef.Name] = CompiledTable{
		Table:         tableDef,
		TableContents: tableContents,
	}
}

func (tables CompiledMorpheTables) GetAllCompiledMorpheTables(morpheName string) map[string]CompiledTable {
	morpheTables, morpheTablesExist := tables[morpheName]
	if !morpheTablesExist {
		return nil
	}
	return morpheTables
}

func (tables CompiledMorpheTables) GetCompiledMorpheTable(morpheName string, tableName string) CompiledTable {
	morpheTables, morpheTablesExist := tables[morpheName]
	if !morpheTablesExist {
		return CompiledTable{}
	}
	compiledTable, compiledTableExists := morpheTables[tableName]
	if !compiledTableExists {
		return CompiledTable{}
	}
	return compiledTable
}

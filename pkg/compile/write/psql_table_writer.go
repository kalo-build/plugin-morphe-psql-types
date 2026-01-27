package write

import "github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"

type PSQLTableWriter interface {
	WriteTable(*psqldef.Table) ([]byte, error)
}

// OrderedPSQLTableWriter extends PSQLTableWriter with order support for migration files.
type OrderedPSQLTableWriter interface {
	PSQLTableWriter
	// WriteTableWithOrder writes a table with an order prefix (e.g., "001_table.sql")
	WriteTableWithOrder(*psqldef.Table, int) ([]byte, error)
}

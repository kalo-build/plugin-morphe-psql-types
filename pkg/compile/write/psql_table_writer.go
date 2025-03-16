package write

import "github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"

type PSQLTableWriter interface {
	WriteTable(*psqldef.Table) ([]byte, error)
}

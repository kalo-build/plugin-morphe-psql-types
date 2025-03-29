package write

import "github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"

type PSQLTableWriter interface {
	WriteTable(*psqldef.Table) ([]byte, error)
}

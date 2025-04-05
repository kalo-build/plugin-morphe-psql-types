package write

import "github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"

type PSQLViewWriter interface {
	WriteView(*psqldef.View) ([]byte, error)
}

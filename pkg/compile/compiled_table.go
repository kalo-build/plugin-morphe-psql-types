package compile

import "github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"

type CompiledTable struct {
	Table         *psqldef.Table
	TableContents []byte
}

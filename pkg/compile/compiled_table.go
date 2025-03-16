package compile

import "github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"

type CompiledTable struct {
	Table         *psqldef.Table
	TableContents []byte
}

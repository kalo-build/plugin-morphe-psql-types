package compile

import "github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"

type CompiledView struct {
	View         *psqldef.View
	ViewContents []byte
}

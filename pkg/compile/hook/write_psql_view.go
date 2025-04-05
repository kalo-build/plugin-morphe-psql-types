package hook

import (
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/write"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

type WritePSQLView struct {
	OnWritePSQLViewStart   OnWritePSQLViewStartHook
	OnWritePSQLViewSuccess OnWritePSQLViewSuccessHook
	OnWritePSQLViewFailure OnWritePSQLViewFailureHook
}

type OnWritePSQLViewStartHook = func(writer write.PSQLViewWriter, view *psqldef.View) (write.PSQLViewWriter, *psqldef.View, error)
type OnWritePSQLViewSuccessHook = func(view *psqldef.View, viewContents []byte) (*psqldef.View, []byte, error)
type OnWritePSQLViewFailureHook = func(writer write.PSQLViewWriter, view *psqldef.View, failureErr error) error

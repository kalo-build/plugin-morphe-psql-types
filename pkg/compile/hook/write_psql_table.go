package hook

import (
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/write"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
)

type WritePSQLTable struct {
	OnWritePSQLTableStart   OnWritePSQLTableStartHook
	OnWritePSQLTableSuccess OnWritePSQLTableSuccessHook
	OnWritePSQLTableFailure OnWritePSQLTableFailureHook
}

type OnWritePSQLTableStartHook = func(writer write.PSQLTableWriter, table *psqldef.Table) (write.PSQLTableWriter, *psqldef.Table, error)
type OnWritePSQLTableSuccessHook = func(table *psqldef.Table, tableContents []byte) (*psqldef.Table, []byte, error)
type OnWritePSQLTableFailureHook = func(writer write.PSQLTableWriter, table *psqldef.Table, failureErr error) error

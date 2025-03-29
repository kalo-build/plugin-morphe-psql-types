package psqldef

import "github.com/kalo-build/clone"

// InsertStatement represents a PSQL INSERT statement
type InsertStatement struct {
	Schema    string
	TableName string
	Columns   []string
	Values    [][]any
}

// DeepClone creates a deep copy of the InsertStatement
func (i InsertStatement) DeepClone() InsertStatement {
	insertCopy := InsertStatement{
		Schema:    i.Schema,
		TableName: i.TableName,
	}

	insertCopy.Columns = clone.Slice(i.Columns)

	if i.Values != nil {
		insertCopy.Values = make([][]any, len(i.Values))
		for rowIdx, row := range i.Values {
			insertCopy.Values[rowIdx] = clone.Slice(row)
		}
	}

	return insertCopy
}

package psqldef

import "github.com/kaloseia/clone"

// Table represents a PSQL table for a model
type Table struct {
	Schema            string
	Name              string
	Columns           []Column
	Indices           []Index
	ForeignKeys       []ForeignKey
	UniqueConstraints []UniqueConstraint
}

// DeepClone creates a deep copy of the Table
func (t Table) DeepClone() Table {
	tableCopy := Table{
		Schema:            t.Schema,
		Name:              t.Name,
		Columns:           clone.DeepCloneSlice(t.Columns),
		Indices:           clone.DeepCloneSlice(t.Indices),
		ForeignKeys:       clone.DeepCloneSlice(t.ForeignKeys),
		UniqueConstraints: clone.DeepCloneSlice(t.UniqueConstraints),
	}

	return tableCopy
}

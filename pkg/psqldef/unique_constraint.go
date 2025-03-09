package psqldef

import "github.com/kaloseia/clone"

// UniqueConstraint represents a unique constraint in a PSQL table
type UniqueConstraint struct {
	Schema      string
	Name        string
	TableName   string
	ColumnNames []string
}

// DeepClone creates a deep copy of the UniqueConstraint
func (c UniqueConstraint) DeepClone() UniqueConstraint {
	constraintCopy := UniqueConstraint{
		Schema:      c.Schema,
		Name:        c.Name,
		TableName:   c.TableName,
		ColumnNames: clone.Slice(c.ColumnNames),
	}

	return constraintCopy
}

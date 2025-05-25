package psqldef

import "github.com/kalo-build/clone"

// Index represents an index in a PSQL table
type Index struct {
	Name      string
	TableName string
	Columns   []string
	IsUnique  bool
	Using     string // e.g., "btree", "gin"
}

// DeepClone creates a deep copy of the Index
func (i Index) DeepClone() Index {
	indexCopy := Index{
		Name:      i.Name,
		TableName: i.TableName,
		Columns:   clone.Slice(i.Columns),
		IsUnique:  i.IsUnique,
		Using:     i.Using,
	}

	return indexCopy
}

package psqldef

import "github.com/kalo-build/clone"

// ForeignKey represents a foreign key in a PSQL table
type ForeignKey struct {
	Schema         string
	Name           string
	TableName      string
	ColumnNames    []string
	RefSchema      string
	RefTableName   string
	RefColumnNames []string
	OnDelete       string // e.g., "CASCADE", "SET NULL"
	OnUpdate       string // e.g., "CASCADE", "SET NULL"
}

// DeepClone creates a deep copy of the ForeignKey
func (fk ForeignKey) DeepClone() ForeignKey {
	foreignKeyCopy := ForeignKey{
		Schema:         fk.Schema,
		Name:           fk.Name,
		TableName:      fk.TableName,
		ColumnNames:    clone.Slice(fk.ColumnNames),
		RefSchema:      fk.RefSchema,
		RefTableName:   fk.RefTableName,
		RefColumnNames: clone.Slice(fk.RefColumnNames),
		OnDelete:       fk.OnDelete,
		OnUpdate:       fk.OnUpdate,
	}

	return foreignKeyCopy
}

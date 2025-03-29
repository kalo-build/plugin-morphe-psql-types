package psqldef

// TableColumn represents a column in a PSQL table
type TableColumn struct {
	Name       string
	Type       PSQLType
	NotNull    bool
	PrimaryKey bool
	Default    string
}

// DeepClone creates a deep copy of the TableColumn
func (c TableColumn) DeepClone() TableColumn {
	columnCopy := TableColumn{
		Name:       c.Name,
		Type:       DeepClonePSQLType(c.Type),
		NotNull:    c.NotNull,
		PrimaryKey: c.PrimaryKey,
		Default:    c.Default,
	}

	return columnCopy
}

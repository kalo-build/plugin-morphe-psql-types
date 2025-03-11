package psqldef

// Column represents a column in a PSQL table
type Column struct {
	Name       string
	Type       PSQLType
	NotNull    bool
	PrimaryKey bool
	Default    string
}

// DeepClone creates a deep copy of the Column
func (c Column) DeepClone() Column {
	columnCopy := Column{
		Name:       c.Name,
		Type:       DeepClonePSQLType(c.Type),
		NotNull:    c.NotNull,
		PrimaryKey: c.PrimaryKey,
		Default:    c.Default,
	}

	return columnCopy
}

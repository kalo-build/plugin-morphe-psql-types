package psqldef

// ViewColumn represents a column in a PostgreSQL view
type ViewColumn struct {
	Name      string
	SourceRef string // Format: "table_alias.column_name"
	Alias     string // Optional, if different from Name
}

// DeepClone creates a deep copy of the ViewColumn
func (vc ViewColumn) DeepClone() ViewColumn {
	return ViewColumn{
		Name:      vc.Name,
		SourceRef: vc.SourceRef,
		Alias:     vc.Alias,
	}
}

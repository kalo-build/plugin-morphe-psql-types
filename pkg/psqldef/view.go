package psqldef

import "github.com/kalo-build/clone"

// View represents a PostgreSQL view generated from a Morphe entity
type View struct {
	Schema      string
	Name        string
	Columns     []ViewColumn
	FromSchema  string
	FromTable   string
	Joins       []JoinClause
	WhereClause string
}

// DeepClone creates a deep copy of the View
func (v View) DeepClone() View {
	viewCopy := View{
		Schema:      v.Schema,
		Name:        v.Name,
		Columns:     clone.DeepCloneSlice(v.Columns),
		FromSchema:  v.FromSchema,
		FromTable:   v.FromTable,
		Joins:       clone.DeepCloneSlice(v.Joins),
		WhereClause: v.WhereClause,
	}

	return viewCopy
}

package psqldef

import "github.com/kalo-build/clone"

// JoinClause represents a JOIN clause in a SQL query
type JoinClause struct {
	Type       string // "LEFT", "INNER", etc.
	Table      string
	Alias      string
	Conditions []JoinCondition
}

// DeepClone creates a deep copy of the JoinClause
func (jc JoinClause) DeepClone() JoinClause {
	return JoinClause{
		Type:       jc.Type,
		Table:      jc.Table,
		Alias:      jc.Alias,
		Conditions: clone.DeepCloneSlice(jc.Conditions),
	}
}

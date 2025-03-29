package psqldef

// JoinCondition represents a condition in a JOIN clause
type JoinCondition struct {
	LeftRef  string // Format: "table_alias.column_name"
	RightRef string // Format: "table_alias.column_name"
}

// DeepClone creates a deep copy of the JoinCondition
func (jc JoinCondition) DeepClone() JoinCondition {
	return JoinCondition{
		LeftRef:  jc.LeftRef,
		RightRef: jc.RightRef,
	}
}

package psqldef

type PSQLTypeRange struct {
	ValueType PSQLType
	Schema    string
	Name      string
}

func (t PSQLTypeRange) IsPrimitive() bool {
	return false
}

func (t PSQLTypeRange) IsArray() bool {
	return false
}

func (t PSQLTypeRange) IsDomain() bool {
	return false
}

func (t PSQLTypeRange) IsComposite() bool {
	return false
}

func (t PSQLTypeRange) IsEnum() bool {
	return false
}

func (t PSQLTypeRange) IsRange() bool {
	return true
}

func (t PSQLTypeRange) GetSchema() string {
	return t.Schema
}

func (t PSQLTypeRange) GetSyntaxLocal() string {
	return t.Name
}

func (t PSQLTypeRange) GetSyntax() string {
	if t.Schema != "" {
		return t.Schema + "." + t.Name
	}
	return t.Name
}

func (t PSQLTypeRange) DeepClone() PSQLTypeRange {
	return PSQLTypeRange{
		ValueType: DeepClonePSQLType(t.ValueType),
		Schema:    t.Schema,
		Name:      t.Name,
	}
}

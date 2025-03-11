package psqldef

type PSQLTypeDomain struct {
	// ValueType supports only primitive, array, and domain types. No composite, enum, range types supported.
	ValueType PSQLType
	Schema    string
	Name      string
}

func (t PSQLTypeDomain) IsPrimitive() bool {
	return false
}

func (t PSQLTypeDomain) IsArray() bool {
	return false
}

func (t PSQLTypeDomain) IsDomain() bool {
	return true
}

func (t PSQLTypeDomain) IsComposite() bool {
	return false
}

func (t PSQLTypeDomain) IsEnum() bool {
	return false
}

func (t PSQLTypeDomain) IsRange() bool {
	return false
}

func (t PSQLTypeDomain) GetSchema() string {
	return t.Schema
}

func (t PSQLTypeDomain) GetSyntaxLocal() string {
	return t.Name
}

func (t PSQLTypeDomain) GetSyntax() string {
	if t.Schema != "" {
		return t.Schema + "." + t.Name
	}
	return t.Name
}

func (t PSQLTypeDomain) DeepClone() PSQLTypeDomain {
	return PSQLTypeDomain{
		ValueType: DeepClonePSQLType(t.ValueType),
		Schema:    t.Schema,
		Name:      t.Name,
	}
}

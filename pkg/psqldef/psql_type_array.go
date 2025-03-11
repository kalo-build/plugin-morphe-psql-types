package psqldef

type PSQLTypeArray struct {
	ValueType PSQLType
}

func (t PSQLTypeArray) IsPrimitive() bool {
	return false
}

func (t PSQLTypeArray) IsArray() bool {
	return true
}

func (t PSQLTypeArray) IsDomain() bool {
	return false
}

func (t PSQLTypeArray) IsComposite() bool {
	return false
}

func (t PSQLTypeArray) IsEnum() bool {
	return false
}

func (t PSQLTypeArray) IsRange() bool {
	return false
}

func (t PSQLTypeArray) GetSchema() string {
	return t.ValueType.GetSchema()
}

func (t PSQLTypeArray) GetSyntaxLocal() string {
	return t.ValueType.GetSyntaxLocal() + "[]"
}

func (t PSQLTypeArray) GetSyntax() string {
	return t.ValueType.GetSyntax() + "[]"
}

func (t PSQLTypeArray) DeepClone() PSQLTypeArray {
	return PSQLTypeArray{
		ValueType: DeepClonePSQLType(t.ValueType),
	}
}

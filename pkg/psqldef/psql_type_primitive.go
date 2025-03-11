package psqldef

type PSQLTypePrimitive struct {
	Syntax string
}

func (t PSQLTypePrimitive) IsPrimitive() bool {
	return true
}

func (t PSQLTypePrimitive) IsArray() bool {
	return false
}

func (t PSQLTypePrimitive) IsDomain() bool {
	return false
}

func (t PSQLTypePrimitive) IsComposite() bool {
	return false
}

func (t PSQLTypePrimitive) IsEnum() bool {
	return false
}

func (t PSQLTypePrimitive) IsRange() bool {
	return false
}

func (t PSQLTypePrimitive) GetSchema() string {
	return ""
}

func (t PSQLTypePrimitive) GetSyntaxLocal() string {
	return t.Syntax
}

func (t PSQLTypePrimitive) GetSyntax() string {
	return t.Syntax
}

func (t PSQLTypePrimitive) DeepClone() PSQLTypePrimitive {
	return PSQLTypePrimitive{
		Syntax: t.Syntax,
	}
}

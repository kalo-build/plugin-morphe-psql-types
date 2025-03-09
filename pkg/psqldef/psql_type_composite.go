package psqldef

type PSQLTypeComposite struct {
	Fields map[string]PSQLType
	Schema string
	Name   string
}

func (t PSQLTypeComposite) IsPrimitive() bool {
	return false
}

func (t PSQLTypeComposite) IsArray() bool {
	return false
}

func (t PSQLTypeComposite) IsDomain() bool {
	return false
}

func (t PSQLTypeComposite) IsComposite() bool {
	return true
}

func (t PSQLTypeComposite) IsEnum() bool {
	return false
}

func (t PSQLTypeComposite) IsRange() bool {
	return false
}

func (t PSQLTypeComposite) GetSchema() string {
	return t.Schema
}

func (t PSQLTypeComposite) GetSyntaxLocal() string {
	return t.Name
}

func (t PSQLTypeComposite) GetSyntax() string {
	if t.Schema != "" {
		return t.Schema + "." + t.Name
	}
	return t.Name
}

func (t PSQLTypeComposite) DeepClone() PSQLTypeComposite {
	return PSQLTypeComposite{
		Fields: DeepClonePSQLTypeMap(t.Fields),
		Schema: t.Schema,
		Name:   t.Name,
	}
}

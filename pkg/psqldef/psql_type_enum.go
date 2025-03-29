package psqldef

import "github.com/kalo-build/clone"

type PSQLTypeEnum struct {
	Values []string
	Schema string
	Name   string
}

func (t PSQLTypeEnum) IsPrimitive() bool {
	return false
}

func (t PSQLTypeEnum) IsArray() bool {
	return false
}

func (t PSQLTypeEnum) IsDomain() bool {
	return false
}

func (t PSQLTypeEnum) IsComposite() bool {
	return false
}

func (t PSQLTypeEnum) IsEnum() bool {
	return true
}

func (t PSQLTypeEnum) IsRange() bool {
	return false
}

func (t PSQLTypeEnum) GetSchema() string {
	return t.Schema
}

func (t PSQLTypeEnum) GetSyntaxLocal() string {
	return t.Name
}

func (t PSQLTypeEnum) GetSyntax() string {
	if t.Schema != "" {
		return t.Schema + "." + t.Name
	}
	return t.Name
}

func (t PSQLTypeEnum) DeepClone() PSQLTypeEnum {
	return PSQLTypeEnum{
		Values: clone.Slice(t.Values),
		Schema: t.Schema,
		Name:   t.Name,
	}
}

package psqldef

type PSQLType interface {
	IsPrimitive() bool
	IsArray() bool
	IsDomain() bool
	IsComposite() bool
	IsEnum() bool
	IsRange() bool

	GetSchema() string
	GetSyntax() string
	GetSyntaxLocal() string
}

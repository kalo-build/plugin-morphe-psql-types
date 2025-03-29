package typemap

import (
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

var MorpheEnumEntryToPSQLEntryType = map[yaml.EnumType]psqldef.PSQLType{
	yaml.EnumTypeString:  psqldef.PSQLTypeText,
	yaml.EnumTypeInteger: psqldef.PSQLTypeInteger,
	yaml.EnumTypeFloat:   psqldef.PSQLTypeDoublePrecision,
}

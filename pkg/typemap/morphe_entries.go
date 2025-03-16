package typemap

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
)

var MorpheEnumEntryToPSQLEntryType = map[yaml.EnumType]psqldef.PSQLType{
	yaml.EnumTypeString:  psqldef.PSQLTypeText,
	yaml.EnumTypeInteger: psqldef.PSQLTypeInteger,
	yaml.EnumTypeFloat:   psqldef.PSQLTypeDoublePrecision,
}

package typemap

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
)

var MorpheModelFieldToGoField = map[yaml.ModelFieldType]psqldef.PSQLType{
	yaml.ModelFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.ModelFieldTypeAutoIncrement: psqldef.PSQLTypeSerial,
	yaml.ModelFieldTypeString:        psqldef.PSQLTypeText,
	yaml.ModelFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.ModelFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.ModelFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.ModelFieldTypeTime:          psqldef.PSQLTypeTime,
	yaml.ModelFieldTypeDate:          psqldef.PSQLTypeDate,
	yaml.ModelFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.ModelFieldTypeSealed:        psqldef.PSQLTypeText,
}

var MorpheStructureFieldToGoField = map[yaml.StructureFieldType]psqldef.PSQLType{
	yaml.StructureFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.StructureFieldTypeAutoIncrement: psqldef.PSQLTypeSerial,
	yaml.StructureFieldTypeString:        psqldef.PSQLTypeText,
	yaml.StructureFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.StructureFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.StructureFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.StructureFieldTypeTime:          psqldef.PSQLTypeDate,
	yaml.StructureFieldTypeDate:          psqldef.PSQLTypeTime,
	yaml.StructureFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.StructureFieldTypeSealed:        psqldef.PSQLTypeText,
}

package typemap

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
)

var MorpheModelFieldToPSQLField = map[yaml.ModelFieldType]psqldef.PSQLType{
	yaml.ModelFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.ModelFieldTypeAutoIncrement: psqldef.PSQLTypeSerial,
	yaml.ModelFieldTypeString:        psqldef.PSQLTypeText,
	yaml.ModelFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.ModelFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.ModelFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.ModelFieldTypeTime:          psqldef.PSQLTypeTimestampTZ,
	yaml.ModelFieldTypeDate:          psqldef.PSQLTypeDate,
	yaml.ModelFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.ModelFieldTypeSealed:        psqldef.PSQLTypeText,
}

var MorpheStructureFieldToPSQLField = map[yaml.StructureFieldType]psqldef.PSQLType{
	yaml.StructureFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.StructureFieldTypeAutoIncrement: psqldef.PSQLTypeSerial,
	yaml.StructureFieldTypeString:        psqldef.PSQLTypeText,
	yaml.StructureFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.StructureFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.StructureFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.StructureFieldTypeTime:          psqldef.PSQLTypeTimestampTZ,
	yaml.StructureFieldTypeDate:          psqldef.PSQLTypeDate,
	yaml.StructureFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.StructureFieldTypeSealed:        psqldef.PSQLTypeText,
}

var MorpheModelFieldToPSQLFieldBigSerial = map[yaml.ModelFieldType]psqldef.PSQLType{
	yaml.ModelFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.ModelFieldTypeAutoIncrement: psqldef.PSQLTypeBigSerial,
	yaml.ModelFieldTypeString:        psqldef.PSQLTypeText,
	yaml.ModelFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.ModelFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.ModelFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.ModelFieldTypeTime:          psqldef.PSQLTypeTimestampTZ,
	yaml.ModelFieldTypeDate:          psqldef.PSQLTypeDate,
	yaml.ModelFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.ModelFieldTypeSealed:        psqldef.PSQLTypeText,
}

var MorpheStructureFieldToPSQLFieldBigSerial = map[yaml.StructureFieldType]psqldef.PSQLType{
	yaml.StructureFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.StructureFieldTypeAutoIncrement: psqldef.PSQLTypeBigSerial,
	yaml.StructureFieldTypeString:        psqldef.PSQLTypeText,
	yaml.StructureFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.StructureFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.StructureFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.StructureFieldTypeTime:          psqldef.PSQLTypeTimestampTZ,
	yaml.StructureFieldTypeDate:          psqldef.PSQLTypeDate,
	yaml.StructureFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.StructureFieldTypeSealed:        psqldef.PSQLTypeText,
}

var MorpheModelFieldToPSQLFieldForeign = map[yaml.ModelFieldType]psqldef.PSQLType{
	yaml.ModelFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.ModelFieldTypeAutoIncrement: psqldef.PSQLTypeInteger,
	yaml.ModelFieldTypeString:        psqldef.PSQLTypeText,
	yaml.ModelFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.ModelFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.ModelFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.ModelFieldTypeTime:          psqldef.PSQLTypeTimestampTZ,
	yaml.ModelFieldTypeDate:          psqldef.PSQLTypeDate,
	yaml.ModelFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.ModelFieldTypeSealed:        psqldef.PSQLTypeText,
}

var MorpheStructureFieldToPSQLFieldForeign = map[yaml.StructureFieldType]psqldef.PSQLType{
	yaml.StructureFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.StructureFieldTypeAutoIncrement: psqldef.PSQLTypeInteger,
	yaml.StructureFieldTypeString:        psqldef.PSQLTypeText,
	yaml.StructureFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.StructureFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.StructureFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.StructureFieldTypeTime:          psqldef.PSQLTypeTimestampTZ,
	yaml.StructureFieldTypeDate:          psqldef.PSQLTypeDate,
	yaml.StructureFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.StructureFieldTypeSealed:        psqldef.PSQLTypeText,
}

var MorpheModelFieldToPSQLFieldBigSerialForeign = map[yaml.ModelFieldType]psqldef.PSQLType{
	yaml.ModelFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.ModelFieldTypeAutoIncrement: psqldef.PSQLTypeBigInt,
	yaml.ModelFieldTypeString:        psqldef.PSQLTypeText,
	yaml.ModelFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.ModelFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.ModelFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.ModelFieldTypeTime:          psqldef.PSQLTypeTimestampTZ,
	yaml.ModelFieldTypeDate:          psqldef.PSQLTypeDate,
	yaml.ModelFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.ModelFieldTypeSealed:        psqldef.PSQLTypeText,
}

var MorpheStructureFieldToPSQLFieldBigSerialForeign = map[yaml.StructureFieldType]psqldef.PSQLType{
	yaml.StructureFieldTypeUUID:          psqldef.PSQLTypeUUID,
	yaml.StructureFieldTypeAutoIncrement: psqldef.PSQLTypeBigInt,
	yaml.StructureFieldTypeString:        psqldef.PSQLTypeText,
	yaml.StructureFieldTypeInteger:       psqldef.PSQLTypeInteger,
	yaml.StructureFieldTypeFloat:         psqldef.PSQLTypeDoublePrecision,
	yaml.StructureFieldTypeBoolean:       psqldef.PSQLTypeBoolean,
	yaml.StructureFieldTypeTime:          psqldef.PSQLTypeTimestampTZ,
	yaml.StructureFieldTypeDate:          psqldef.PSQLTypeDate,
	yaml.StructureFieldTypeProtected:     psqldef.PSQLTypeText,
	yaml.StructureFieldTypeSealed:        psqldef.PSQLTypeText,
}

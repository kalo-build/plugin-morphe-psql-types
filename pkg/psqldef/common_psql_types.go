package psqldef

var (
	PSQLTypeText = PSQLTypePrimitive{
		Syntax: "TEXT",
	}
	PSQLTypeVarchar = PSQLTypePrimitive{
		Syntax: "VARCHAR",
	}
	PSQLTypeChar = PSQLTypePrimitive{
		Syntax: "CHAR",
	}
	PSQLTypeBoolean = PSQLTypePrimitive{
		Syntax: "BOOLEAN",
	}
	PSQLTypeSmallInt = PSQLTypePrimitive{
		Syntax: "SMALLINT",
	}
	PSQLTypeInteger = PSQLTypePrimitive{
		Syntax: "INTEGER",
	}
	PSQLTypeBigInt = PSQLTypePrimitive{
		Syntax: "BIGINT",
	}
	PSQLTypeSerial = PSQLTypePrimitive{
		Syntax: "SERIAL",
	}
	PSQLTypeBigSerial = PSQLTypePrimitive{
		Syntax: "BIGSERIAL",
	}
	PSQLTypeReal = PSQLTypePrimitive{
		Syntax: "REAL",
	}
	PSQLTypeDoublePrecision = PSQLTypePrimitive{
		Syntax: "DOUBLE PRECISION",
	}
	PSQLTypeNumeric = PSQLTypePrimitive{
		Syntax: "NUMERIC",
	}
	PSQLTypeUUID = PSQLTypePrimitive{
		Syntax: "UUID",
	}
	PSQLTypeBytea = PSQLTypePrimitive{
		Syntax: "BYTEA",
	}
	PSQLTypeTimestamp = PSQLTypePrimitive{
		Syntax: "TIMESTAMP",
	}
	PSQLTypeTimestampTZ = PSQLTypePrimitive{
		Syntax: "TIMESTAMPTZ",
	}
	PSQLTypeDate = PSQLTypePrimitive{
		Syntax: "DATE",
	}
	PSQLTypeTime = PSQLTypePrimitive{
		Syntax: "TIME",
	}
	PSQLTypeTimeTZ = PSQLTypePrimitive{
		Syntax: "TIMETZ",
	}
	PSQLTypeInterval = PSQLTypePrimitive{
		Syntax: "INTERVAL",
	}
	PSQLTypeJSON = PSQLTypePrimitive{
		Syntax: "JSON",
	}
	PSQLTypeJSONB = PSQLTypePrimitive{
		Syntax: "JSONB",
	}
	PSQLTypeCIDR = PSQLTypePrimitive{
		Syntax: "CIDR",
	}
	PSQLTypeINET = PSQLTypePrimitive{
		Syntax: "INET",
	}
	PSQLTypeMACADDR = PSQLTypePrimitive{
		Syntax: "MACADDR",
	}
)

package compile_test

import (
	"fmt"
	"testing"

	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
	"github.com/stretchr/testify/suite"
)

type CompileEnumsTestSuite struct {
	suite.Suite
}

func TestCompileEnumsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEnumsTestSuite))
}

func (suite *CompileEnumsTestSuite) getMorpheEnumsConfig() cfg.MorpheEnumsConfig {
	return cfg.MorpheEnumsConfig{
		Schema: "public",
	}
}

func (suite *CompileEnumsTestSuite) SetupTest() {
}

func (suite *CompileEnumsTestSuite) TearDownTest() {
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_String() {
	enumHooks := hook.CompileMorpheEnum{}
	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin":  "ADMIN",
			"Editor": "EDITOR",
			"Viewer": "VIEWER",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.Nil(enumErr)
	suite.NotNil(lookupTable)
	suite.NotEmpty(lookupTable.SeedData)

	suite.Equal(config.Schema, lookupTable.Schema)
	suite.Equal("user_roles", lookupTable.Name)

	suite.Len(lookupTable.Columns, 4)

	columns := lookupTable.Columns

	column0 := columns[0]
	suite.Equal("id", column0.Name)
	suite.Equal(psqldef.PSQLTypeSerial, column0.Type)
	suite.True(column0.PrimaryKey)

	column1 := columns[1]
	suite.Equal("key", column1.Name)
	suite.Equal(psqldef.PSQLTypeText, column1.Type)
	suite.True(column1.NotNull)

	column2 := columns[2]
	suite.Equal("value", column2.Name)
	suite.Equal(psqldef.PSQLTypeText, column2.Type)
	suite.True(column2.NotNull)

	column3 := columns[3]
	suite.Equal("value_type", column3.Name)
	suite.Equal(psqldef.PSQLTypeText, column3.Type)
	suite.True(column3.NotNull)

	seedData := lookupTable.SeedData[0]
	suite.Len(seedData.Values, 3)

	seedData0 := seedData.Values[0]
	suite.Equal("Admin", seedData0[0])
	suite.Equal("ADMIN", seedData0[1])
	suite.Equal("String", seedData0[2])

	seedData1 := seedData.Values[1]
	suite.Equal("Editor", seedData1[0])
	suite.Equal("EDITOR", seedData1[1])
	suite.Equal("String", seedData1[2])

	seedData2 := seedData.Values[2]
	suite.Equal("Viewer", seedData2[0])
	suite.Equal("VIEWER", seedData2[1])
	suite.Equal("String", seedData2[2])

	suite.Len(lookupTable.UniqueConstraints, 1)
	uniqueConstraint00 := lookupTable.UniqueConstraints[0]
	suite.Equal("uk_user_roles_key", uniqueConstraint00.Name)
	suite.Equal("user_roles", uniqueConstraint00.TableName)
	suite.Len(uniqueConstraint00.ColumnNames, 1)
	suite.Equal("key", uniqueConstraint00.ColumnNames[0])
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_String_UseBigSerial() {
	enumHooks := hook.CompileMorpheEnum{}
	config := suite.getMorpheEnumsConfig()
	config.UseBigSerial = true

	enum0 := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin":  "ADMIN",
			"Editor": "EDITOR",
			"Viewer": "VIEWER",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.Nil(enumErr)
	suite.NotNil(lookupTable)
	suite.NotEmpty(lookupTable.SeedData)

	suite.Equal(config.Schema, lookupTable.Schema)
	suite.Equal("user_roles", lookupTable.Name)

	suite.Len(lookupTable.Columns, 4)

	columns := lookupTable.Columns

	column0 := columns[0]
	suite.Equal("id", column0.Name)
	suite.Equal(psqldef.PSQLTypeBigSerial, column0.Type)
	suite.True(column0.PrimaryKey)

	column1 := columns[1]
	suite.Equal("key", column1.Name)
	suite.Equal(psqldef.PSQLTypeText, column1.Type)
	suite.True(column1.NotNull)

	column2 := columns[2]
	suite.Equal("value", column2.Name)
	suite.Equal(psqldef.PSQLTypeText, column2.Type)
	suite.True(column2.NotNull)

	column3 := columns[3]
	suite.Equal("value_type", column3.Name)
	suite.Equal(psqldef.PSQLTypeText, column3.Type)
	suite.True(column3.NotNull)

	seedData := lookupTable.SeedData[0]
	suite.Len(seedData.Values, 3)

	seedData0 := seedData.Values[0]
	suite.Equal("Admin", seedData0[0])
	suite.Equal("ADMIN", seedData0[1])
	suite.Equal("String", seedData0[2])

	seedData1 := seedData.Values[1]
	suite.Equal("Editor", seedData1[0])
	suite.Equal("EDITOR", seedData1[1])
	suite.Equal("String", seedData1[2])

	seedData2 := seedData.Values[2]
	suite.Equal("Viewer", seedData2[0])
	suite.Equal("VIEWER", seedData2[1])
	suite.Equal("String", seedData2[2])

	suite.Len(lookupTable.UniqueConstraints, 1)
	uniqueConstraint00 := lookupTable.UniqueConstraints[0]
	suite.Equal("uk_user_roles_key", uniqueConstraint00.Name)
	suite.Equal("user_roles", uniqueConstraint00.TableName)
	suite.Len(uniqueConstraint00.ColumnNames, 1)
	suite.Equal("key", uniqueConstraint00.ColumnNames[0])
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_Float() {
	enumHooks := hook.CompileMorpheEnum{}
	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "Analytics",
		Type: yaml.EnumTypeFloat,
		Entries: map[string]any{
			"Pi":    3.141,
			"Euler": 2.718,
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.Nil(enumErr)
	suite.NotNil(lookupTable)
	suite.NotEmpty(lookupTable.SeedData)

	suite.Equal(config.Schema, lookupTable.Schema)
	suite.Equal("analytics", lookupTable.Name)

	suite.Len(lookupTable.Columns, 4)

	columns := lookupTable.Columns

	column0 := columns[0]
	suite.Equal("id", column0.Name)
	suite.Equal(psqldef.PSQLTypeSerial, column0.Type)
	suite.True(column0.PrimaryKey)

	column1 := columns[1]
	suite.Equal("key", column1.Name)
	suite.Equal(psqldef.PSQLTypeText, column1.Type)
	suite.True(column1.NotNull)

	column2 := columns[2]
	suite.Equal("value", column2.Name)
	suite.Equal(psqldef.PSQLTypeText, column2.Type)
	suite.True(column2.NotNull)

	column3 := columns[3]
	suite.Equal("value_type", column3.Name)
	suite.Equal(psqldef.PSQLTypeText, column3.Type)
	suite.True(column3.NotNull)

	seedData := lookupTable.SeedData[0]
	suite.Len(seedData.Values, 2)

	seedData0 := seedData.Values[0]
	suite.Equal("Euler", seedData0[0])
	suite.Equal("2.718", seedData0[1])
	suite.Equal("Float", seedData0[2])

	seedData1 := seedData.Values[1]
	suite.Equal("Pi", seedData1[0])
	suite.Equal("3.141", seedData1[1])
	suite.Equal("Float", seedData1[2])

	suite.Len(lookupTable.UniqueConstraints, 1)
	uniqueConstraint00 := lookupTable.UniqueConstraints[0]
	suite.Equal("uk_analytics_key", uniqueConstraint00.Name)
	suite.Equal("analytics", uniqueConstraint00.TableName)
	suite.Len(uniqueConstraint00.ColumnNames, 1)
	suite.Equal("key", uniqueConstraint00.ColumnNames[0])
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_Integer() {
	enumHooks := hook.CompileMorpheEnum{}
	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "Analytics",
		Type: yaml.EnumTypeInteger,
		Entries: map[string]any{
			"AnswerToLife":  42,
			"FineStructure": 317,
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.Nil(enumErr)
	suite.NotNil(lookupTable)
	suite.NotEmpty(lookupTable.SeedData)

	suite.Equal(config.Schema, lookupTable.Schema)
	suite.Equal("analytics", lookupTable.Name)

	suite.Len(lookupTable.Columns, 4)

	columns := lookupTable.Columns

	column0 := columns[0]
	suite.Equal("id", column0.Name)
	suite.Equal(psqldef.PSQLTypeSerial, column0.Type)
	suite.True(column0.PrimaryKey)

	column1 := columns[1]
	suite.Equal("key", column1.Name)
	suite.Equal(psqldef.PSQLTypeText, column1.Type)
	suite.True(column1.NotNull)

	column2 := columns[2]
	suite.Equal("value", column2.Name)
	suite.Equal(psqldef.PSQLTypeText, column2.Type)
	suite.True(column2.NotNull)

	column3 := columns[3]
	suite.Equal("value_type", column3.Name)
	suite.Equal(psqldef.PSQLTypeText, column3.Type)
	suite.True(column3.NotNull)

	seedData := lookupTable.SeedData[0]
	suite.Len(seedData.Values, 2)

	seedData0 := seedData.Values[0]
	suite.Equal("AnswerToLife", seedData0[0])
	suite.Equal("42", seedData0[1])
	suite.Equal("Integer", seedData0[2])

	seedData1 := seedData.Values[1]
	suite.Equal("FineStructure", seedData1[0])
	suite.Equal("317", seedData1[1])
	suite.Equal("Integer", seedData1[2])

	suite.Len(lookupTable.UniqueConstraints, 1)
	uniqueConstraint00 := lookupTable.UniqueConstraints[0]
	suite.Equal("uk_analytics_key", uniqueConstraint00.Name)
	suite.Equal("analytics", uniqueConstraint00.TableName)
	suite.Len(uniqueConstraint00.ColumnNames, 1)
	suite.Equal("key", uniqueConstraint00.ColumnNames[0])
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_NoName() {
	enumHooks := hook.CompileMorpheEnum{}
	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin":  "ADMIN",
			"Editor": "EDITOR",
			"Viewer": "VIEWER",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.ErrorIs(enumErr, yaml.ErrNoMorpheEnumName)
	suite.Nil(lookupTable)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_NoType() {
	enumHooks := hook.CompileMorpheEnum{}
	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "UserRole",
		Entries: map[string]any{
			"Admin":  "ADMIN",
			"Editor": "EDITOR",
			"Viewer": "VIEWER",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.ErrorIs(enumErr, yaml.ErrNoMorpheEnumType)
	suite.Nil(lookupTable)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_NoEntries() {
	enumHooks := hook.CompileMorpheEnum{}
	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name:    "UserRole",
		Type:    yaml.EnumTypeString,
		Entries: map[string]any{},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.ErrorIs(enumErr, yaml.ErrNoMorpheEnumEntries)
	suite.Nil(lookupTable)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_EntryTypeMismatch() {
	enumHooks := hook.CompileMorpheEnum{}
	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeInteger,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.ErrorContains(enumErr, "enum entry 'Blue' value 'rgb(0,0,255)' with type 'string' does not match the enum type of 'Integer'")
	suite.Nil(lookupTable)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_StartHook_Successful() {
	var featureFlag = "otherName"
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumStart: func(config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error) {
			if featureFlag != "otherName" {
				return config, enum, nil
			}
			enum.Name = enum.Name + "CHANGED"
			delete(enum.Entries, "Green")
			return config, enum, nil
		},
	}

	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.Nil(enumErr)
	suite.NotNil(lookupTable)
	suite.NotEmpty(lookupTable.SeedData)

	suite.Equal(config.Schema, lookupTable.Schema)
	suite.Equal("color_changeds", lookupTable.Name)

	suite.Len(lookupTable.Columns, 4)

	columns := lookupTable.Columns

	column0 := columns[0]
	suite.Equal("id", column0.Name)
	suite.Equal(psqldef.PSQLTypeSerial, column0.Type)
	suite.True(column0.PrimaryKey)

	column1 := columns[1]
	suite.Equal("key", column1.Name)
	suite.Equal(psqldef.PSQLTypeText, column1.Type)
	suite.True(column1.NotNull)

	column2 := columns[2]
	suite.Equal("value", column2.Name)
	suite.Equal(psqldef.PSQLTypeText, column2.Type)
	suite.True(column2.NotNull)

	column3 := columns[3]
	suite.Equal("value_type", column3.Name)
	suite.Equal(psqldef.PSQLTypeText, column3.Type)
	suite.True(column3.NotNull)

	seedData := lookupTable.SeedData[0]
	suite.Len(seedData.Values, 2)

	seedData0 := seedData.Values[0]
	suite.Equal("Blue", seedData0[0])
	suite.Equal("rgb(0,0,255)", seedData0[1])
	suite.Equal("String", seedData0[2])

	seedData1 := seedData.Values[1]
	suite.Equal("Red", seedData1[0])
	suite.Equal("rgb(255,0,0)", seedData1[1])
	suite.Equal("String", seedData1[2])

	suite.Len(lookupTable.UniqueConstraints, 1)
	uniqueConstraint00 := lookupTable.UniqueConstraints[0]
	suite.Equal("uk_color_changeds_key", uniqueConstraint00.Name)
	suite.Equal("color_changeds", uniqueConstraint00.TableName)
	suite.Len(uniqueConstraint00.ColumnNames, 1)
	suite.Equal("key", uniqueConstraint00.ColumnNames[0])
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_StartHook_Failure() {
	var featureFlag = "otherName"
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumStart: func(config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error) {
			if featureFlag != "otherName" {
				return config, enum, nil
			}
			return config, enum, fmt.Errorf("compile enum start hook error")
		},
	}

	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.ErrorContains(enumErr, "compile enum start hook error")
	suite.Nil(lookupTable)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_SuccessHook_Successful() {
	var featureFlag = "otherName"
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumSuccess: func(table *psqldef.Table) (*psqldef.Table, error) {
			if featureFlag != "otherName" {
				return table, nil
			}
			table.Name = table.Name + "_changed"
			table.Columns = append(table.Columns, psqldef.Column{
				Name:    "description",
				Type:    psqldef.PSQLTypeText,
				NotNull: false,
			})
			table.UniqueConstraints = []psqldef.UniqueConstraint{}

			// Update seed data directly in the table
			if len(table.SeedData) > 0 {
				table.SeedData[0].Values = [][]any{
					{
						"Orange",
						"rgb(255,165,0)",
						"String",
					},
				}
			}

			return table, nil
		},
	}

	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.Nil(enumErr)
	suite.NotNil(lookupTable)
	suite.NotEmpty(lookupTable.SeedData)

	suite.Equal(config.Schema, lookupTable.Schema)
	suite.Equal("colors_changed", lookupTable.Name)

	suite.Len(lookupTable.Columns, 5)

	columns := lookupTable.Columns

	column0 := columns[0]
	suite.Equal("id", column0.Name)
	suite.Equal(psqldef.PSQLTypeSerial, column0.Type)
	suite.True(column0.PrimaryKey)

	column1 := columns[1]
	suite.Equal("key", column1.Name)
	suite.Equal(psqldef.PSQLTypeText, column1.Type)
	suite.True(column1.NotNull)

	column2 := columns[2]
	suite.Equal("value", column2.Name)
	suite.Equal(psqldef.PSQLTypeText, column2.Type)
	suite.True(column2.NotNull)

	column3 := columns[3]
	suite.Equal("value_type", column3.Name)
	suite.Equal(psqldef.PSQLTypeText, column3.Type)
	suite.True(column3.NotNull)

	column4 := columns[4]
	suite.Equal("description", column4.Name)
	suite.Equal(psqldef.PSQLTypeText, column4.Type)
	suite.False(column4.NotNull)

	seedData := lookupTable.SeedData[0]
	suite.Len(seedData.Values, 1)

	seedData0 := seedData.Values[0]
	suite.Equal("Orange", seedData0[0])
	suite.Equal("rgb(255,165,0)", seedData0[1])
	suite.Equal("String", seedData0[2])

	suite.Len(lookupTable.UniqueConstraints, 0)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_SuccessHook_Failure() {
	var featureFlag = "otherName"
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumSuccess: func(table *psqldef.Table) (*psqldef.Table, error) {
			if featureFlag != "otherName" {
				return table, nil
			}
			return table, fmt.Errorf("compile enum success hook error")
		},
	}

	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.ErrorContains(enumErr, "compile enum success hook error")
	suite.Nil(lookupTable)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_FailureHook() {
	var featureFlag = "otherName"
	var failureErr error
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumFailure: func(config cfg.MorpheEnumsConfig, enum yaml.Enum, err error) error {
			if featureFlag != "otherName" {
				return err
			}
			failureErr = err
			return fmt.Errorf("compile enum failure hook error: %w", err)
		},
	}

	config := suite.getMorpheEnumsConfig()

	enum0 := yaml.Enum{
		Name: "",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	lookupTable, enumErr := compile.MorpheEnumToPSQLTable(enumHooks, config, enum0)

	suite.ErrorIs(failureErr, yaml.ErrNoMorpheEnumName)
	suite.ErrorContains(enumErr, "compile enum failure hook error")
	suite.Nil(lookupTable)
}

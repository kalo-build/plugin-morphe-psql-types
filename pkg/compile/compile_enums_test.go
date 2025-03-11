package compile_test

import (
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

func (suite *CompileEnumsTestSuite) getMorpheConfig() cfg.MorpheConfig {
	return cfg.MorpheConfig{
		MorpheEnumsConfig: cfg.MorpheEnumsConfig{
			Schema: "public",
		},
	}
}

func (suite *CompileEnumsTestSuite) getCompileConfig() compile.MorpheCompileConfig {
	return compile.MorpheCompileConfig{
		MorpheConfig: suite.getMorpheConfig(),
		EnumHooks:    hook.CompileMorpheEnum{},
	}
}

func (suite *CompileEnumsTestSuite) SetupTest() {
}

func (suite *CompileEnumsTestSuite) TearDownTest() {
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToPSQLTable_String() {
	enumHooks := hook.CompileMorpheEnum{}
	config := suite.getCompileConfig()
	config.EnumHooks = enumHooks

	enum0 := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin":  "ADMIN",
			"Editor": "EDITOR",
			"Viewer": "VIEWER",
		},
	}

	lookupTable, seedData, enumErr := compile.MorpheEnumToPSQLTable(config, enum0)

	suite.Nil(enumErr)
	suite.NotNil(lookupTable)
	suite.NotNil(seedData)

	suite.Equal(config.MorpheConfig.MorpheEnumsConfig.Schema, lookupTable.Schema)
	suite.Equal("user_roles", lookupTable.Name)

	suite.Len(lookupTable.Columns, 3)

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

	suite.Len(seedData.Values, 3)

	seedData0 := seedData.Values[0]
	suite.Equal("Admin", seedData0[0])
	suite.Equal("ADMIN", seedData0[1])

	seedData1 := seedData.Values[1]
	suite.Equal("Editor", seedData1[0])
	suite.Equal("EDITOR", seedData1[1])

	seedData2 := seedData.Values[2]
	suite.Equal("Viewer", seedData2[0])
	suite.Equal("VIEWER", seedData2[1])

	suite.Len(lookupTable.UniqueConstraints, 1)
	uniqueConstraint00 := lookupTable.UniqueConstraints[0]
	suite.Equal("uk_user_roles_key", uniqueConstraint00.Name)
	suite.Equal("user_roles", uniqueConstraint00.TableName)
	suite.Len(uniqueConstraint00.ColumnNames, 1)
	suite.Equal("key", uniqueConstraint00.ColumnNames[0])
}

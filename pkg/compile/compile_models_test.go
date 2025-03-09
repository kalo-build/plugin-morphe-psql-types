package compile_test

import (
	"testing"

	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
	"github.com/stretchr/testify/suite"
)

type CompileModelsTestSuite struct {
	suite.Suite
}

func TestCompileModelsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileModelsTestSuite))
}

func (suite *CompileModelsTestSuite) getMorpheConfig() cfg.MorpheConfig {
	modelsConfig := cfg.MorpheModelsConfig{
		Schema:       "public",
		UseBigSerial: false,
	}
	return cfg.MorpheConfig{
		MorpheModelsConfig: modelsConfig,
	}
}

func (suite *CompileModelsTestSuite) getCompileConfig() compile.MorpheCompileConfig {
	morpheConfig := suite.getMorpheConfig()
	return compile.MorpheCompileConfig{
		MorpheConfig: morpheConfig,
		ModelHooks:   hook.CompileMorpheModel{},
	}
}

// func (suite *CompileModelsTestSuite) getCompileConfigWithHooks(hooks hook.CompileMorpheModel) compile.MorpheCompileConfig {
// 	config := suite.getCompileConfig()
// 	config.ModelHooks = hooks
// 	return config
// }

func (suite *CompileModelsTestSuite) SetupTest() {
}

func (suite *CompileModelsTestSuite) TearDownTest() {
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables() {
	config := suite.getCompileConfig()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Boolean": {
				Type: yaml.ModelFieldTypeBoolean,
			},
			"Date": {
				Type: yaml.ModelFieldTypeDate,
			},
			"Float": {
				Type: yaml.ModelFieldTypeFloat,
			},
			"Integer": {
				Type: yaml.ModelFieldTypeInteger,
			},
			"Protected": {
				Type: yaml.ModelFieldTypeProtected,
			},
			"Sealed": {
				Type: yaml.ModelFieldTypeSealed,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
			"Time": {
				Type: yaml.ModelFieldTypeTime,
			},
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
					"primary",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"UUID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table0 := allTables[0]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table0.Schema)
	suite.Equal("basics", table0.Name)

	columns := table0.Columns
	suite.Len(columns, 10)

	columns00 := columns[0]
	suite.Equal("auto_increment", columns00.Name)
	suite.Equal(psqldef.PSQLTypeSerial, columns00.Type)
	suite.False(columns00.NotNull)
	suite.False(columns00.PrimaryKey)
	suite.Equal("", columns00.Default)

	columns01 := columns[1]
	suite.Equal("boolean", columns01.Name)
	suite.Equal(psqldef.PSQLTypeBoolean, columns01.Type)
	suite.False(columns01.NotNull)
	suite.False(columns01.PrimaryKey)
	suite.Equal("", columns01.Default)

	columns02 := columns[2]
	suite.Equal("date", columns02.Name)
	suite.Equal(psqldef.PSQLTypeDate, columns02.Type)
	suite.False(columns02.NotNull)
	suite.False(columns02.PrimaryKey)
	suite.Equal("", columns02.Default)

	columns03 := columns[3]
	suite.Equal("float", columns03.Name)
	suite.Equal(psqldef.PSQLTypeDoublePrecision, columns03.Type)
	suite.False(columns03.NotNull)
	suite.False(columns03.PrimaryKey)
	suite.Equal("", columns03.Default)

	columns04 := columns[4]
	suite.Equal("integer", columns04.Name)
	suite.Equal(psqldef.PSQLTypeInteger, columns04.Type)
	suite.False(columns04.NotNull)
	suite.False(columns04.PrimaryKey)
	suite.Equal("", columns04.Default)

	columns05 := columns[5]
	suite.Equal("protected", columns05.Name)
	suite.Equal(psqldef.PSQLTypeText, columns05.Type)
	suite.False(columns05.NotNull)
	suite.False(columns05.PrimaryKey)
	suite.Equal("", columns05.Default)

	columns06 := columns[6]
	suite.Equal("sealed", columns06.Name)
	suite.Equal(psqldef.PSQLTypeText, columns06.Type)
	suite.False(columns06.NotNull)
	suite.False(columns06.PrimaryKey)
	suite.Equal("", columns06.Default)

	columns07 := columns[7]
	suite.Equal("string", columns07.Name)
	suite.Equal(psqldef.PSQLTypeText, columns07.Type)
	suite.False(columns07.NotNull)
	suite.False(columns07.PrimaryKey)
	suite.Equal("", columns07.Default)

	columns08 := columns[8]
	suite.Equal("time", columns08.Name)
	suite.Equal(psqldef.PSQLTypeTimestampTZ, columns08.Type)
	suite.False(columns08.NotNull)
	suite.False(columns08.PrimaryKey)
	suite.Equal("", columns08.Default)

	columns09 := columns[9]
	suite.Equal("uuid", columns09.Name)
	suite.Equal(psqldef.PSQLTypeUUID, columns09.Type)
	suite.False(columns09.NotNull)
	suite.True(columns09.PrimaryKey)
	suite.Equal("", columns09.Default)

	suite.Len(table0.Indices, 0)
	suite.Len(table0.ForeignKeys, 0)
	suite.Len(table0.UniqueConstraints, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_UseBigSerial() {
	config := suite.getCompileConfig()
	config.UseBigSerial = true

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
					"primary",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"UUID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table0 := allTables[0]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table0.Schema)
	suite.Equal("basics", table0.Name)

	columns := table0.Columns
	suite.Len(columns, 2)

	columns00 := columns[0]
	suite.Equal("auto_increment", columns00.Name)
	suite.Equal(psqldef.PSQLTypeBigSerial, columns00.Type)
	suite.False(columns00.NotNull)
	suite.False(columns00.PrimaryKey)
	suite.Equal("", columns00.Default)

	columns01 := columns[1]
	suite.Equal("uuid", columns01.Name)
	suite.Equal(psqldef.PSQLTypeUUID, columns01.Type)
	suite.False(columns01.NotNull)
	suite.True(columns01.PrimaryKey)
	suite.Equal("", columns01.Default)

	suite.Len(table0.Indices, 0)
	suite.Len(table0.ForeignKeys, 0)
	suite.Len(table0.UniqueConstraints, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_NoSchema() {
	config := suite.getCompileConfig()
	config.Schema = ""

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
					"primary",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"UUID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.NotNil(allTablesErr)
	suite.ErrorIs(allTablesErr, cfg.ErrNoModelSchema)
	suite.Len(allTables, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_NoModelName() {
	config := suite.getCompileConfig()

	model0 := yaml.Model{
		Name: "",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
					"primary",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"UUID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "morphe model has no name")
	suite.Len(allTables, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_NoFields() {
	config := suite.getCompileConfig()

	model0 := yaml.Model{
		Name:   "Basic",
		Fields: map[string]yaml.ModelField{},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"UUID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "morphe model has no fields")
	suite.Len(allTables, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_NoIdentifiers() {
	config := suite.getCompileConfig()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
					"primary",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{},
		Related:     map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "morphe model has no identifiers")
	suite.Len(allTables, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForOne() {
	config := suite.getCompileConfig()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	model1 := yaml.Model{
		Name: "BasicParent",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r := registry.NewRegistry()
	r.SetModel("Basic", model0)
	r.SetModel("BasicParent", model1)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table0 := allTables[0]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table0.Schema)
	suite.Equal("basics", table0.Name)

	columns0 := table0.Columns
	suite.Len(columns0, 3)

	columns00 := columns0[0]
	suite.Equal("id", columns00.Name)
	suite.Equal(psqldef.PSQLTypeSerial, columns00.Type)
	suite.False(columns00.NotNull)
	suite.True(columns00.PrimaryKey)
	suite.Equal("", columns00.Default)

	columns01 := columns0[1]
	suite.Equal("string", columns01.Name)
	suite.Equal(psqldef.PSQLTypeText, columns01.Type)
	suite.False(columns01.NotNull)
	suite.False(columns01.PrimaryKey)
	suite.Equal("", columns01.Default)

	columns02 := columns0[2]
	suite.Equal("basic_parent_id", columns02.Name)
	suite.Equal(psqldef.PSQLTypeInteger, columns02.Type)
	suite.False(columns01.NotNull)
	suite.False(columns01.PrimaryKey)
	suite.Equal("", columns01.Default)

	suite.Len(table0.ForeignKeys, 1)

	foreignKey0 := table0.ForeignKeys[0]
	suite.Equal("public", foreignKey0.Schema)
	suite.Equal("fk_basic_basic_parent_id", foreignKey0.Name)
	suite.Equal("basics", foreignKey0.TableName)
	suite.Len(foreignKey0.ColumnNames, 1)
	fkColumn00 := foreignKey0.ColumnNames[0]
	suite.Equal("basic_parent_id", fkColumn00)
	suite.Equal("basic_parents", foreignKey0.RefTableName)
	suite.Len(foreignKey0.RefColumnNames, 1)
	fkColumnRef00 := foreignKey0.RefColumnNames[0]
	suite.Equal("id", fkColumnRef00)
	suite.Equal("CASCADE", foreignKey0.OnDelete)
	suite.Equal("", foreignKey0.OnUpdate)

	suite.Len(table0.Indices, 1)
	index0 := table0.Indices[0]
	suite.Equal("idx_basics_basic_parent_id", index0.Name)
	suite.Equal("basics", index0.TableName)
	suite.Len(index0.Columns, 1)
	suite.Equal("basic_parent_id", index0.Columns[0])
	suite.False(index0.IsUnique)

	suite.Len(table0.UniqueConstraints, 0)
}

// TODO: Test cases
// _Related*
// _EnumField
// _*Hook*

package compile_test

import (
	"fmt"
	"testing"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
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
	enumsConfig := cfg.MorpheEnumsConfig{
		Schema:       "public",
		UseBigSerial: false,
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Schema:         "public",
		ViewNameSuffix: "_entities",
	}
	return cfg.MorpheConfig{
		MorpheModelsConfig:   modelsConfig,
		MorpheEnumsConfig:    enumsConfig,
		MorpheEntitiesConfig: entitiesConfig,
	}
}

func (suite *CompileModelsTestSuite) getCompileConfig() compile.MorpheCompileConfig {
	morpheConfig := suite.getMorpheConfig()
	return compile.MorpheCompileConfig{
		MorpheConfig: morpheConfig,
		ModelHooks:   hook.CompileMorpheModel{},
	}
}

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
	config.MorpheModelsConfig.UseBigSerial = true

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
	config.MorpheModelsConfig.Schema = ""

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
	suite.Equal("fk_basics_basic_parent_id", foreignKey0.Name)
	suite.Equal("basics", foreignKey0.TableName)
	suite.Len(foreignKey0.ColumnNames, 1)
	fkColumn00 := foreignKey0.ColumnNames[0]
	suite.Equal("basic_parent_id", fkColumn00)
	suite.Equal("public", foreignKey0.RefSchema)
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

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForMany_HasOne() {
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
				Type: "ForMany",
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
				Type: "HasOne",
			},
		},
	}
	r := registry.NewRegistry()
	r.SetModel("Basic", model0)
	r.SetModel("BasicParent", model1)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 2)

	table0 := allTables[0]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table0.Schema)
	suite.Equal("basics", table0.Name)

	columns0 := table0.Columns
	suite.Len(columns0, 2)

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

	suite.Len(table0.ForeignKeys, 0)
	suite.Len(table0.Indices, 0)
	suite.Len(table0.UniqueConstraints, 0)

	// Junction table basics <-> basic_parents
	table1 := allTables[1]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table1.Schema)
	suite.Equal("basic_basic_parents", table1.Name)

	columns1 := table1.Columns
	suite.Len(columns1, 3)

	columns10 := columns1[0]
	suite.Equal("id", columns10.Name)
	suite.Equal(psqldef.PSQLTypeSerial, columns10.Type)
	suite.False(columns10.NotNull)
	suite.True(columns10.PrimaryKey)
	suite.Equal("", columns10.Default)

	columns11 := columns1[1]
	suite.Equal("basic_id", columns11.Name)
	suite.Equal(psqldef.PSQLTypeInteger, columns11.Type)
	suite.False(columns11.NotNull)
	suite.False(columns11.PrimaryKey)
	suite.Equal("", columns11.Default)

	columns12 := columns1[2]
	suite.Equal("basic_parent_id", columns12.Name)
	suite.Equal(psqldef.PSQLTypeInteger, columns12.Type)
	suite.False(columns12.NotNull)
	suite.False(columns12.PrimaryKey)
	suite.Equal("", columns12.Default)

	suite.Len(table1.ForeignKeys, 2)
	foreignKey10 := table1.ForeignKeys[0]
	suite.Equal("public", foreignKey10.Schema)
	suite.Equal("fk_basic_basic_parents_basic_id", foreignKey10.Name)
	suite.Equal("basic_basic_parents", foreignKey10.TableName)
	suite.Len(foreignKey10.ColumnNames, 1)
	fkColumn10 := foreignKey10.ColumnNames[0]
	suite.Equal("basic_id", fkColumn10)
	suite.Equal("public", foreignKey10.RefSchema)
	suite.Equal("basics", foreignKey10.RefTableName)
	suite.Len(foreignKey10.RefColumnNames, 1)
	fkColumnRef10 := foreignKey10.RefColumnNames[0]
	suite.Equal("id", fkColumnRef10)
	suite.Equal("CASCADE", foreignKey10.OnDelete)
	suite.Equal("", foreignKey10.OnUpdate)

	foreignKey11 := table1.ForeignKeys[1]
	suite.Equal("public", foreignKey11.Schema)
	suite.Equal("fk_basic_basic_parents_basic_parent_id", foreignKey11.Name)
	suite.Equal("basic_basic_parents", foreignKey11.TableName)
	suite.Len(foreignKey11.ColumnNames, 1)
	fkColumn11 := foreignKey11.ColumnNames[0]
	suite.Equal("basic_parent_id", fkColumn11)
	suite.Equal("public", foreignKey11.RefSchema)
	suite.Equal("basic_parents", foreignKey11.RefTableName)
	suite.Len(foreignKey11.RefColumnNames, 1)
	fkColumnRef11 := foreignKey11.RefColumnNames[0]
	suite.Equal("id", fkColumnRef11)
	suite.Equal("CASCADE", foreignKey11.OnDelete)
	suite.Equal("", foreignKey11.OnUpdate)

	suite.Len(table1.Indices, 2)
	index10 := table1.Indices[0]
	suite.Equal("idx_basic_basic_parents_basic_id", index10.Name)
	suite.Equal("basic_basic_parents", index10.TableName)
	suite.Len(index10.Columns, 1)
	suite.Equal("basic_id", index10.Columns[0])
	suite.False(index10.IsUnique)

	index11 := table1.Indices[1]
	suite.Equal("idx_basic_basic_parents_basic_parent_id", index11.Name)
	suite.Equal("basic_basic_parents", index11.TableName)
	suite.Len(index11.Columns, 1)
	suite.Equal("basic_parent_id", index11.Columns[0])
	suite.False(index11.IsUnique)

	suite.Len(table1.UniqueConstraints, 1)
	uniqueConstraint10 := table1.UniqueConstraints[0]
	suite.Equal("uk_basic_basic_parents_basic_id_basic_parent_id", uniqueConstraint10.Name)
	suite.Equal("basic_basic_parents", uniqueConstraint10.TableName)
	suite.Len(uniqueConstraint10.ColumnNames, 2)
	suite.Equal("basic_id", uniqueConstraint10.ColumnNames[0])
	suite.Equal("basic_parent_id", uniqueConstraint10.ColumnNames[1])
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForMany_HasMany() {
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
				Type: "ForMany",
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
	suite.Len(allTables, 2)

	table0 := allTables[0]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table0.Schema)
	suite.Equal("basics", table0.Name)

	columns0 := table0.Columns
	suite.Len(columns0, 2)

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

	suite.Len(table0.ForeignKeys, 0)
	suite.Len(table0.Indices, 0)
	suite.Len(table0.UniqueConstraints, 0)

	// Junction table basics <-> basic_parents
	table1 := allTables[1]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table1.Schema)
	suite.Equal("basic_basic_parents", table1.Name)

	columns1 := table1.Columns
	suite.Len(columns1, 3)

	columns10 := columns1[0]
	suite.Equal("id", columns10.Name)
	suite.Equal(psqldef.PSQLTypeSerial, columns10.Type)
	suite.False(columns10.NotNull)
	suite.True(columns10.PrimaryKey)
	suite.Equal("", columns10.Default)

	columns11 := columns1[1]
	suite.Equal("basic_id", columns11.Name)
	suite.Equal(psqldef.PSQLTypeInteger, columns11.Type)
	suite.False(columns11.NotNull)
	suite.False(columns11.PrimaryKey)
	suite.Equal("", columns11.Default)

	columns12 := columns1[2]
	suite.Equal("basic_parent_id", columns12.Name)
	suite.Equal(psqldef.PSQLTypeInteger, columns12.Type)
	suite.False(columns12.NotNull)
	suite.False(columns12.PrimaryKey)
	suite.Equal("", columns12.Default)

	suite.Len(table1.ForeignKeys, 2)
	foreignKey10 := table1.ForeignKeys[0]
	suite.Equal("public", foreignKey10.Schema)
	suite.Equal("fk_basic_basic_parents_basic_id", foreignKey10.Name)
	suite.Equal("basic_basic_parents", foreignKey10.TableName)
	suite.Len(foreignKey10.ColumnNames, 1)
	fkColumn10 := foreignKey10.ColumnNames[0]
	suite.Equal("basic_id", fkColumn10)
	suite.Equal("public", foreignKey10.RefSchema)
	suite.Equal("basics", foreignKey10.RefTableName)
	suite.Len(foreignKey10.RefColumnNames, 1)
	fkColumnRef10 := foreignKey10.RefColumnNames[0]
	suite.Equal("id", fkColumnRef10)
	suite.Equal("CASCADE", foreignKey10.OnDelete)
	suite.Equal("", foreignKey10.OnUpdate)

	foreignKey11 := table1.ForeignKeys[1]
	suite.Equal("public", foreignKey11.Schema)
	suite.Equal("fk_basic_basic_parents_basic_parent_id", foreignKey11.Name)
	suite.Equal("basic_basic_parents", foreignKey11.TableName)
	suite.Len(foreignKey11.ColumnNames, 1)
	fkColumn11 := foreignKey11.ColumnNames[0]
	suite.Equal("basic_parent_id", fkColumn11)
	suite.Equal("public", foreignKey11.RefSchema)
	suite.Equal("basic_parents", foreignKey11.RefTableName)
	suite.Len(foreignKey11.RefColumnNames, 1)
	fkColumnRef11 := foreignKey11.RefColumnNames[0]
	suite.Equal("id", fkColumnRef11)
	suite.Equal("CASCADE", foreignKey11.OnDelete)
	suite.Equal("", foreignKey11.OnUpdate)

	suite.Len(table1.Indices, 2)
	index10 := table1.Indices[0]
	suite.Equal("idx_basic_basic_parents_basic_id", index10.Name)
	suite.Equal("basic_basic_parents", index10.TableName)
	suite.Len(index10.Columns, 1)
	suite.Equal("basic_id", index10.Columns[0])
	suite.False(index10.IsUnique)

	index11 := table1.Indices[1]
	suite.Equal("idx_basic_basic_parents_basic_parent_id", index11.Name)
	suite.Equal("basic_basic_parents", index11.TableName)
	suite.Len(index11.Columns, 1)
	suite.Equal("basic_parent_id", index11.Columns[0])
	suite.False(index11.IsUnique)

	suite.Len(table1.UniqueConstraints, 1)
	uniqueConstraint10 := table1.UniqueConstraints[0]
	suite.Equal("uk_basic_basic_parents_basic_id_basic_parent_id", uniqueConstraint10.Name)
	suite.Equal("basic_basic_parents", uniqueConstraint10.TableName)
	suite.Len(uniqueConstraint10.ColumnNames, 2)
	suite.Equal("basic_id", uniqueConstraint10.ColumnNames[0])
	suite.Equal("basic_parent_id", uniqueConstraint10.ColumnNames[1])
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_HasOne() {
	config := suite.getCompileConfig()

	model0 := yaml.Model{
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
				Type: "HasOne",
			},
		},
	}

	model1 := yaml.Model{
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
	r := registry.NewRegistry()
	r.SetModel("BasicParent", model0)
	r.SetModel("Basic", model1)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table0 := allTables[0]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table0.Schema)
	suite.Equal("basic_parents", table0.Name)

	columns0 := table0.Columns
	suite.Len(columns0, 2)

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

	suite.Len(table0.ForeignKeys, 0)
	suite.Len(table0.Indices, 0)
	suite.Len(table0.UniqueConstraints, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_HasMany() {
	config := suite.getCompileConfig()

	model0 := yaml.Model{
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

	model1 := yaml.Model{
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
	r := registry.NewRegistry()
	r.SetModel("BasicParent", model0)
	r.SetModel("Basic", model1)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table0 := allTables[0]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table0.Schema)
	suite.Equal("basic_parents", table0.Name)

	columns0 := table0.Columns
	suite.Len(columns0, 2)

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

	suite.Len(table0.ForeignKeys, 0)
	suite.Len(table0.Indices, 0)
	suite.Len(table0.UniqueConstraints, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_StartHook_Successful() {
	var featureFlag = "otherName"
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelStart: func(config cfg.MorpheConfig, model yaml.Model) (cfg.MorpheConfig, yaml.Model, error) {
			if featureFlag != "otherName" {
				return config, model, nil
			}
			config.MorpheModelsConfig.UseBigSerial = true
			model.Name = model.Name + "CHANGED"
			delete(model.Fields, "Float")
			return config, model, nil
		},
	}
	config := suite.getCompileConfig()
	config.ModelHooks = modelHooks

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
	suite.Equal("basic_changeds", table0.Name)

	columns0 := table0.Columns
	suite.Len(columns0, 9)

	column00 := columns0[0]
	suite.Equal("auto_increment", column00.Name)
	suite.Equal(psqldef.PSQLTypeBigSerial, column00.Type)

	column01 := columns0[1]
	suite.Equal("boolean", column01.Name)
	suite.Equal(psqldef.PSQLTypeBoolean, column01.Type)

	column02 := columns0[2]
	suite.Equal("date", column02.Name)
	suite.Equal(psqldef.PSQLTypeDate, column02.Type)

	column03 := columns0[3]
	suite.Equal("integer", column03.Name)
	suite.Equal(psqldef.PSQLTypeInteger, column03.Type)

	column04 := columns0[4]
	suite.Equal("protected", column04.Name)
	suite.Equal(psqldef.PSQLTypeText, column04.Type)

	column05 := columns0[5]
	suite.Equal("sealed", column05.Name)
	suite.Equal(psqldef.PSQLTypeText, column05.Type)

	column06 := columns0[6]
	suite.Equal("string", column06.Name)
	suite.Equal(psqldef.PSQLTypeText, column06.Type)

	column07 := columns0[7]
	suite.Equal("time", column07.Name)
	suite.Equal(psqldef.PSQLTypeTimestampTZ, column07.Type)

	column08 := columns0[8]
	suite.Equal("uuid", column08.Name)
	suite.Equal(psqldef.PSQLTypeUUID, column08.Type)
	suite.True(column08.PrimaryKey)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_StartHook_Failure() {
	var featureFlag = "otherName"
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelStart: func(config cfg.MorpheConfig, model yaml.Model) (cfg.MorpheConfig, yaml.Model, error) {
			if featureFlag != "otherName" {
				return config, model, nil
			}
			return config, model, fmt.Errorf("compile model start hook error")
		},
	}
	config := suite.getCompileConfig()
	config.ModelHooks = modelHooks

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
	suite.ErrorContains(allTablesErr, "compile model start hook error")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_SuccessHook_Successful() {
	var featureFlag = "otherName"
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelSuccess: func(allModelTables []*psqldef.Table) ([]*psqldef.Table, error) {
			if featureFlag != "otherName" {
				return allModelTables, nil
			}
			for _, modelTablePtr := range allModelTables {
				modelTablePtr.Name = modelTablePtr.Name + "_changed"
				newColumns := []psqldef.TableColumn{}
				for _, modelTableColumn := range modelTablePtr.Columns {
					if modelTableColumn.Name == "float" {
						continue
					}
					newColumns = append(newColumns, modelTableColumn)
				}
				modelTablePtr.Columns = newColumns
			}
			return allModelTables, nil
		},
	}
	config := suite.getCompileConfig()
	config.ModelHooks = modelHooks

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
	suite.Equal("basics_changed", table0.Name)

	columns0 := table0.Columns
	suite.Len(columns0, 9)

	column00 := columns0[0]
	suite.Equal("auto_increment", column00.Name)
	suite.Equal(psqldef.PSQLTypeSerial, column00.Type)

	column01 := columns0[1]
	suite.Equal("boolean", column01.Name)
	suite.Equal(psqldef.PSQLTypeBoolean, column01.Type)

	column02 := columns0[2]
	suite.Equal("date", column02.Name)
	suite.Equal(psqldef.PSQLTypeDate, column02.Type)

	column03 := columns0[3]
	suite.Equal("integer", column03.Name)
	suite.Equal(psqldef.PSQLTypeInteger, column03.Type)

	column04 := columns0[4]
	suite.Equal("protected", column04.Name)
	suite.Equal(psqldef.PSQLTypeText, column04.Type)

	column05 := columns0[5]
	suite.Equal("sealed", column05.Name)
	suite.Equal(psqldef.PSQLTypeText, column05.Type)

	column06 := columns0[6]
	suite.Equal("string", column06.Name)
	suite.Equal(psqldef.PSQLTypeText, column06.Type)

	column07 := columns0[7]
	suite.Equal("time", column07.Name)
	suite.Equal(psqldef.PSQLTypeTimestampTZ, column07.Type)

	column08 := columns0[8]
	suite.Equal("uuid", column08.Name)
	suite.Equal(psqldef.PSQLTypeUUID, column08.Type)
	suite.True(column08.PrimaryKey)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_SuccessHook_Failure() {
	var featureFlag = "otherName"
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelSuccess: func(allModelTables []*psqldef.Table) ([]*psqldef.Table, error) {
			if featureFlag != "otherName" {
				return allModelTables, nil
			}
			return nil, fmt.Errorf("compile model success hook error")
		},
	}
	config := suite.getCompileConfig()
	config.ModelHooks = modelHooks

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
	suite.ErrorContains(allTablesErr, "compile model success hook error")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_FailureHook_NoSchema() {
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelFailure: func(config cfg.MorpheConfig, model yaml.Model, compileFailure error) error {
			return fmt.Errorf("Model %s: %w", model.Name, compileFailure)
		},
	}
	config := suite.getCompileConfig()
	config.ModelHooks = modelHooks
	config.MorpheConfig.MorpheModelsConfig.Schema = "" // Clear schema to cause validation error

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
	suite.ErrorContains(allTablesErr, "Model Basic: model schema cannot be empty")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_EnumField() {
	config := suite.getCompileConfig()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Nationality": {
				Type: "Nationality",
			},
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
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

	enum0 := yaml.Enum{
		Name: "Nationality",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"US": "American",
			"DE": "German",
			"FR": "French",
		},
	}

	r := registry.NewRegistry()
	r.SetEnum("Nationality", enum0)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table0 := allTables[0]
	suite.Equal(table0.Name, "basics")

	columns0 := table0.Columns
	suite.Len(columns0, 3)

	column00 := columns0[0]
	suite.Equal(column00.Name, "auto_increment")
	suite.Equal(column00.Type, psqldef.PSQLTypeSerial)

	column01 := columns0[1]
	suite.Equal(column01.Name, "nationality_id")
	suite.Equal(column01.Type, psqldef.PSQLTypeInteger)
	suite.True(column01.NotNull)

	column02 := columns0[2]
	suite.Equal(column02.Name, "uuid")
	suite.Equal(column02.Type, psqldef.PSQLTypeUUID)
	suite.True(column02.PrimaryKey)

	foreignKeys0 := table0.ForeignKeys
	suite.Len(foreignKeys0, 1)

	foreignKey0 := foreignKeys0[0]
	suite.Equal(foreignKey0.Schema, config.MorpheConfig.MorpheModelsConfig.Schema)
	suite.Equal(foreignKey0.Name, "fk_basics_nationality_id")
	suite.Equal(foreignKey0.TableName, "basics")
	suite.Equal(foreignKey0.ColumnNames, []string{"nationality_id"})
	suite.Equal(foreignKey0.RefSchema, "public")
	suite.Equal(foreignKey0.RefTableName, "nationalities")
	suite.Equal(foreignKey0.RefColumnNames, []string{"id"})
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForOnePoly() {
	config := suite.getCompileConfig()

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	articleModel := yaml.Model{
		Name: "Article",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"id":      {Type: yaml.ModelFieldTypeUUID},
			"content": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Post", postModel)
	r.SetModel("Article", articleModel)
	r.SetModel("Comment", commentModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, commentModel)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table := allTables[0]
	suite.Equal("comments", table.Name)

	suite.Len(table.Columns, 4)

	suite.Equal("content", table.Columns[0].Name)
	suite.Equal("id", table.Columns[1].Name)

	suite.Equal("commentable_type", table.Columns[2].Name)
	suite.Equal(psqldef.PSQLTypeText, table.Columns[2].Type)
	suite.True(table.Columns[2].NotNull)
	suite.False(table.Columns[2].PrimaryKey)

	suite.Equal("commentable_id", table.Columns[3].Name)
	suite.Equal(psqldef.PSQLTypeText, table.Columns[3].Type)
	suite.True(table.Columns[3].NotNull)
	suite.False(table.Columns[3].PrimaryKey)

	suite.Len(table.ForeignKeys, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForOnePoly_LongRelationName() {
	config := suite.getCompileConfig()

	userModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	auditLogModel := yaml.Model{
		Name: "AuditLog",
		Fields: map[string]yaml.ModelField{
			"id":     {Type: yaml.ModelFieldTypeUUID},
			"action": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"AuditableResource": {
				Type: "ForOnePoly",
				For:  []string{"User", "Post"},
			},
		},
	}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()
	r.SetModel("User", userModel)
	r.SetModel("Post", postModel)
	r.SetModel("AuditLog", auditLogModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, auditLogModel)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table := allTables[0]
	suite.Equal("audit_logs", table.Name)

	suite.Len(table.Columns, 4)

	suite.Equal("action", table.Columns[0].Name)
	suite.Equal("id", table.Columns[1].Name)

	suite.Equal("auditable_resource_type", table.Columns[2].Name)
	suite.Equal(psqldef.PSQLTypeText, table.Columns[2].Type)
	suite.True(table.Columns[2].NotNull)
	suite.False(table.Columns[2].PrimaryKey)

	suite.Equal("auditable_resource_id", table.Columns[3].Name)
	suite.Equal(psqldef.PSQLTypeText, table.Columns[3].Type)
	suite.True(table.Columns[3].NotNull)
	suite.False(table.Columns[3].PrimaryKey)

	suite.Len(table.ForeignKeys, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForManyPoly() {
	config := suite.getCompileConfig()

	model0 := yaml.Model{
		Name: "Tag",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Name": {
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
			"Taggable": {
				Type: "ForManyPoly",
				For: []string{
					"Post",
					"Product",
				},
			},
		},
	}
	model1 := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Title": {
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
			"Tag": {
				Type:    "HasManyPoly",
				Through: "Taggable",
			},
		},
	}
	model2 := yaml.Model{
		Name: "Product",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Name": {
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
			"Tag": {
				Type:    "HasManyPoly",
				Through: "Taggable",
			},
		},
	}
	r := registry.NewRegistry()
	r.SetModel("Tag", model0)
	r.SetModel("Post", model1)
	r.SetModel("Product", model2)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, model0)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 2)

	table0 := allTables[0]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table0.Schema)
	suite.Equal("tags", table0.Name)

	columns0 := table0.Columns
	suite.Len(columns0, 2)

	columns00 := columns0[0]
	suite.Equal("id", columns00.Name)
	suite.Equal(psqldef.PSQLTypeSerial, columns00.Type)
	suite.False(columns00.NotNull)
	suite.True(columns00.PrimaryKey)
	suite.Equal("", columns00.Default)

	columns01 := columns0[1]
	suite.Equal("name", columns01.Name)
	suite.Equal(psqldef.PSQLTypeText, columns01.Type)
	suite.False(columns01.NotNull)
	suite.False(columns01.PrimaryKey)
	suite.Equal("", columns01.Default)

	suite.Len(table0.ForeignKeys, 0)
	suite.Len(table0.Indices, 0)
	suite.Len(table0.UniqueConstraints, 0)

	// Polymorphic junction table tags <-> taggables
	table1 := allTables[1]

	suite.Equal(config.MorpheConfig.MorpheModelsConfig.Schema, table1.Schema)
	suite.Equal("tag_taggables", table1.Name)

	columns1 := table1.Columns
	suite.Len(columns1, 4)

	columns10 := columns1[0]
	suite.Equal("id", columns10.Name)
	suite.Equal(psqldef.PSQLTypeSerial, columns10.Type)
	suite.False(columns10.NotNull)
	suite.True(columns10.PrimaryKey)
	suite.Equal("", columns10.Default)

	columns11 := columns1[1]
	suite.Equal("tag_id", columns11.Name)
	suite.Equal(psqldef.PSQLTypeInteger, columns11.Type)
	suite.False(columns11.NotNull)
	suite.False(columns11.PrimaryKey)
	suite.Equal("", columns11.Default)

	columns12 := columns1[2]
	suite.Equal("taggable_type", columns12.Name)
	suite.Equal(psqldef.PSQLTypeText, columns12.Type)
	suite.False(columns12.NotNull)
	suite.False(columns12.PrimaryKey)
	suite.Equal("", columns12.Default)

	columns13 := columns1[3]
	suite.Equal("taggable_id", columns13.Name)
	suite.Equal(psqldef.PSQLTypeText, columns13.Type)
	suite.False(columns13.NotNull)
	suite.False(columns13.PrimaryKey)
	suite.Equal("", columns13.Default)

	// Should have one foreign key constraint (for the source model)
	suite.Len(table1.ForeignKeys, 1)
	foreignKey10 := table1.ForeignKeys[0]
	suite.Equal("public", foreignKey10.Schema)
	suite.Equal("fk_tag_taggables_tag_id", foreignKey10.Name)
	suite.Equal("tag_taggables", foreignKey10.TableName)
	suite.Len(foreignKey10.ColumnNames, 1)
	fkColumn10 := foreignKey10.ColumnNames[0]
	suite.Equal("tag_id", fkColumn10)
	suite.Equal("public", foreignKey10.RefSchema)
	suite.Equal("tags", foreignKey10.RefTableName)
	suite.Len(foreignKey10.RefColumnNames, 1)
	fkColumnRef10 := foreignKey10.RefColumnNames[0]
	suite.Equal("id", fkColumnRef10)
	suite.Equal("CASCADE", foreignKey10.OnDelete)
	suite.Equal("", foreignKey10.OnUpdate)

	suite.Len(table1.Indices, 1)
	index10 := table1.Indices[0]
	suite.Equal("idx_tag_taggables_tag_id", index10.Name)
	suite.Equal("tag_taggables", index10.TableName)
	suite.Len(index10.Columns, 1)
	suite.Equal("tag_id", index10.Columns[0])
	suite.False(index10.IsUnique)

	// Should have unique constraint on (source_id, target_type, target_id)
	suite.Len(table1.UniqueConstraints, 1)
	uniqueConstraint10 := table1.UniqueConstraints[0]
	suite.Equal("uk_tag_taggables_tag_id_taggable_type_taggable_id", uniqueConstraint10.Name)
	suite.Equal("tag_taggables", uniqueConstraint10.TableName)
	suite.Len(uniqueConstraint10.ColumnNames, 3)
	suite.Equal("tag_id", uniqueConstraint10.ColumnNames[0])
	suite.Equal("taggable_type", uniqueConstraint10.ColumnNames[1])
	suite.Equal("taggable_id", uniqueConstraint10.ColumnNames[2])
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_HasOnePoly() {
	config := suite.getCompileConfig()

	// Post model with HasOnePoly relationship
	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id":    {Type: yaml.ModelFieldTypeUUID},
			"title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Comment": {
				Type:    "HasOnePoly",
				Through: "Commentable",
			},
		},
	}

	// Comment model with ForOnePoly relationship
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"id":      {Type: yaml.ModelFieldTypeUUID},
			"content": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}

	// Article model (target of the polymorphic relationship)
	articleModel := yaml.Model{
		Name: "Article",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()
	r.SetModel("Post", postModel)
	r.SetModel("Comment", commentModel)
	r.SetModel("Article", articleModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, postModel)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table := allTables[0]
	suite.Equal("posts", table.Name)

	// Should only have the model's own fields, no additional columns for HasOnePoly
	suite.Len(table.Columns, 2)
	suite.Equal("id", table.Columns[0].Name)
	suite.Equal("title", table.Columns[1].Name)

	suite.Len(table.ForeignKeys, 0)
	suite.Len(table.Indices, 0)
	suite.Len(table.UniqueConstraints, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_HasManyPoly() {
	config := suite.getCompileConfig()

	// Post model with HasManyPoly relationship
	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id":    {Type: yaml.ModelFieldTypeUUID},
			"title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Tag": {
				Type:    "HasManyPoly",
				Through: "Taggable",
			},
		},
	}

	// Tag model with ForManyPoly relationship
	tagModel := yaml.Model{
		Name: "Tag",
		Fields: map[string]yaml.ModelField{
			"id":   {Type: yaml.ModelFieldTypeUUID},
			"name": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Taggable": {
				Type: "ForManyPoly",
				For:  []string{"Post", "Product"},
			},
		},
	}

	// Product model (target of the polymorphic relationship)
	productModel := yaml.Model{
		Name: "Product",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()
	r.SetModel("Post", postModel)
	r.SetModel("Tag", tagModel)
	r.SetModel("Product", productModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, postModel)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1)

	table := allTables[0]
	suite.Equal("posts", table.Name)

	// Should only have the model's own fields, no additional columns for HasManyPoly
	suite.Len(table.Columns, 2)
	suite.Equal("id", table.Columns[0].Name)
	suite.Equal("title", table.Columns[1].Name)

	suite.Len(table.ForeignKeys, 0)
	suite.Len(table.Indices, 0)
	suite.Len(table.UniqueConstraints, 0)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForOnePoly_NonExistentForProperty() {
	config := suite.getCompileConfig()

	// Comment model with ForOnePoly relationship referencing non-existent model
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"id":      {Type: yaml.ModelFieldTypeUUID},
			"content": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "NonExistentModel"},
			},
		},
	}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)
	r.SetModel("Post", postModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, commentModel)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "model 'NonExistentModel' referenced in 'for' property not found")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForManyPoly_EmptyForProperty() {
	config := suite.getCompileConfig()

	// Tag model with ForManyPoly relationship with empty for array
	tagModel := yaml.Model{
		Name: "Tag",
		Fields: map[string]yaml.ModelField{
			"id":   {Type: yaml.ModelFieldTypeUUID},
			"name": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Taggable": {
				Type: "ForManyPoly",
				For:  []string{}, // Empty for array
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Tag", tagModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, tagModel)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "polymorphic relation 'Taggable' must have at least one model in 'for' property")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForOnePoly_MissingForProperty() {
	config := suite.getCompileConfig()

	// Comment model with ForOnePoly relationship missing for property
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"id":      {Type: yaml.ModelFieldTypeUUID},
			"content": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				// Missing For property
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, commentModel)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "polymorphic relation 'Commentable' must have at least one model in 'for' property")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_HasOnePoly_InvalidThrough() {
	config := suite.getCompileConfig()

	// Post model with HasOnePoly relationship pointing to non-existent through
	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id":    {Type: yaml.ModelFieldTypeUUID},
			"title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Comment": {
				Type:    "HasOnePoly",
				Through: "NonExistentRelation",
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Post", postModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, postModel)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "through property 'NonExistentRelation' not found")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_HasManyPoly_InvalidThrough() {
	config := suite.getCompileConfig()

	// Post model with HasManyPoly relationship pointing to non-polymorphic through
	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id":    {Type: yaml.ModelFieldTypeUUID},
			"title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Tag": {
				Type:    "HasManyPoly",
				Through: "RegularRelation",
			},
			"RegularRelation": {
				Type: "ForOne",
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Post", postModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, postModel)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "through property 'RegularRelation' must reference a polymorphic relation")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_HasOnePoly_MissingThroughProperty() {
	config := suite.getCompileConfig()

	// Post model with HasOnePoly relationship missing through property
	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id":    {Type: yaml.ModelFieldTypeUUID},
			"title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Comment": {
				Type: "HasOnePoly",
				// Missing Through property
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Post", postModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, postModel)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "polymorphic relation 'Comment' of type 'HasOnePoly' must have a 'through' property")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_CircularPolymorphicReference() {
	config := suite.getCompileConfig()

	// Simple circular reference: A -> B -> A
	modelA := yaml.Model{
		Name: "UserProfile",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"LinkedContent": {
				Type: "ForOnePoly",
				For:  []string{"BlogPost"},
			},
		},
	}

	modelB := yaml.Model{
		Name: "BlogPost",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Author": {
				Type: "ForOnePoly",
				For:  []string{"UserProfile"},
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("UserProfile", modelA)
	r.SetModel("BlogPost", modelB)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, modelA)

	suite.NotNil(allTablesErr)
	suite.ErrorContains(allTablesErr, "circular polymorphic reference detected")
	suite.ErrorContains(allTablesErr, "infinite loop")
	suite.ErrorContains(allTablesErr, "UserProfile -[LinkedContent]-> BlogPost")
	suite.Nil(allTables)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForOne_Aliased() {
	config := suite.getCompileConfig()

	// Company model with aliased ForOne relationships
	companyModel := yaml.Model{
		Name: "Company",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
				Attributes: []string{
					"mandatory",
				},
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
			"TaxID": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"name":    {Fields: []string{"Name"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Person": {
				Type: "ForOne",
			},
		},
	}

	// Person model  
	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
				Attributes: []string{
					"mandatory",
				},
			},
			"FirstName": {
				Type: yaml.ModelFieldTypeString,
			},
			"LastName": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"name":    {Fields: []string{"FirstName", "LastName"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()
	r.SetModel("Company", companyModel)
	r.SetModel("Person", personModel)

	allTables, allTablesErr := compile.MorpheModelToPSQLTables(config, r, companyModel)

	suite.Nil(allTablesErr)
	suite.Len(allTables, 1) // Only Company table

	// Check Company table
	companyTable := allTables[0]
	suite.Equal("companies", companyTable.Name)
	
	// Should have 4 columns: id, name, tax_id, person_id (for Person relationship)
	suite.Len(companyTable.Columns, 4)
	
	// Check that the foreign key column uses the relation name
	personFkColumn := companyTable.Columns[3]
	suite.Equal("person_id", personFkColumn.Name) // Current behavior: uses relation name
	suite.Equal(psqldef.PSQLTypeInteger, personFkColumn.Type)
	
	// Check foreign key constraint
	suite.Len(companyTable.ForeignKeys, 1)
	fk := companyTable.ForeignKeys[0]
	suite.Equal("fk_companies_person_id", fk.Name)
	suite.Equal([]string{"person_id"}, fk.ColumnNames)
	suite.Equal("people", fk.RefTableName) // References the actual target table
	suite.Equal([]string{"id"}, fk.RefColumnNames)
	
	// This test shows CURRENT behavior - after aliasing implementation, 
	// we want to be able to have "Owner" and "Employee" relations both 
	// pointing to "Person" model with different column names
}

// TODO: These tests will be implemented after adding Aliased field support
// func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_ForMany_Aliased() {
// func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_Aliased_MissingTarget() {
// func (suite *CompileModelsTestSuite) TestMorpheModelToPSQLTables_Related_Aliased_WithPolymorphic() {

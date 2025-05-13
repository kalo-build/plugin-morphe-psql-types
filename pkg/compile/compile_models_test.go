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
	suite.Equal("public", index0.Schema)
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
	suite.Equal("public", index10.Schema)
	suite.Equal("basic_basic_parents", index10.TableName)
	suite.Len(index10.Columns, 1)
	suite.Equal("basic_id", index10.Columns[0])
	suite.False(index10.IsUnique)

	index11 := table1.Indices[1]
	suite.Equal("idx_basic_basic_parents_basic_parent_id", index11.Name)
	suite.Equal("public", index11.Schema)
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
	suite.Equal("public", index10.Schema)
	suite.Equal("basic_basic_parents", index10.TableName)
	suite.Len(index10.Columns, 1)
	suite.Equal("basic_id", index10.Columns[0])
	suite.False(index10.IsUnique)

	index11 := table1.Indices[1]
	suite.Equal("idx_basic_basic_parents_basic_parent_id", index11.Name)
	suite.Equal("public", index11.Schema)
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

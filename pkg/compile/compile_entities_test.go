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

type CompileEntitiesTestSuite struct {
	suite.Suite
}

func TestCompileEntitiesTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEntitiesTestSuite))
}

func (suite *CompileEntitiesTestSuite) getMorpheConfig() cfg.MorpheConfig {
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

func (suite *CompileEntitiesTestSuite) getCompileConfig() compile.MorpheCompileConfig {
	return compile.MorpheCompileConfig{
		MorpheConfig: suite.getMorpheConfig(),
		EntityHooks:  hook.CompileMorpheEntity{},
	}
}

func (suite *CompileEntitiesTestSuite) SetupTest() {
	// Setup code if needed
}

func (suite *CompileEntitiesTestSuite) TearDownTest() {
	// Teardown code if needed
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView() {
	config := suite.getCompileConfig()

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Child": {
				Type: "HasOne",
			},
		},
	}
	r.SetModel("User", model0)

	model1 := yaml.Model{
		Name: "Child",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
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
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
			"Time": {
				Type: yaml.ModelFieldTypeTime,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"User": {
				Type: "ForOne",
			},
		},
	}
	r.SetModel("Child", model1)

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
				Attributes: []string{
					"immutable",
				},
			},
			"AutoIncrement": {
				Type: "User.Child.AutoIncrement",
			},
			"Boolean": {
				Type: "User.Child.Boolean",
			},
			"Date": {
				Type: "User.Child.Date",
			},
			"Float": {
				Type: "User.Child.Float",
			},
			"Integer": {
				Type: "User.Child.Integer",
			},
			"String": {
				Type: "User.Child.String",
			},
			"Time": {
				Type: "User.Child.Time",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	r.SetEntity("User", entity0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.Nil(err)
	suite.NotNil(view)
	suite.Equal("public", view.Schema)
	suite.Equal("user_entities", view.Name)

	suite.Equal("public", view.FromSchema)
	suite.Equal("users", view.FromTable)

	suite.Len(view.Columns, 8)

	column0 := view.Columns[0]
	suite.Equal(column0.Name, "auto_increment")
	suite.Equal(column0.Alias, "")
	suite.Equal(column0.SourceRef, "children.auto_increment")

	column1 := view.Columns[1]
	suite.Equal(column1.Name, "boolean")
	suite.Equal(column1.Alias, "")
	suite.Equal(column1.SourceRef, "children.boolean")

	column2 := view.Columns[2]
	suite.Equal(column2.Name, "date")
	suite.Equal(column2.Alias, "")
	suite.Equal(column2.SourceRef, "children.date")

	column3 := view.Columns[3]
	suite.Equal(column3.Name, "float")
	suite.Equal(column3.Alias, "")
	suite.Equal(column3.SourceRef, "children.float")

	column4 := view.Columns[4]
	suite.Equal(column4.Name, "integer")
	suite.Equal(column4.Alias, "")
	suite.Equal(column4.SourceRef, "children.integer")

	column5 := view.Columns[5]
	suite.Equal(column5.Name, "string")
	suite.Equal(column5.Alias, "")
	suite.Equal(column5.SourceRef, "children.string")

	column6 := view.Columns[6]
	suite.Equal(column6.Name, "time")
	suite.Equal(column6.Alias, "")
	suite.Equal(column6.SourceRef, "children.time")

	column7 := view.Columns[7]
	suite.Equal(column7.Name, "uuid")
	suite.Equal(column7.Alias, "")
	suite.Equal(column7.SourceRef, "users.uuid")

	suite.Equal(1, len(view.Joins))
	join := view.Joins[0]
	suite.Equal("public", join.Schema)
	suite.Equal("children", join.Table)
	suite.Equal("children", join.Alias)
	suite.Equal("LEFT", join.Type)

	suite.Len(join.Conditions, 1)
	joinCondition0 := join.Conditions[0]
	suite.Equal("users.uuid", joinCondition0.LeftRef)
	suite.Equal("children.uuid", joinCondition0.RightRef)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_AlternativeSuffix() {
	config := suite.getCompileConfig()
	config.MorpheEntitiesConfig.ViewNameSuffix = "_alt"

	r := registry.NewRegistry()

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
				Attributes: []string{
					"immutable",
				},
			},
			"Name": {
				Type: "User.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	model0 := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Name": {
				Type: "String",
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", model0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.Nil(err)
	suite.NotNil(view)
	suite.Equal("public", view.Schema)
	suite.Equal("user_alt", view.Name)

	suite.Equal("public", view.FromSchema)
	suite.Equal("users", view.FromTable)

	suite.Len(view.Columns, 2)

	column0 := view.Columns[0]
	suite.Equal(column0.Name, "name")
	suite.Equal(column0.Alias, "")
	suite.Equal(column0.SourceRef, "users.name")

	column1 := view.Columns[1]
	suite.Equal(column1.Name, "uuid")
	suite.Equal(column1.Alias, "")
	suite.Equal(column1.SourceRef, "users.uuid")

	suite.Equal(0, len(view.Joins))
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_NoEntityName() {
	config := suite.getCompileConfig()

	r := registry.NewRegistry()

	entity0 := yaml.Entity{
		Name: "",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
				Attributes: []string{
					"immutable",
				},
			},
			"Name": {
				Type: "User.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	model0 := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Name": {
				Type: "String",
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", model0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.ErrorContains(err, "morphe entity has no name")
	suite.Nil(view)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_NoFields() {
	config := suite.getCompileConfig()

	r := registry.NewRegistry()

	entity0 := yaml.Entity{
		Name:   "User",
		Fields: map[string]yaml.EntityField{},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	model0 := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Name": {
				Type: "String",
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", model0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.ErrorContains(err, "morphe entity User has no fields")
	suite.Nil(view)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_NoIdentifiers() {
	config := suite.getCompileConfig()

	r := registry.NewRegistry()

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
				Attributes: []string{
					"immutable",
				},
			},
			"Name": {
				Type: "User.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{},
		Related:     map[string]yaml.EntityRelation{},
	}

	model0 := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Name": {
				Type: "String",
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", model0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.ErrorContains(err, "entity 'User' has no identifiers")
	suite.Nil(view)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_EnumField() {
	config := suite.getCompileConfig()

	r := registry.NewRegistry()

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
				Attributes: []string{
					"immutable",
				},
			},
			"Nationality": {
				Type: "User.Nationality",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	model0 := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Nationality": {
				Type: "Nationality",
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", model0)

	enum0 := yaml.Enum{
		Name: "Nationality",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"US": "American",
			"DE": "German",
			"FR": "French",
		},
	}
	r.SetEnum("Nationality", enum0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.Nil(err)
	suite.NotNil(view)
	suite.Equal("public", view.Schema)
	suite.Equal("user_entities", view.Name)

	suite.Equal("public", view.FromSchema)
	suite.Equal("users", view.FromTable)

	suite.Len(view.Columns, 2)

	column0 := view.Columns[0]
	suite.Equal(column0.Name, "nationality")
	suite.Equal(column0.Alias, "")
	suite.Equal(column0.SourceRef, "users.nationality")

	column1 := view.Columns[1]
	suite.Equal(column1.Name, "uuid")
	suite.Equal(column1.Alias, "")
	suite.Equal(column1.SourceRef, "users.uuid")

	suite.Equal(0, len(view.Joins))
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_NoSchema() {
	r := registry.NewRegistry()

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
	}
	r.SetEntity("User", entity0)

	config := suite.getCompileConfig()
	config.MorpheEntitiesConfig.Schema = ""

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.ErrorContains(err, "schema is required")
	suite.Nil(view)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_StartHook_Successful() {
	var hookCalled = false
	var hooksEntity yaml.Entity

	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntityStart: func(config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error) {
			hookCalled = true
			hooksEntity = entity

			modifiedConfig := config
			modifiedConfig.MorpheEntitiesConfig.ViewNameSuffix = "_hook_modified"

			return modifiedConfig, entity, nil
		},
	}

	config := suite.getCompileConfig()
	config.EntityHooks = entityHooks

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.True(hookCalled)
	suite.Equal(entity0.Name, hooksEntity.Name)

	suite.Nil(err)
	suite.NotNil(view)

	suite.Equal("public", view.Schema)
	suite.Equal("basic_hook_modified", view.Name)

	suite.Equal("public", view.FromSchema)
	suite.Equal("basics", view.FromTable)
	suite.Len(view.Columns, 1)

	column0 := view.Columns[0]
	suite.Equal("id", column0.Name)
	suite.Equal("basics.id", column0.SourceRef)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_StartHook_Failure() {
	var featureFlag = "otherName"
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntityStart: func(config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error) {
			if featureFlag != "otherName" {
				return config, entity, nil
			}
			return config, entity, fmt.Errorf("compile entity start hook error")
		},
	}

	config := suite.getCompileConfig()
	config.EntityHooks = entityHooks

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.NotNil(err)
	suite.ErrorContains(err, "compile entity start hook error")
	suite.Nil(view)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_SuccessHook_Successful() {
	var hookCalled = false

	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntitySuccess: func(view *psqldef.View) (*psqldef.View, error) {
			hookCalled = true

			modifiedView := view.DeepClone()
			modifiedView.Name = view.Name + "_modified"

			return &modifiedView, nil
		},
	}

	config := suite.getCompileConfig()
	config.EntityHooks = entityHooks

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.True(hookCalled)

	suite.Nil(err)
	suite.NotNil(view)

	suite.Equal("public", view.Schema)
	suite.Equal("basic_entities_modified", view.Name)

	suite.Equal("public", view.FromSchema)
	suite.Equal("basics", view.FromTable)
	suite.Len(view.Columns, 1)

	column0 := view.Columns[0]
	suite.Equal("id", column0.Name)
	suite.Equal("basics.id", column0.SourceRef)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_SuccessHook_Failure() {
	var featureFlag = "otherName"
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntitySuccess: func(view *psqldef.View) (*psqldef.View, error) {
			if featureFlag != "otherName" {
				return view, nil
			}
			return nil, fmt.Errorf("compile entity success hook error")
		},
	}

	config := suite.getCompileConfig()
	config.EntityHooks = entityHooks

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.NotNil(err)
	suite.ErrorContains(err, "compile entity success hook error")
	suite.Nil(view)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_FailureHook() {
	var hookCalled = false
	var originalError error

	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntityStart: func(config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error) {
			return config, entity, fmt.Errorf("original error")
		},
		OnCompileMorpheEntityFailure: func(config cfg.MorpheConfig, entity yaml.Entity, err error) error {
			hookCalled = true
			originalError = err
			return fmt.Errorf("enhanced error: %w", err)
		},
	}

	config := suite.getCompileConfig()
	config.EntityHooks = entityHooks

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	view, err := compile.MorpheEntityToPSQLView(config, r, entity0)

	suite.True(hookCalled)
	suite.NotNil(originalError)
	suite.Equal("original error", originalError.Error())

	suite.NotNil(err)
	suite.ErrorContains(err, "enhanced error")
	suite.ErrorContains(err, "original error")

	suite.Nil(view)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_ForOnePoly() {
	config := suite.getCompileConfig()
	r := registry.NewRegistry()

	// Comment entity for a model that has ForOnePoly relationships
	// Test that entity compilation works correctly for models with polymorphic relationships
	commentEntity := yaml.Entity{
		Name: "Comment",
		Fields: map[string]yaml.EntityField{
			"id": {
				Type:       "Comment.id",
				Attributes: []string{"immutable"},
			},
			"content": {
				Type: "Comment.content",
			},
			// Only reference regular fields, not polymorphic columns
			// This tests that entity compilation works for models with polymorphic relationships
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"id"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Comment model with ForOnePoly relationship
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"id": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"content": {
				Type: yaml.ModelFieldTypeString,
			},
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

	// Target models
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

	r.SetModel("Comment", commentModel)
	r.SetModel("Post", postModel)
	r.SetModel("Article", articleModel)

	view, err := compile.MorpheEntityToPSQLView(config, r, commentEntity)

	suite.Nil(err)
	suite.NotNil(view)
	suite.Equal("comment_entities", view.Name)
	suite.Equal("comments", view.FromTable)

	// Should have only regular fields (polymorphic relationship doesn't affect entity views)
	suite.Len(view.Columns, 2)

	suite.Equal("content", view.Columns[0].Name)
	suite.Equal("comments.content", view.Columns[0].SourceRef)

	suite.Equal("id", view.Columns[1].Name)
	suite.Equal("comments.id", view.Columns[1].SourceRef)

	// Should have no joins
	suite.Len(view.Joins, 0)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_ForManyPoly() {
	config := suite.getCompileConfig()
	r := registry.NewRegistry()

	// Tag entity with ForManyPoly relationship
	tagEntity := yaml.Entity{
		Name: "Tag",
		Fields: map[string]yaml.EntityField{
			"id": {
				Type:       "Tag.id",
				Attributes: []string{"immutable"},
			},
			"name": {
				Type: "Tag.name",
			},
			// Don't reference the polymorphic relationship directly to avoid validation errors
			// The test will verify that only regular fields are included
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"id"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Tag model with ForManyPoly relationship
	tagModel := yaml.Model{
		Name: "Tag",
		Fields: map[string]yaml.ModelField{
			"id": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"name": {
				Type: yaml.ModelFieldTypeString,
			},
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

	// Target models
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

	r.SetModel("Tag", tagModel)
	r.SetModel("Post", postModel)
	r.SetModel("Product", productModel)

	view, err := compile.MorpheEntityToPSQLView(config, r, tagEntity)

	suite.Nil(err)
	suite.NotNil(view)
	suite.Equal("tag_entities", view.Name)
	suite.Equal("tags", view.FromTable)

	// ForManyPoly should be simple - only regular fields, no junction table materialization
	// Should have columns: id, name
	suite.Len(view.Columns, 2)

	suite.Equal("id", view.Columns[0].Name)
	suite.Equal("tags.id", view.Columns[0].SourceRef)

	suite.Equal("name", view.Columns[1].Name)
	suite.Equal("tags.name", view.Columns[1].SourceRef)

	// Should have no joins for ForManyPoly relationships
	suite.Len(view.Joins, 0)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_HasOnePoly() {
	config := suite.getCompileConfig()
	r := registry.NewRegistry()

	// Post entity with HasOnePoly relationship
	postEntity := yaml.Entity{
		Name: "Post",
		Fields: map[string]yaml.EntityField{
			"id": {
				Type:       "Post.id",
				Attributes: []string{"immutable"},
			},
			"title": {
				Type: "Post.title",
			},
			// Don't reference the polymorphic relationship directly to avoid validation errors
			// The test will verify that only regular fields are included
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"id"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Post model with HasOnePoly relationship
	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"title": {
				Type: yaml.ModelFieldTypeString,
			},
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

	// Comment model with the forward ForOnePoly relationship
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

	r.SetModel("Post", postModel)
	r.SetModel("Comment", commentModel)
	r.SetModel("Article", articleModel)

	view, err := compile.MorpheEntityToPSQLView(config, r, postEntity)

	suite.Nil(err)
	suite.NotNil(view)
	suite.Equal("post_entities", view.Name)
	suite.Equal("posts", view.FromTable)

	// HasOnePoly should not be materialized - only regular fields
	// Should have columns: id, title
	suite.Len(view.Columns, 2)

	suite.Equal("id", view.Columns[0].Name)
	suite.Equal("posts.id", view.Columns[0].SourceRef)

	suite.Equal("title", view.Columns[1].Name)
	suite.Equal("posts.title", view.Columns[1].SourceRef)

	// Should have no joins for HasOnePoly relationships
	suite.Len(view.Joins, 0)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_HasManyPoly() {
	config := suite.getCompileConfig()
	r := registry.NewRegistry()

	// Post entity with HasManyPoly relationship
	postEntity := yaml.Entity{
		Name: "Post",
		Fields: map[string]yaml.EntityField{
			"id": {
				Type:       "Post.id",
				Attributes: []string{"immutable"},
			},
			"title": {
				Type: "Post.title",
			},
			// Don't reference the polymorphic relationship directly to avoid validation errors
			// The test will verify that only regular fields are included
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"id"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Post model with HasManyPoly relationship
	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"title": {
				Type: yaml.ModelFieldTypeString,
			},
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

	// Tag model with the forward ForManyPoly relationship
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

	r.SetModel("Post", postModel)
	r.SetModel("Tag", tagModel)
	r.SetModel("Product", productModel)

	view, err := compile.MorpheEntityToPSQLView(config, r, postEntity)

	suite.Nil(err)
	suite.NotNil(view)
	suite.Equal("post_entities", view.Name)
	suite.Equal("posts", view.FromTable)

	// HasManyPoly should not be materialized - only regular fields
	// Should have columns: id, title
	suite.Len(view.Columns, 2)

	suite.Equal("id", view.Columns[0].Name)
	suite.Equal("posts.id", view.Columns[0].SourceRef)

	suite.Equal("title", view.Columns[1].Name)
	suite.Equal("posts.title", view.Columns[1].SourceRef)

	// Should have no joins for HasManyPoly relationships
	suite.Len(view.Joins, 0)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_Mixed_Polymorphic_And_Regular() {
	config := suite.getCompileConfig()
	r := registry.NewRegistry()

	// Entity with mixed polymorphic and regular relationships
	mixedEntity := yaml.Entity{
		Name: "Mixed",
		Fields: map[string]yaml.EntityField{
			"id": {
				Type:       "Mixed.id",
				Attributes: []string{"immutable"},
			},
			"name": {
				Type: "Mixed.name",
			},
			"user": {
				Type: "Mixed.User.email", // Regular relationship - should create join
			},
			// Note: Cannot reference polymorphic relationships directly due to entity validation
			// This test verifies that regular relationships work correctly alongside polymorphic ones
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"id"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Mixed model with various relationship types
	mixedModel := yaml.Model{
		Name: "Mixed",
		Fields: map[string]yaml.ModelField{
			"id": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
			"User": {
				Type: "ForOne", // Regular relationship
			},
			"Tag": {
				Type:    "HasManyPoly",
				Through: "Taggable",
			},
		},
	}

	// Supporting models
	userModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"id":    {Type: yaml.ModelFieldTypeUUID},
			"email": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{},
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
				For:  []string{"Mixed", "Product"},
			},
		},
	}

	r.SetModel("Mixed", mixedModel)
	r.SetModel("User", userModel)
	r.SetModel("Post", postModel)
	r.SetModel("Article", articleModel)
	r.SetModel("Tag", tagModel)

	view, err := compile.MorpheEntityToPSQLView(config, r, mixedEntity)

	suite.Nil(err)
	suite.NotNil(view)
	suite.Equal("mixed_entities", view.Name)
	suite.Equal("mixeds", view.FromTable)

	// Should have columns: id, name, user
	// (polymorphic relationships don't interfere with regular relationships)
	suite.Len(view.Columns, 3)

	// Check regular columns
	suite.Equal("id", view.Columns[0].Name)
	suite.Equal("mixeds.id", view.Columns[0].SourceRef)

	suite.Equal("name", view.Columns[1].Name)
	suite.Equal("mixeds.name", view.Columns[1].SourceRef)

	suite.Equal("user", view.Columns[2].Name)
	suite.Equal("users.email", view.Columns[2].SourceRef)

	// Should have one join for the regular User relationship
	suite.Len(view.Joins, 1)

	join := view.Joins[0]
	suite.Equal("LEFT", join.Type)
	suite.Equal("users", join.Table)
	suite.Equal("users", join.Alias)
	suite.Len(join.Conditions, 1)
	suite.Equal("mixeds.id", join.Conditions[0].LeftRef)
	suite.Equal("users.id", join.Conditions[0].RightRef)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_FieldPath_AliasedRelationships() {
	config := suite.getCompileConfig()
	r := registry.NewRegistry()

	// Person entity that references fields through aliased relationships
	personEntity := yaml.Entity{
		Name: "PersonProfile",
		Fields: map[string]yaml.EntityField{
			"id": {
				Type: "Person.ID",
			},
			"name": {
				Type: "Person.Name",
			},
			"workEmail": {
				Type: "Person.WorkContact.Email",
			},
			"workPhone": {
				Type: "Person.WorkContact.Phone",
			},
			"personalEmail": {
				Type: "Person.PersonalContact.Email",
			},
			"personalPhone": {
				Type: "Person.PersonalContact.Phone",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"id"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Person model with aliased relationships to Contact
	personModel := yaml.Model{
		Name: "Person",
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
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
			"PersonalContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
		},
	}

	// Contact model
	contactModel := yaml.Model{
		Name: "Contact",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Email": {
				Type: yaml.ModelFieldTypeString,
			},
			"Phone": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"PersonsAsWork": {
				Type:    "HasMany",
				Aliased: "Person.WorkContact",
			},
			"PersonsAsPersonal": {
				Type:    "HasMany",
				Aliased: "Person.PersonalContact",
			},
		},
	}

	r.SetEntity("PersonProfile", personEntity)
	r.SetModel("Person", personModel)
	r.SetModel("Contact", contactModel)

	view, err := compile.MorpheEntityToPSQLView(config, r, personEntity)

	if err != nil {
		suite.T().Logf("Error compiling entity: %v", err)
	}

	suite.Nil(err)
	suite.NotNil(view)

	// Check view basics
	suite.Equal("public", view.Schema)
	suite.Equal("person_profile_entities", view.Name)
	suite.Equal("public", view.FromSchema)
	suite.Equal("person_profiles", view.FromTable)

	// Check columns
	suite.Len(view.Columns, 6) // id, name, work_email, work_phone, personal_email, personal_phone

	suite.Equal("id", view.Columns[0].Name)
	suite.Equal("person_profiles.id", view.Columns[0].SourceRef)

	suite.Equal("name", view.Columns[1].Name)
	suite.Equal("person_profiles.name", view.Columns[1].SourceRef)

	suite.Equal("personal_email", view.Columns[2].Name)
	suite.Equal("personal_contacts.email", view.Columns[2].SourceRef)

	suite.Equal("personal_phone", view.Columns[3].Name)
	suite.Equal("personal_contacts.phone", view.Columns[3].SourceRef)

	suite.Equal("work_email", view.Columns[4].Name)
	suite.Equal("work_contacts.email", view.Columns[4].SourceRef)

	suite.Equal("work_phone", view.Columns[5].Name)
	suite.Equal("work_contacts.phone", view.Columns[5].SourceRef)

	// Check joins - there should be 2 joins for the aliased relationships
	suite.Len(view.Joins, 2)

	// Join for PersonalContact (processed first alphabetically)
	personalJoin := view.Joins[0]
	suite.Equal("LEFT", personalJoin.Type)
	suite.Equal("public", personalJoin.Schema)
	suite.Equal("personal_contacts", personalJoin.Table)
	suite.Equal("personal_contacts", personalJoin.Alias)
	suite.Len(personalJoin.Conditions, 1)
	suite.Equal("person_profiles.id", personalJoin.Conditions[0].LeftRef)
	suite.Equal("personal_contacts.id", personalJoin.Conditions[0].RightRef)

	// Join for WorkContact (processed second alphabetically)
	workJoin := view.Joins[1]
	suite.Equal("LEFT", workJoin.Type)
	suite.Equal("public", workJoin.Schema)
	suite.Equal("work_contacts", workJoin.Table)
	suite.Equal("work_contacts", workJoin.Alias)
	suite.Len(workJoin.Conditions, 1)
	suite.Equal("person_profiles.id", workJoin.Conditions[0].LeftRef)
	suite.Equal("work_contacts.id", workJoin.Conditions[0].RightRef)
}

// TestMorpheModelToPSQLTables_Aliased_ErrorHandling tests error cases for aliased relationships
func (suite *CompileEntitiesTestSuite) TestMorpheModelToPSQLTables_Aliased_ErrorHandling() {
	config := suite.getCompileConfig()
	r := registry.NewRegistry()

	// Person model with aliased relationship to non-existent model
	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "NonExistentModel", // This should cause an error
			},
		},
	}

	r.SetModel("Person", personModel)

	// Try to compile - should fail
	_, err := compile.MorpheModelToPSQLTables(config, r, personModel)

	suite.NotNil(err)
	suite.Contains(err.Error(), "NonExistentModel")
}

// Test entity views with polymorphic aliased relationships
func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_ForOnePoly_Aliased() {
	config := suite.getCompileConfig()

	// Document model
	documentModel := yaml.Model{
		Name: "Document",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: yaml.ModelFieldTypeAutoIncrement},
			"Title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	// Video model
	videoModel := yaml.Model{
		Name: "Video",
		Fields: map[string]yaml.ModelField{
			"ID":       {Type: yaml.ModelFieldTypeAutoIncrement},
			"Duration": {Type: yaml.ModelFieldTypeInteger},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	// Comment model with ForOnePoly using alias
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Text": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"CommentableResource": {
				Type:    "ForOnePoly",
				For:     []string{"Document", "Video"},
				Aliased: "Resource", // Alias that doesn't map to any model
			},
		},
	}

	// Comment entity
	commentEntity := yaml.Entity{
		Name: "CommentView",
		Fields: map[string]yaml.EntityField{
			"id":   {Type: "Comment.ID"},
			"text": {Type: "Comment.Text"},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		// We can't reference CommentableResource fields due to polymorphic nature
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()
	r.SetModel("Document", documentModel)
	r.SetModel("Video", videoModel)
	r.SetModel("Comment", commentModel)
	r.SetEntity("CommentView", commentEntity)

	view, viewErr := compile.MorpheEntityToPSQLView(config, r, commentEntity)

	suite.Nil(viewErr)
	suite.NotNil(view)
	suite.Equal("comment_view_entities", view.Name)

	// Should have only the explicitly referenced fields (polymorphic columns are not auto-included)
	suite.Len(view.Columns, 2) // id, text

	// Verify the columns
	suite.Equal("id", view.Columns[0].Name)
	suite.Equal("comment_views.id", view.Columns[0].SourceRef)

	suite.Equal("text", view.Columns[1].Name)
	suite.Equal("comment_views.text", view.Columns[1].SourceRef)

	// Should have no joins for polymorphic relationships
	suite.Len(view.Joins, 0)
}

// Test entity views with ForManyPoly aliased relationships
func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_ForManyPoly_Aliased() {
	config := suite.getCompileConfig()

	// Document model
	documentModel := yaml.Model{
		Name: "Document",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	// Project model
	projectModel := yaml.Model{
		Name: "Project",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	// Tag model with ForManyPoly using alias
	tagModel := yaml.Model{
		Name: "Tag",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Name": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"TaggedItems": {
				Type:    "ForManyPoly",
				For:     []string{"Document", "Project"},
				Aliased: "Item", // Alias
			},
		},
	}

	// Tag entity
	tagEntity := yaml.Entity{
		Name: "TagView",
		Fields: map[string]yaml.EntityField{
			"id":   {Type: "Tag.ID"},
			"name": {Type: "Tag.Name"},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()
	r.SetModel("Document", documentModel)
	r.SetModel("Project", projectModel)
	r.SetModel("Tag", tagModel)
	r.SetEntity("TagView", tagEntity)

	view, viewErr := compile.MorpheEntityToPSQLView(config, r, tagEntity)

	suite.Nil(viewErr)
	suite.NotNil(view)
	suite.Equal("tag_view_entities", view.Name)

	// ForManyPoly should be simple - only regular fields, no junction table materialization
	suite.Len(view.Columns, 2) // Only id and name

	// Should have no joins for ForManyPoly relationships
	suite.Len(view.Joins, 0)
}

// Test the polymorphic inverse aliasing pattern in entities
func (suite *CompileEntitiesTestSuite) TestMorpheEntityToPSQLView_HasOnePoly_Aliased() {
	config := suite.getCompileConfig()

	// Comment model with ForOnePoly
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Text": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Task"},
			},
		},
	}

	// Post model with semantic alias
	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: yaml.ModelFieldTypeAutoIncrement},
			"Title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Note": { // Semantic field name
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	// Post entity
	postEntity := yaml.Entity{
		Name: "PostView",
		Fields: map[string]yaml.EntityField{
			"id":    {Type: "Post.ID"},
			"title": {Type: "Post.Title"},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)
	r.SetModel("Post", postModel)
	r.SetEntity("PostView", postEntity)

	view, viewErr := compile.MorpheEntityToPSQLView(config, r, postEntity)

	suite.Nil(viewErr)
	suite.NotNil(view)
	suite.Equal("post_view_entities", view.Name)

	// HasOnePoly should not be materialized - only regular fields
	suite.Len(view.Columns, 2) // Only id and title

	// Should have no joins for HasOnePoly relationships
	suite.Len(view.Joins, 0)
}

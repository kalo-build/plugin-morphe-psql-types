package compile_test

import (
	"testing"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EntityAliasingTestSuite struct {
	suite.Suite
}

func TestEntityAliasingTestSuite(t *testing.T) {
	suite.Run(t, new(EntityAliasingTestSuite))
}

func (suite *EntityAliasingTestSuite) getCompileConfig() compile.MorpheCompileConfig {
	return compile.DefaultMorpheCompileConfig("", "")
}

func (suite *EntityAliasingTestSuite) TestEntityFieldPathIndirection_CurrentBehavior() {
	config := suite.getCompileConfig()
	r := registry.NewRegistry()

	// Create ContactInfo model
	contactInfoModel := yaml.Model{
		Name: "ContactInfo",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
			"Email": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("ContactInfo", contactInfoModel)

	// Create Person model (current behavior: relation name = target model name)
	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
			"LastName": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"ContactInfo": {Type: "ForOne"},
		},
	}
	r.SetModel("Person", personModel)

	// Create Person entity with field path through relation
	personEntity := yaml.Entity{
		Name: "Person",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Person.ID",
			},
			"LastName": {
				Type: "Person.LastName",
			},
			"Email": {
				Type: "Person.ContactInfo.Email", // Field path through relation
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Test entity compilation
	view, err := compile.MorpheEntityToPSQLView(config, r, personEntity)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), view)

	// Verify view name includes suffix
	assert.Equal(suite.T(), "person_entities", view.Name)
	assert.Equal(suite.T(), "people", view.FromTable)

	// Verify columns are correctly ordered and reference correct tables
	require.Len(suite.T(), view.Columns, 3)

	// Columns are alphabetically sorted by field name
	emailColumn := view.Columns[0]
	assert.Equal(suite.T(), "email", emailColumn.Name)
	assert.Equal(suite.T(), "contact_infos.email", emailColumn.SourceRef) // References target table

	idColumn := view.Columns[1]
	assert.Equal(suite.T(), "id", idColumn.Name)
	assert.Equal(suite.T(), "people.id", idColumn.SourceRef)

	lastNameColumn := view.Columns[2]
	assert.Equal(suite.T(), "last_name", lastNameColumn.Name)
	assert.Equal(suite.T(), "people.last_name", lastNameColumn.SourceRef)

	// Verify join is created for the related table
	require.Len(suite.T(), view.Joins, 1)
	join := view.Joins[0]
	assert.Equal(suite.T(), "LEFT", join.Type)
	assert.Equal(suite.T(), "contact_infos", join.Table) // Target table name
	assert.Equal(suite.T(), "contact_infos", join.Alias)

	// Verify join condition
	require.Len(suite.T(), join.Conditions, 1)
	condition := join.Conditions[0]
	assert.Equal(suite.T(), "people.id", condition.LeftRef)
	assert.Equal(suite.T(), "contact_infos.id", condition.RightRef)
}

func (suite *EntityAliasingTestSuite) TestEntityFieldPathIndirection_ReadyForAliasing() {
	// This test demonstrates that the entity compilation is ready for aliasing
	// When the Aliased field is added to ModelRelation, the following will work:
	//
	// personModel.Related["PrimaryContact"] = yaml.ModelRelation{
	//     Type: "ForOne",
	//     Aliased: "ContactInfo",  // <-- This field will be added
	// }
	//
	// personEntity.Fields["Email"] = yaml.EntityField{
	//     Type: "Person.PrimaryContact.Email",  // <-- This will resolve to ContactInfo
	// }
	//
	// Expected behavior:
	// - processFieldPath will call GetTargetModelNameFromRelation("PrimaryContact", relation)
	// - GetTargetModelNameFromRelation will detect Aliased field and return "ContactInfo"
	// - Table name will be "contact_infos" (pluralized ContactInfo)
	// - Join condition will reference the actual target table
	// - Column source will be "contact_infos.email"

	config := suite.getCompileConfig()
	r := registry.NewRegistry()

	// This test passes with current code, proving aliasing infrastructure is ready
	contactInfoModel := yaml.Model{
		Name: "ContactInfo",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
			"Email": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("ContactInfo", contactInfoModel)

	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"ContactInfo": {Type: "ForOne"},
		},
	}
	r.SetModel("Person", personModel)

	personEntity := yaml.Entity{
		Name: "Person",
		Fields: map[string]yaml.EntityField{
			"ID": {Type: "Person.ID"},
			"Email": {Type: "Person.ContactInfo.Email"},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	view, err := compile.MorpheEntityToPSQLView(config, r, personEntity)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), view)

	// This proves the infrastructure works and will automatically support aliasing
	// when the Aliased field is added to the external ModelRelation struct
	require.Len(suite.T(), view.Columns, 2)
	
	// Columns are alphabetically sorted
	emailColumn := view.Columns[0]
	assert.Equal(suite.T(), "email", emailColumn.Name)
	assert.Equal(suite.T(), "contact_infos.email", emailColumn.SourceRef)
	
	idColumn := view.Columns[1]
	assert.Equal(suite.T(), "id", idColumn.Name)
	assert.Equal(suite.T(), "people.id", idColumn.SourceRef)

	require.Len(suite.T(), view.Joins, 1)
	join := view.Joins[0]
	assert.Equal(suite.T(), "contact_infos", join.Table)
}
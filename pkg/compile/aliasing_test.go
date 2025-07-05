package compile

import (
	"testing"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/stretchr/testify/assert"
)

func TestGetTargetModelNameFromRelation_WithoutAlias(t *testing.T) {
	// Test current behavior when no aliasing is present
	relation := yaml.ModelRelation{
		Type: "ForOne",
	}
	
	result := GetTargetModelNameFromRelation("Person", relation)
	assert.Equal(t, "Person", result, "Should return relation name when no alias is present")
}

func TestGetTargetModelNameFromRelation_WithEmptyAlias(t *testing.T) {
	// Test behavior when alias field exists but is empty
	// Note: This test will pass once the Aliased field is added to ModelRelation
	relation := yaml.ModelRelation{
		Type: "ForOne",
		// If Aliased field exists and is empty, should fallback to relation name
	}
	
	result := GetTargetModelNameFromRelation("Owner", relation)
	assert.Equal(t, "Owner", result, "Should return relation name when alias is empty")
}

func TestValidateAliasedRelations_ValidRelations(t *testing.T) {
	// Create a registry with test models
	r := registry.NewRegistry()
	
	// Add Person model to registry
	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Person", personModel)
	
	// Create Company model with non-aliased relations
	companyModel := yaml.Model{
		Name: "Company",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Person": {Type: "ForOne"},
		},
	}
	
	err := ValidateAliasedRelations(r, companyModel)
	assert.NoError(t, err, "Should not return error for valid non-aliased relations")
}

func TestValidateAliasedRelations_MissingTarget(t *testing.T) {
	// Create a registry without the target model
	r := registry.NewRegistry()
	
	// Create Company model that tries to reference non-existent model
	// Note: This test currently passes because ValidateAliasedRelations only validates
	// when aliasing is actually used (relation name != target name).
	// Once the Aliased field is added, this test should be updated.
	companyModel := yaml.Model{
		Name: "Company",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"NonExistentModel": {Type: "ForOne"},
		},
	}
	
	err := ValidateAliasedRelations(r, companyModel)
	// Currently this will NOT error because no aliasing is detected
	// This validates the current behavior - aliasing validation only occurs when aliasing is used
	assert.NoError(t, err, "Should not error for non-aliased relations (current behavior)")
}

func TestGetForeignKeyColumnNameWithAlias_CurrentBehavior(t *testing.T) {
	relation := yaml.ModelRelation{
		Type: "ForOne",
	}
	
	result := GetForeignKeyColumnNameWithAlias("Person", relation, "ID")
	expected := GetForeignKeyColumnName("Person", "ID") // Should be same as current behavior
	assert.Equal(t, expected, result, "Should match current behavior when no alias is present")
}

func TestGetJunctionTableNameWithAlias_CurrentBehavior(t *testing.T) {
	relation := yaml.ModelRelation{
		Type: "ForMany",
	}
	
	result := GetJunctionTableNameWithAlias("Company", "Person", relation)
	expected := GetJunctionTableName("Company", "Person") // Should be same as current behavior
	assert.Equal(t, expected, result, "Should match current behavior when no alias is present")
}

// Integration test to ensure the current system still works
func TestCompileWithoutAliasing_Integration(t *testing.T) {
	config := DefaultMorpheCompileConfig("", "")
	
	// Create registry with test models
	r := registry.NewRegistry()
	
	// Person model
	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Person", personModel)
	
	// Company model with ForOne relationship to Person
	companyModel := yaml.Model{
		Name: "Company",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Person": {Type: "ForOne"},
		},
	}
	r.SetModel("Company", companyModel)
	
	// Test compilation
	tables, err := MorpheModelToPSQLTables(config, r, companyModel)
	assert.NoError(t, err, "Should compile successfully without aliasing")
	assert.Len(t, tables, 1, "Should generate one table for Company")
	
	// Verify the table structure
	companyTable := tables[0]
	assert.Equal(t, "companies", companyTable.Name)
	
	// Should have foreign key column
	assert.Len(t, companyTable.Columns, 2, "Should have ID and person_id columns")
	
	personFkColumn := companyTable.Columns[1] // Second column should be the FK
	assert.Equal(t, "person_id", personFkColumn.Name)
	
	// Should have foreign key constraint
	assert.Len(t, companyTable.ForeignKeys, 1)
	fk := companyTable.ForeignKeys[0]
	assert.Equal(t, "people", fk.RefTableName) // Should reference people table
}
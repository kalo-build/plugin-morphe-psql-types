package compile

import (
	"fmt"
	"reflect"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
)

// GetTargetModelNameFromRelation returns the target model name for a relation.
// If the relation has an "Aliased" field, it returns that value.
// Otherwise, it returns the relation name (current behavior).
func GetTargetModelNameFromRelation(relationName string, relation yaml.ModelRelation) string {
	// Use reflection to check if the ModelRelation has an "Aliased" field
	relationValue := reflect.ValueOf(relation)
	aliasedField := relationValue.FieldByName("Aliased")
	
	if aliasedField.IsValid() && aliasedField.Kind() == reflect.String {
		aliasedValue := aliasedField.String()
		if aliasedValue != "" {
			return aliasedValue
		}
	}
	
	// Fallback to the relation name (current behavior)
	return relationName
}

// ValidateAliasedRelations validates that all aliased target models exist in the registry
func ValidateAliasedRelations(r *registry.Registry, model yaml.Model) error {
	for relationName, relation := range model.Related {
		targetModelName := GetTargetModelNameFromRelation(relationName, relation)
		
		// If the target model name is different from the relation name, validate it exists
		if targetModelName != relationName {
			_, err := r.GetModel(targetModelName)
			if err != nil {
				return fmt.Errorf("aliased target model '%s' for relation '%s' not found: %w", targetModelName, relationName, err)
			}
		}
	}
	return nil
}

// GetForeignKeyColumnNameWithAlias generates a column name for a foreign key using aliased targets
func GetForeignKeyColumnNameWithAlias(relationName string, relation yaml.ModelRelation, targetPrimaryFieldName string) string {
	targetModelName := GetTargetModelNameFromRelation(relationName, relation)
	return GetForeignKeyColumnName(targetModelName, targetPrimaryFieldName)
}

// GetJunctionTableNameWithAlias generates a junction table name using aliased targets
func GetJunctionTableNameWithAlias(sourceModelName, relationName string, relation yaml.ModelRelation) string {
	targetModelName := GetTargetModelNameFromRelation(relationName, relation)
	return GetJunctionTableName(sourceModelName, targetModelName)
}
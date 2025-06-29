package compile

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/kalo-build/go-util/strcase"
)

// PostgreSQL identifier length limit
const maxIdentifierLength = 63

// Client for pluralization, initialized once
var pluralizeClient = pluralize.NewClient()

// AbbreviateIdentifier shortens an identifier if it exceeds PostgreSQL's 63-character limit
func AbbreviateIdentifier(identifier string, useHash bool) string {
	// If the identifier is already within limits, return it as is
	if len(identifier) <= maxIdentifierLength {
		return identifier
	}

	// Split the identifier into parts (e.g., by underscores for snake_case)
	parts := strings.Split(identifier, "_")
	abbreviated := make([]string, len(parts))

	// For non-prefixes (like "fk", "uk", "idx"), keep them as is
	// These are typically important for identifying the type of object
	prefixParts := 0
	if len(parts) > 0 && (parts[0] == "fk" || parts[0] == "uk" || parts[0] == "idx") {
		abbreviated[0] = parts[0]
		prefixParts = 1
	}

	// Abbreviate each remaining part to its first two letters
	// If the part is already short (1-2 chars), keep it as is
	for i := prefixParts; i < len(parts); i++ {
		if len(parts[i]) <= 2 {
			abbreviated[i] = parts[i]
		} else {
			abbreviated[i] = parts[i][:2]
		}
	}

	// Join the abbreviated parts
	result := strings.Join(abbreviated, "_")

	// If we still exceed the limit and hash is requested, add a short hash
	if len(result) > maxIdentifierLength && useHash {
		// Calculate how much space we have left for the main identifier
		// allowing for an underscore and 8 character hash
		mainLength := maxIdentifierLength - 9 // underscore + 8 chars for hash

		// Truncate the main part and add a hash of the original to avoid collisions
		hash := fmt.Sprintf("%x", md5.Sum([]byte(identifier)))[:8]
		result = result[:mainLength] + "_" + hash
	} else if len(result) > maxIdentifierLength {
		// If we're still over the limit and not using hash, just truncate
		result = result[:maxIdentifierLength]
	}

	return result
}

// GetTableNameFromModel returns the snake_case, pluralized table name for a model
func GetTableNameFromModel(modelName string) string {
	tableName := Pluralize(strcase.ToSnakeCaseLower(modelName))
	return AbbreviateIdentifier(tableName, false)
}

// GetColumnNameFromField returns the snake_case column name for a field
func GetColumnNameFromField(fieldName string) string {
	columnName := strcase.ToSnakeCaseLower(fieldName)
	return AbbreviateIdentifier(columnName, false)
}

// GetForeignKeyColumnName generates a column name for a foreign key
func GetForeignKeyColumnName(relatedModelName, relatedFieldName string) string {
	columnName := fmt.Sprintf("%s_%s",
		strcase.ToSnakeCaseLower(relatedModelName),
		strcase.ToSnakeCaseLower(relatedFieldName))
	return AbbreviateIdentifier(columnName, false)
}

// GetForeignKeyConstraintName generates a name for a foreign key constraint
func GetForeignKeyConstraintName(tableName, columnName string) string {
	constraintName := fmt.Sprintf("fk_%s_%s",
		strcase.ToSnakeCaseLower(tableName),
		columnName)
	return AbbreviateIdentifier(constraintName, true)
}

// GetJunctionTableForeignKeyConstraintName generates a name for a junction table foreign key constraint
func GetJunctionTableForeignKeyConstraintName(junctionTableName, modelName, idFieldName string) string {
	constraintName := fmt.Sprintf("fk_%s_%s_%s",
		strcase.ToSnakeCaseLower(junctionTableName),
		strcase.ToSnakeCaseLower(modelName),
		strcase.ToSnakeCaseLower(idFieldName))
	return AbbreviateIdentifier(constraintName, true)
}

// GetIndexName generates a name for an index
func GetIndexName(tableName, columnName string) string {
	indexName := fmt.Sprintf("idx_%s_%s", tableName, columnName)
	return AbbreviateIdentifier(indexName, true)
}

// GetUniqueConstraintName generates a name for a unique constraint
func GetUniqueConstraintName(tableName string, columnNames ...string) string {
	parts := []string{tableName}
	parts = append(parts, columnNames...)
	constraintName := fmt.Sprintf("uk_%s", strings.Join(parts, "_"))
	return AbbreviateIdentifier(constraintName, true)
}

// GetJunctionTableName generates a name for a junction table
func GetJunctionTableName(sourceModelName, targetModelName string) string {
	// Generate the singular form of the junction table name
	tableName := fmt.Sprintf("%s_%s",
		strcase.ToSnakeCaseLower(sourceModelName),
		strcase.ToSnakeCaseLower(targetModelName))

	// Return the pluralized form
	tableName = Pluralize(tableName)
	return AbbreviateIdentifier(tableName, false)
}

// GetJunctionTableUniqueConstraintName generates a name for a junction table unique constraint
func GetJunctionTableUniqueConstraintName(
	junctionTableName string,
	model1Name, model1IdName string,
	model2Name, model2IdName string,
) string {
	constraintName := fmt.Sprintf("uk_%s_%s_%s_%s_%s",
		strcase.ToSnakeCaseLower(junctionTableName),
		strcase.ToSnakeCaseLower(model1Name),
		strcase.ToSnakeCaseLower(model1IdName),
		strcase.ToSnakeCaseLower(model2Name),
		strcase.ToSnakeCaseLower(model2IdName))
	return AbbreviateIdentifier(constraintName, true)
}

// GetPolymorphicJunctionTableUniqueConstraintName generates a name for a polymorphic junction table unique constraint
func GetPolymorphicJunctionTableUniqueConstraintName(
	junctionTableName string,
	sourceName, sourceIdName string,
	relationName string,
) string {
	constraintName := fmt.Sprintf("uk_%s_%s_%s_%s_type_%s_id",
		strcase.ToSnakeCaseLower(junctionTableName),
		strcase.ToSnakeCaseLower(sourceName),
		strcase.ToSnakeCaseLower(sourceIdName),
		strcase.ToSnakeCaseLower(relationName),
		strcase.ToSnakeCaseLower(relationName))
	return AbbreviateIdentifier(constraintName, true)
}

// Pluralize a word using simple English rules
func Pluralize(word string) string {
	return pluralizeClient.Plural(word)
}

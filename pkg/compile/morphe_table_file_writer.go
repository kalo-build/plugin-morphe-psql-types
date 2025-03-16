package compile

import (
	"fmt"
	"strings"

	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/psqldef"
	"github.com/kaloseia/plugin-morphe-psql-types/pkg/sqlfile"
)

type MorpheTableFileWriter struct {
	Type          MorpheTableType
	TargetDirPath string
}

func (w *MorpheTableFileWriter) WriteTable(tableDefinition *psqldef.Table) ([]byte, error) {
	allTableLines, allLinesErr := w.getAllTableLines(tableDefinition)
	if allLinesErr != nil {
		return nil, allLinesErr
	}

	tableFileContents, tableContentsErr := core.LinesToString(allTableLines)
	if tableContentsErr != nil {
		return nil, tableContentsErr
	}

	return sqlfile.WriteSQLDefinitionFile(w.TargetDirPath, tableDefinition.Name, tableFileContents)
}

func (w *MorpheTableFileWriter) getAllTableLines(tableDefinition *psqldef.Table) ([]string, error) {
	allTableLines := []string{}

	// Add header comment
	allTableLines = append(allTableLines, fmt.Sprintf("-- Table definition for %s", tableDefinition.Name))
	allTableLines = append(allTableLines, "")

	// Create schema if specified
	if tableDefinition.Schema != "" {
		allTableLines = append(allTableLines, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", tableDefinition.Schema))
		allTableLines = append(allTableLines, "")
	}

	// Create table
	tableLines, tableErr := w.getCreateTableLines(tableDefinition)
	if tableErr != nil {
		return nil, tableErr
	}
	allTableLines = append(allTableLines, tableLines...)
	allTableLines = append(allTableLines, "")

	// Add indices
	if len(tableDefinition.Indices) > 0 {
		indexLines, indexErr := w.getIndexLines(tableDefinition)
		if indexErr != nil {
			return nil, indexErr
		}
		allTableLines = append(allTableLines, indexLines...)
		allTableLines = append(allTableLines, "")
	}

	// Add seed data
	if len(tableDefinition.SeedData) > 0 {
		seedDataLines, seedErr := w.getSeedDataLines(tableDefinition)
		if seedErr != nil {
			return nil, seedErr
		}
		allTableLines = append(allTableLines, seedDataLines...)
		allTableLines = append(allTableLines, "")
	}

	return allTableLines, nil
}

func (w *MorpheTableFileWriter) getCreateTableLines(tableDefinition *psqldef.Table) ([]string, error) {
	tableName := tableDefinition.Name
	if tableDefinition.Schema != "" {
		tableName = tableDefinition.Schema + "." + tableName
	}

	tableLines := []string{
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName),
	}

	// Add columns
	for colIdx, column := range tableDefinition.Columns {
		columnDef := w.formatColumnDefinition(column)

		// Add comma if not the last column or if we have constraints to add
		if colIdx < len(tableDefinition.Columns)-1 ||
			len(tableDefinition.ForeignKeys) > 0 ||
			len(tableDefinition.UniqueConstraints) > 0 {
			columnDef += ","
		}

		tableLines = append(tableLines, "\t"+columnDef)
	}

	// Add unique constraints
	for uqIdx, uniqueConstraint := range tableDefinition.UniqueConstraints {
		constraintLine := fmt.Sprintf("\tUNIQUE (%s)", strings.Join(uniqueConstraint.ColumnNames, ", "))

		// Add comma if not the last constraint or if we have foreign keys to add
		if uqIdx < len(tableDefinition.UniqueConstraints)-1 || len(tableDefinition.ForeignKeys) > 0 {
			constraintLine += ","
		}

		tableLines = append(tableLines, constraintLine)
	}

	// Add foreign key constraints with proper formatting
	for fkIdx, foreignKey := range tableDefinition.ForeignKeys {
		// Format according to the spec with named constraints
		if foreignKey.Name != "" {
			// Format with CONSTRAINT and multiline for readability
			fkLine := fmt.Sprintf("\tCONSTRAINT %s FOREIGN KEY (%s)",
				foreignKey.Name,
				strings.Join(foreignKey.ColumnNames, ", "))
			tableLines = append(tableLines, fkLine)

			refLine := fmt.Sprintf("\t\tREFERENCES %s(%s)",
				foreignKey.RefTableName,
				strings.Join(foreignKey.RefColumnNames, ", "))

			if foreignKey.OnDelete != "" {
				refLine += fmt.Sprintf("\n\t\tON DELETE %s", foreignKey.OnDelete)
			}

			// Only add comma if not the last foreign key
			if fkIdx < len(tableDefinition.ForeignKeys)-1 {
				refLine += ","
			}

			tableLines = append(tableLines, refLine)
		} else {
			// Fallback to simple single-line format for unnamed constraints
			fkLine := fmt.Sprintf("\tFOREIGN KEY (%s) REFERENCES %s (%s)",
				strings.Join(foreignKey.ColumnNames, ", "),
				foreignKey.RefTableName,
				strings.Join(foreignKey.RefColumnNames, ", "))

			// Only add comma if not the last foreign key
			if fkIdx < len(tableDefinition.ForeignKeys)-1 {
				fkLine += ","
			}

			tableLines = append(tableLines, fkLine)
		}
	}

	tableLines = append(tableLines, ");")
	return tableLines, nil
}

func (w *MorpheTableFileWriter) formatColumnDefinition(column psqldef.Column) string {
	parts := []string{column.Name, column.Type.GetSyntax()}

	if column.NotNull {
		parts = append(parts, "NOT NULL")
	}

	if column.PrimaryKey {
		parts = append(parts, "PRIMARY KEY")
	}

	if column.Default != "" {
		parts = append(parts, "DEFAULT", column.Default)
	}

	return strings.Join(parts, " ")
}

func (w *MorpheTableFileWriter) getIndexLines(tableDefinition *psqldef.Table) ([]string, error) {
	indexLines := []string{
		"-- Indices",
	}

	tableName := tableDefinition.Name
	if tableDefinition.Schema != "" {
		tableName = tableDefinition.Schema + "." + tableName
	}

	for _, index := range tableDefinition.Indices {
		indexName := index.Name
		if indexName == "" {
			indexName = fmt.Sprintf("idx_%s_%s", tableDefinition.Name, strings.Join(index.Columns, "_"))
		}

		indexType := ""
		if index.Using != "" {
			indexType = "USING " + index.Using + " "
		}

		unique := ""
		if index.IsUnique {
			unique = "UNIQUE "
		}

		indexLine := fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s %s(%s);",
			unique, indexName, tableName, indexType, strings.Join(index.Columns, ", "))

		indexLines = append(indexLines, indexLine)
	}

	return indexLines, nil
}

func (w *MorpheTableFileWriter) getSeedDataLines(tableDefinition *psqldef.Table) ([]string, error) {
	seedDataLines := []string{
		"-- Seed Data",
	}

	// Create a column map for quick lookups
	columnMap := make(map[string]psqldef.Column)
	for _, col := range tableDefinition.Columns {
		columnMap[col.Name] = col
	}

	for _, insertStmt := range tableDefinition.SeedData {
		tableName := insertStmt.TableName
		if insertStmt.Schema != "" {
			tableName = insertStmt.Schema + "." + tableName
		}

		// Validate table name matches
		if tableDefinition.Name != insertStmt.TableName {
			return nil, fmt.Errorf("seed data refers to table '%s', but expected '%s'",
				insertStmt.TableName, tableDefinition.Name)
		}

		columnList := strings.Join(insertStmt.Columns, ", ")
		for rowIdx, valueRow := range insertStmt.Values {
			// Validate row length matches column count
			if len(valueRow) != len(insertStmt.Columns) {
				return nil, fmt.Errorf("row %d has %d values but expected %d columns",
					rowIdx, len(valueRow), len(insertStmt.Columns))
			}

			formattedValues := make([]string, len(valueRow))

			for rowIdx, val := range valueRow {
				colName := insertStmt.Columns[rowIdx]

				// Check if column exists in table definition
				col, exists := columnMap[colName]
				if !exists {
					return nil, fmt.Errorf("column '%s' in seed data not found in table definition", colName)
				}

				// Validate value type against column type
				if err := w.validateValueType(val, col); err != nil {
					return nil, fmt.Errorf("invalid value for column '%s' (row %d): %v", colName, rowIdx, err)
				}

				formattedValues[rowIdx] = w.formatSQLValue(val, col.Type)
			}

			valueList := strings.Join(formattedValues, ", ")
			insertLine := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
				tableName, columnList, valueList)

			seedDataLines = append(seedDataLines, insertLine)
		}
	}

	return seedDataLines, nil
}

// validateValueType checks if a value is compatible with the column type
func (w *MorpheTableFileWriter) validateValueType(value any, column psqldef.Column) error {
	if value == nil {
		if column.NotNull {
			return fmt.Errorf("NULL value not allowed for NOT NULL column")
		}
		return nil
	}

	switch column.Type.GetSyntax() {
	case "boolean", "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean value")
		}
	case "integer", "int", "int4", "smallint", "int2", "bigint", "int8":
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return nil
		default:
			return fmt.Errorf("expected integer value")
		}
	case "real", "float4", "double precision", "float8", "numeric", "decimal":
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			return nil
		default:
			return fmt.Errorf("expected numeric value")
		}
	case "text", "varchar", "char", "character", "character varying":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string value")
		}
	}

	return nil
}

// formatSQLValue formats a value for SQL, taking into account the column type
func (w *MorpheTableFileWriter) formatSQLValue(value any, columnType psqldef.PSQLType) string {
	if value == nil {
		return "NULL"
	}

	// Handle special PostgreSQL type formatting
	typeSyntax := columnType.GetSyntax()

	// Special case handling for certain PostgreSQL types
	if strings.HasPrefix(typeSyntax, "timestamp") ||
		strings.HasPrefix(typeSyntax, "date") ||
		strings.HasPrefix(typeSyntax, "time") {
		if str, ok := value.(string); ok {
			return fmt.Sprintf("'%s'::timestamptz", strings.ReplaceAll(str, "'", "''"))
		}
	}

	if strings.HasPrefix(typeSyntax, "uuid") {
		if str, ok := value.(string); ok {
			return fmt.Sprintf("'%s'::uuid", strings.ReplaceAll(str, "'", "''"))
		}
	}

	if strings.HasPrefix(typeSyntax, "json") || strings.HasPrefix(typeSyntax, "jsonb") {
		if str, ok := value.(string); ok {
			return fmt.Sprintf("'%s'::jsonb", strings.ReplaceAll(str, "'", "''"))
		}
	}

	// Default formatting by Go type
	switch v := value.(type) {
	case string:
		// Escape single quotes for SQL strings
		escaped := strings.ReplaceAll(v, "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		// For complex types, try to convert to string
		return fmt.Sprintf("'%v'", v)
	}
}

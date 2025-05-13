package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/sqlfile"
)

type MorpheViewFileWriter struct {
	TargetDirPath string
}

func (w *MorpheViewFileWriter) WriteView(viewDefinition *psqldef.View) ([]byte, error) {
	allViewLines, allLinesErr := w.getAllViewLines(viewDefinition)
	if allLinesErr != nil {
		return nil, allLinesErr
	}

	viewFileContents, viewContentsErr := core.LinesToString(allViewLines)
	if viewContentsErr != nil {
		return nil, viewContentsErr
	}

	return sqlfile.WriteSQLDefinitionFile(w.TargetDirPath, viewDefinition.Name, viewFileContents)
}

func (w *MorpheViewFileWriter) getAllViewLines(viewDefinition *psqldef.View) ([]string, error) {
	allViewLines := []string{}

	// Add header comment
	allViewLines = append(allViewLines, fmt.Sprintf("-- View definition for %s", viewDefinition.Name))
	allViewLines = append(allViewLines, "")

	// Create schema if specified
	if viewDefinition.Schema != "" {
		allViewLines = append(allViewLines, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", viewDefinition.Schema))
		allViewLines = append(allViewLines, "")
	}

	// Create view
	viewLines, viewErr := w.getCreateViewLines(viewDefinition)
	if viewErr != nil {
		return nil, viewErr
	}
	allViewLines = append(allViewLines, viewLines...)
	allViewLines = append(allViewLines, "")

	return allViewLines, nil
}

func (w *MorpheViewFileWriter) getCreateViewLines(viewDefinition *psqldef.View) ([]string, error) {
	if len(viewDefinition.Columns) == 0 {
		return nil, fmt.Errorf("view has no columns")
	}

	viewName := viewDefinition.Name
	if viewDefinition.Schema != "" {
		viewName = viewDefinition.Schema + "." + viewName
	}

	viewLines := []string{
		fmt.Sprintf("CREATE OR REPLACE VIEW %s AS", viewName),
		"SELECT",
	}

	columnRefs := []string{}
	for _, column := range viewDefinition.Columns {
		columnRef := column.SourceRef
		if column.Alias != "" {
			columnRef += fmt.Sprintf(" AS %s", column.Name)
		} else {
			parts := strings.Split(column.SourceRef, ".")
			if len(parts) > 1 && parts[len(parts)-1] != column.Name {
				columnRef += fmt.Sprintf(" AS %s", column.Name)
			}
		}
		columnRefs = append(columnRefs, "\t"+columnRef)
	}

	viewLines = append(viewLines, strings.Join(columnRefs, ",\n"))

	if viewDefinition.FromTable == "" {
		return nil, fmt.Errorf("view has no source table")
	}

	fromSchema := viewDefinition.FromSchema
	if fromSchema == "" {
		fromSchema = "public"
	}

	fromTable := viewDefinition.FromTable

	viewLines = append(viewLines, fmt.Sprintf("FROM %s.%s", fromSchema, fromTable))

	for _, join := range viewDefinition.Joins {
		joinTable := join.Table
		if join.Alias != "" && join.Alias != join.Table {
			joinTable += " AS " + join.Alias
		}
		joinSchema := join.Schema
		if joinSchema == "" {
			joinSchema = "public"
		}

		joinLine := fmt.Sprintf("%s JOIN %s.%s", join.Type, joinSchema, joinTable)
		viewLines = append(viewLines, joinLine)

		if len(join.Conditions) > 0 {
			conditions := []string{}
			for _, condition := range join.Conditions {
				conditions = append(conditions, fmt.Sprintf("%s = %s", condition.LeftRef, condition.RightRef))
			}
			viewLines = append(viewLines, "\tON "+strings.Join(conditions, " AND "))
		}
	}

	if viewDefinition.WhereClause != "" {
		viewLines = append(viewLines, fmt.Sprintf("WHERE %s", viewDefinition.WhereClause))
	}

	viewLines[len(viewLines)-1] += ";"

	return viewLines, nil
}

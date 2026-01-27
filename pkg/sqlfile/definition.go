package sqlfile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/go-util/strcase"
)

// WriteSQLDefinitionFile writes a SQL definition file without an order prefix.
// Deprecated: Use WriteSQLDefinitionFileWithOrder for ordered migrations.
func WriteSQLDefinitionFile(dirPath string, definitionName string, psqlFileContents string) ([]byte, error) {
	return WriteSQLDefinitionFileWithOrder(dirPath, definitionName, psqlFileContents, 0)
}

// WriteSQLDefinitionFileWithOrder writes a SQL definition file with an optional order prefix.
// If order is 0, no prefix is added. Otherwise, the file is named like "001_table_name.sql".
func WriteSQLDefinitionFileWithOrder(dirPath string, definitionName string, psqlFileContents string, order int) ([]byte, error) {
	definitionFileName := strcase.ToSnakeCaseLower(definitionName)

	// Add order prefix if order > 0
	if order > 0 {
		definitionFileName = fmt.Sprintf("%03d_%s", order, definitionFileName)
	}

	definitionFilePath := filepath.Join(dirPath, definitionFileName+".sql")
	if _, readErr := os.ReadDir(dirPath); readErr != nil && os.IsNotExist(readErr) {
		mkDirErr := os.MkdirAll(dirPath, 0644)
		if mkDirErr != nil {
			return nil, mkDirErr
		}
	}
	return []byte(psqlFileContents), os.WriteFile(definitionFilePath, []byte(psqlFileContents), 0644)
}

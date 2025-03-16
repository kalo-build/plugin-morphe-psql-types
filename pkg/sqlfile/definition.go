package sqlfile

import (
	"os"
	"path/filepath"

	"github.com/kaloseia/go-util/strcase"
)

func WriteSQLDefinitionFile(dirPath string, definitionName string, psqlFileContents string) ([]byte, error) {
	definitionFileName := strcase.ToSnakeCaseLower(definitionName)
	definitionFilePath := filepath.Join(dirPath, definitionFileName+".sql")
	if _, readErr := os.ReadDir(dirPath); readErr != nil && os.IsNotExist(readErr) {
		mkDirErr := os.MkdirAll(dirPath, 0644)
		if mkDirErr != nil {
			return nil, mkDirErr
		}
	}
	return []byte(psqlFileContents), os.WriteFile(definitionFilePath, []byte(psqlFileContents), 0644)
}

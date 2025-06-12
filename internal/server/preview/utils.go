package preview

import (
	"os"
	"path/filepath"
	"portal/internal/patcher"
	"portal/shared"
)

func patchPreview(filePath string, variables shared.FileVariables) error {
	globalFilePath := filepath.Join("component-preview/src/components", filePath)

	rawFile, err := os.ReadFile(globalFilePath)
	if err != nil {
		return err
	}

	newContent, err := patcher.PatchFile(string(rawFile), variables)
	if err != nil {
		return err
	}

	return os.WriteFile(globalFilePath, []byte(newContent), 0644)
}

package preview

import (
	"fmt"
	"os"
	"portal/internal/patcher"
	"portal/internal/server/github"
	"portal/shared"
)

func patchPreview(variables shared.PortalVariables) error {
	for filePath, fileVars := range variables {
		globalFilePath := fmt.Sprintf("%s/%s", github.RepoFolderName, filePath)

		rawFile, err := os.ReadFile(globalFilePath)
		if err != nil {
			return err
		}

		newContent, err := patcher.PatchFile(string(rawFile), fileVars)
		if err != nil {
			return err
		}

		err = os.WriteFile(globalFilePath, []byte(newContent), 0644)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

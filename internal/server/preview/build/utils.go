package build

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"portal/internal/server/github"
	"portal/shared"
	"strings"
)

// getComponentImports returns the imported dependencies imported (internal and external) from the file located at componentFilePath inside the cloned project.
func getComponentImports(componentFilePath string) (imports []string, err error) {
	fileExtensions := []string{"js", "ts", "jsx", "tsx", "mjs"}

	isValid := false

	for _, fileExt := range fileExtensions {
		if strings.HasSuffix(componentFilePath, "."+fileExt) {
			isValid = true
			break
		}
	}

	if !isValid {
		return []string{}, nil
	}

	file, err := os.ReadFile(filepath.Join(github.RepoFolderName, componentFilePath))
	if err != nil {
		return nil, err
	}

	out, err := shared.ExecuteTool("getComponentImports", map[string]any{
		"sourceCode": string(file),
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		Imports []string `json:"imports"`
	}

	if err = json.NewDecoder(&out).Decode(&result); err != nil {
		return nil, err
	}

	return result.Imports, nil
}

// installPackage runs npm install on the given dependency.
func installPackage(importPath string) error {
	cmd := exec.Command("npm", "install", importPath)

	cmd.Dir = "component-preview"
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

package build

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"portal/internal/server/github"
	"portal/shared"
)

func BuildComponentPage(componentFilePath string) error {
	err := os.MkdirAll("component-preview/src/components", 0755)
	if err != nil {
		return err
	}

	visitedImports := make(map[string]struct{})

	err = handleImports(componentFilePath, visitedImports)
	if err != nil {
		return err
	}

	slog.Info("Built component preview")

	return nil
}

func handleImports(componentFilePath string, visited map[string]struct{}) error {
	if _, ok := visited[componentFilePath]; ok {
		return nil
	}

	slog.Info("Importing " + componentFilePath)

	visited[componentFilePath] = struct{}{}

	file, err := os.ReadFile(filepath.Join(github.RepoFolderName, componentFilePath))
	if err != nil {
		return err
	}

	out, err := shared.ExecuteTool("previewComponent", map[string]any{
		"sourceCode": string(file),
	})
	if err != nil {
		return err
	}

	var result struct {
		SourceCode string   `json:"sourceCode"`
		Imports    []string `json:"imports"`
	}

	if err = json.NewDecoder(&out).Decode(&result); err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join("component-preview/src/components", filepath.Base(componentFilePath)), []byte(result.SourceCode), 0644)
	if err != nil {
		return err
	}

	for _, importPath := range result.Imports {
		if importPath[0] == '.' {
			// Other imports
			importedFilePath := filepath.Join(filepath.Dir(componentFilePath), importPath)
			fileExt, err := seekExtension(filepath.Join(github.RepoFolderName, importedFilePath), []string{"", "jsx", "tsx", "js", "ts"})
			if err != nil {
				return err
			}

			if err = handleImports(importedFilePath+"."+fileExt, visited); err != nil {
				return err
			}
		} else {
			continue
			// Packages
			cmd := exec.Command("npm", "install", importPath)

			cmd.Dir = "component-preview"
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

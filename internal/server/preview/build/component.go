package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"portal/internal/server/github"
)

func BuildComponentPage(componentFilePath string) error {
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", github.RepoFolderName, componentFilePath))
	if err != nil {
		return err
	}

	var input bytes.Buffer
	err = json.NewEncoder(&input).Encode(map[string]any{
		"sourceCode": string(file),
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("node", "tools/previewComponent.js")

	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Stdin = &input

	err = cmd.Run()
	if err != nil {
		return err
	}

	var result struct {
		SourceCode string   `json:"sourceCode"`
		Imports    []string `json:"imports"`
	}

	err = json.NewDecoder(&out).Decode(&result)
	if err != nil {
		return err
	}

	slog.Info("Initializing single component preview...")
	err = initComponentProject(result.SourceCode)
	if err != nil {
		return err
	}

	slog.Info("Handling imports...")
	err = handleImports(componentFilePath, result.Imports)
	if err != nil {
		return err
	}

	slog.Info("Built component preview")

	return nil
}

func initComponentProject(newSourceCode string) error {
	err := os.MkdirAll("component-preview/src/components", 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile("component-preview/src/components/ComponentPreview.tsx", []byte(newSourceCode), 0644)
	if err != nil {
		return err
	}

	return nil
}

func handleImports(componentFilePath string, imports []string) error {
	for _, importPath := range imports {
		if importPath[0] == '.' {
			// Relative imports
			srcPath := filepath.Join(github.RepoFolderName, filepath.Dir(componentFilePath), importPath)
			destPath := filepath.Join("component-preview/src/components", filepath.Base(importPath))

			fileExt, err := seekExtension(srcPath, []string{"", "jsx", "tsx", "js", "ts"})
			if err != nil {
				return err
			}

			err = copyFile(srcPath+fileExt, destPath+fileExt)
			if err != nil {
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

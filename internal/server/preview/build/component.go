package build

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"portal/internal/server/github"
	"portal/shared"
	"strings"
)

func BuildComponentPage(componentFilePath string) error {
	err := os.MkdirAll("component-preview/src/components", 0755)
	if err != nil {
		return err
	}

	visitedImports := make(map[string]struct{})

	if true {
		err = handleDependencies(componentFilePath, visitedImports)
		if err != nil {
			return err
		}
	}

	componentName, err := mockComponent(componentFilePath)
	if err != nil {
		return err
	}

	if err = makeEntryPoint(componentName, componentFilePath); err != nil {
		return err
	}

	slog.Info("Built component preview")

	return nil
}

func makeEntryPoint(componentName string, componentFilePath string) error {
	relPath, err := filepath.Rel("component-preview/src", filepath.Join("component-preview/src/components", componentFilePath))
	if err != nil {
		return err
	}

	fileContent := fmt.Sprintf(`import React from 'react';
import ReactDOM from 'react-dom/client';
import %s from './%s';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
	<React.StrictMode>
		<%s />
	</React.StrictMode>
);
`, componentName, relPath, componentName)

	return os.WriteFile("component-preview/src/index.js", []byte(fileContent), os.ModePerm)
}

func handleDependencies(componentFilePath string, visited map[string]struct{}) error {
	if _, ok := visited[componentFilePath]; ok {
		return nil
	}

	slog.Info("Importing " + componentFilePath)

	visited[componentFilePath] = struct{}{}

	file, err := os.ReadFile(filepath.Join(github.RepoFolderName, componentFilePath))
	if err != nil {
		return err
	}

	out, err := shared.ExecuteTool("getComponentImports", map[string]any{
		"sourceCode": string(file),
	})
	if err != nil {
		return err
	}

	var result struct {
		Imports []string `json:"imports"`
	}

	if err = json.NewDecoder(&out).Decode(&result); err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Join("component-preview/src/components", filepath.Dir(componentFilePath)), 0755); err != nil {
		return err
	}
	if err = copyFile(filepath.Join(github.RepoFolderName, componentFilePath), filepath.Join("component-preview/src/components", componentFilePath)); err != nil {
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

			if err = handleDependencies(importedFilePath+"."+fileExt, visited); err != nil {
				return err
			}
		} else {
			if importPath[0] != '@' {
				importPath = strings.Split(importPath, "/")[0]
			}

			if _, ok := visited[importPath]; ok {
				continue
			}

			visited[importPath] = struct{}{}

			slog.Info("Installing package " + importPath)

			// Packages
			cmd := exec.Command("npm", "install", importPath)

			cmd.Dir = "component-preview"
			//cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func mockComponent(componentFilePath string) (componentName string, err error) {
	file, err := os.ReadFile(filepath.Join(github.RepoFolderName, componentFilePath))
	if err != nil {
		return "", err
	}

	out, err := shared.ExecuteTool("mockComponentPreview", map[string]any{
		"sourceCode": string(file),
	})
	if err != nil {
		return "", err
	}

	var result struct {
		ComponentName string `json:"componentName"`
	}

	if err = json.NewDecoder(&out).Decode(&result); err != nil {
		return "", err
	}

	return result.ComponentName, nil
}

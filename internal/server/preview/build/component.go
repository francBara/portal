package build

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
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

	if err = handleDependencies(componentFilePath, visitedImports); err != nil {
		return err
	}

	component, err := scanComponent(componentFilePath)
	if err != nil {
		return err
	}

	if err = makeEntryPoint(component, componentFilePath); err != nil {
		return err
	}

	// Env files
	envPath := seekFiles([]string{".env.test", ".env.dev", ".env.prod"})
	if envPath != "" {
		envFile := filepath.Base(envPath)

		slog.Info(fmt.Sprintf("Copying %s into component-preview", envFile))
		copyFile(envPath, filepath.Join("component-preview", envFile))
	}

	// Configuration files
	for _, fileName := range []string{"tailwind.config.js", "postcss.config.mjs", "tailwind.config.mjs", "postcss.config.js"} {
		imports, err := getComponentImports(fileName)
		if err != nil {
			return err
		}

		for _, importPath := range imports {
			if importPath[0] != '@' {
				importPath = strings.Split(importPath, "/")[0]
			}
			slog.Info("Installing package " + importPath)
			installPackage(importPath)
		}

		copyFile(filepath.Join(github.RepoFolderName, fileName), filepath.Join("component-preview", fileName))
	}

	slog.Info("Built component preview")

	return nil
}

func makeEntryPoint(component componentMock, componentFilePath string) error {
	relPath, err := filepath.Rel("component-preview/src", filepath.Join("component-preview/src/components", componentFilePath))
	if err != nil {
		return err
	}

	variableDeclarations := ""
	componentProps := ""

	for name, value := range component.Mock {
		marshaledValue, err := json.Marshal(value)
		if err != nil {
			return err
		}

		variableDeclarations += fmt.Sprintf("const %s = %s;\n", name, string(marshaledValue))
		componentProps += fmt.Sprintf("%s={%s} ", name, name)
	}

	fileContent := fmt.Sprintf(`import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import "./index.css";
import %s from './%s';

%s
const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
	<React.StrictMode>
		<BrowserRouter>
			<div className="min-h-screen flex items-center justify-center">
				<div className="h-200 w-64">
					<%s %s/>
				</div>
			</div>
		</BrowserRouter>
	</React.StrictMode>
);
`, component.ComponentName, relPath, variableDeclarations, component.ComponentName, componentProps)

	if err = os.WriteFile("component-preview/src/index.jsx", []byte(fileContent), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile("component-preview/src/index.css", []byte("@tailwind base;\n@tailwind components;\n@tailwind utilities;\n"), os.ModePerm)

}

func handleDependencies(componentFilePath string, visited map[string]struct{}) error {
	if _, ok := visited[componentFilePath]; ok {
		return nil
	}

	slog.Info("Importing " + componentFilePath)

	visited[componentFilePath] = struct{}{}

	if err := os.MkdirAll(filepath.Join("component-preview/src/components", filepath.Dir(componentFilePath)), 0755); err != nil {
		return err
	}
	if err := copyFile(filepath.Join(github.RepoFolderName, componentFilePath), filepath.Join("component-preview/src/components", componentFilePath)); err != nil {
		return err
	}

	imports, err := getComponentImports(componentFilePath)
	if err != nil {
		return err
	}

	for _, importPath := range imports {
		if importPath[0] == '.' {
			// Other imports
			importedFilePath := filepath.Join(filepath.Dir(componentFilePath), importPath)
			fileExt, err := seekExtension(filepath.Join(github.RepoFolderName, importedFilePath), []string{"jsx", "tsx", "js", "ts"})
			if err != nil {
				return fmt.Errorf("seekExtension for %s: %w", importedFilePath, err)
			}

			if fileExt != "" {
				importedFilePath += "." + fileExt
			}

			if err = handleDependencies(importedFilePath, visited); err != nil {
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
			err := installPackage(importPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type componentMock struct {
	ComponentName string         `json:"componentName"`
	Mock          map[string]any `json:"mock"`
}

func scanComponent(componentFilePath string) (mock componentMock, err error) {
	file, err := os.ReadFile(filepath.Join(github.RepoFolderName, componentFilePath))
	if err != nil {
		return componentMock{}, err
	}

	out, err := shared.ExecuteTool("scanComponentPreview", map[string]any{
		"sourceCode": string(file),
	})
	if err != nil {
		return componentMock{}, err
	}

	var result componentMock

	if err = json.NewDecoder(&out).Decode(&result); err != nil {
		return componentMock{}, err
	}

	return result, nil
}

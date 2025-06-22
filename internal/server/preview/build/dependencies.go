package build

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"portal/internal/server/github"
)

// handleDependencies copies the component located at componentFilePath to its corresponding location in component-preview,
// then calls handleDependencies on its internal dependencies, and returns external dependencies.
func handleDependencies(componentFilePath string, visited map[string]struct{}) (externalDependencies map[string]string, err error) {
	// Avoids cycles
	if _, ok := visited[componentFilePath]; ok {
		return map[string]string{}, nil
	}

	slog.Info("Importing " + componentFilePath)

	visited[componentFilePath] = struct{}{}

	if err := os.MkdirAll(filepath.Join("component-preview/src/components", filepath.Dir(componentFilePath)), 0755); err != nil {
		return map[string]string{}, err
	}
	if err := copyFile(filepath.Join(github.RepoFolderName, componentFilePath), filepath.Join("component-preview/src/components", componentFilePath)); err != nil {
		return map[string]string{}, err
	}

	imports, err := getComponentImports(componentFilePath)
	if err != nil {
		return map[string]string{}, err
	}

	externalDependencies = make(map[string]string)

	for _, importPath := range imports {
		if importPath[0] == '.' {
			// Internal dependencies
			importedFilePath := filepath.Join(filepath.Dir(componentFilePath), importPath)
			fileExt, err := seekExtension(filepath.Join(github.RepoFolderName, importedFilePath), []string{"jsx", "tsx", "js", "ts"})
			if err != nil {
				return map[string]string{}, fmt.Errorf("seekExtension for %s: %w", importedFilePath, err)
			}

			if fileExt != "" {
				importedFilePath += "." + fileExt
			}

			currExtDependencies, err := handleDependencies(importedFilePath, visited)
			if err != nil {
				return map[string]string{}, err
			}

			for dep := range currExtDependencies {
				externalDependencies[dep] = ""
			}
		} else {
			// External dependencies
			if _, ok := visited[importPath]; ok {
				continue
			}

			visited[importPath] = struct{}{}

			externalDependencies[importPath] = ""
		}
	}

	return externalDependencies, nil
}

// applyVersions reads the project's package.json, and enforces the same versions to the necessary packages.
func applyVersions(dependencies map[string]string) error {
	var projectPackage struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	rawPackage, err := os.ReadFile(filepath.Join(github.RepoFolderName, "package.json"))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(rawPackage, &projectPackage); err != nil {
		return err
	}

	for dep := range dependencies {
		if _, ok := projectPackage.Dependencies[dep]; ok {
			dependencies[dep] = projectPackage.Dependencies[dep]
		}
		if _, ok := projectPackage.DevDependencies[dep]; ok {
			dependencies[dep] = projectPackage.DevDependencies[dep]
		}
	}

	return nil
}

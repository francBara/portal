package build

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"portal/internal/server/github"
	"strings"
)

// handleDependencies copies the component located at componentFilePath to its corresponding location in component-preview,
// then calls handleDependencies on its internal dependencies, and installs external dependencies.
func handleDependencies(componentFilePath string, visited map[string]struct{}) error {
	// Avoids cycles
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
			// Internal dependencies
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
			// External dependencies
			if importPath[0] != '@' {
				importPath = strings.Split(importPath, "/")[0]
			}

			if _, ok := visited[importPath]; ok {
				continue
			}

			visited[importPath] = struct{}{}

			slog.Info("Installing package " + importPath)

			if err = installPackage(importPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// handleDevDependencies install necessary dev dependencies, applying the project versions if present.
func handleDevDependencies() error {
	// Get project dev dependencies
	var projectPackage struct {
		DevDependencies map[string]string `json:"devDependencies"`
	}

	rawPackage, err := os.ReadFile(filepath.Join(github.RepoFolderName, "package.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(rawPackage, &projectPackage)
	if err != nil {
		return err
	}

	// Necessary
	devDependencies := map[string]string{
		"postcss":              "",
		"autoprefixer":         "",
		"tailwindcss":          "",
		"@vitejs/plugin-react": "",
		"vite":                 "",
	}

	for dependency, version := range projectPackage.DevDependencies {
		if _, ok := devDependencies[dependency]; ok {
			devDependencies[dependency] = version
		}
	}

	for dependency, version := range devDependencies {
		slog.Info("installing dev dependency " + dependency + ", version " + version)

		if version != "" {
			dependency = dependency + "@" + version
		}

		cmd := exec.Command("npm", "install", dependency, "--save-dev")

		cmd.Dir = "component-preview"
		cmd.Stderr = os.Stderr

		if err = cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

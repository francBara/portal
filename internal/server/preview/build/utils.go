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

// getComponentImports returns the imported dependencies imported (internal and external) from the file located at componentFilePath inside the cloned project.
func getComponentImports(componentFilePath string) (imports []string, err error) {
	fileExtensions := []string{"js", "ts", "jsx", "tsx", "mjs", "mts"}

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

	for i := range result.Imports {
		// External dependencies
		if len(result.Imports[i]) > 0 && result.Imports[i][0] != '.' && result.Imports[i][0] != '@' {
			result.Imports[i] = strings.Split(result.Imports[i], "/")[0]
		}
	}

	return result.Imports, nil
}

func installAll(dependencies map[string]string) error {
	verDependencies := []string{}

	for dep, version := range dependencies {
		if version != "" {
			verDependencies = append(verDependencies, fmt.Sprintf("%s@%s", dep, version))
		} else {
			verDependencies = append(verDependencies, dep)
		}
	}

	cmd := exec.Command("yarn", append([]string{"add"}, verDependencies...)...)

	for _, dep := range verDependencies {
		slog.Info("installing package " + dep)
	}

	cmd.Dir = "component-preview"

	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("yarn", "install")

	cmd.Dir = "component-preview"

	return cmd.Run()
}

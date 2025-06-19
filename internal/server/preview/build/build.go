package build

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"portal/internal/server/github"
	"strings"
	"sync"
)

var mutex sync.Mutex

var lastBuilt string

func BuildComponentPage(componentFilePath string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if componentFilePath == lastBuilt {
		return nil
	}

	err := os.MkdirAll("component-preview/src/components", 0755)
	if err != nil {
		return err
	}

	visitedImports := make(map[string]struct{})

	if err = makePackage(); err != nil {
		return err
	}

	component, err := scanComponent(componentFilePath)
	if err != nil {
		return err
	}

	if err = handleDependencies(componentFilePath, visitedImports); err != nil {
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
		if !fileExists(filepath.Join(github.RepoFolderName, fileName)) {
			continue
		}

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

	if err = handleDevDependencies(); err != nil {
		return err
	}

	slog.Info("Built component preview")

	lastBuilt = componentFilePath

	return nil
}

func makePackage() error {
	return os.WriteFile("component-preview/package.json", []byte(`{
    "name": "component-preview",
    "version": "1.0.0",
    "private": true,
    "dependencies": {
        "react": "^19.1.0",
        "react-dom": "^19.1.0",
        "react-router-dom": "^7.6.2",
        "react-scripts": "5.0.1"
    }
}
`), os.ModePerm)
}

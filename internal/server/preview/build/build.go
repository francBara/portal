package build

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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

	externalDependencies, err := handleDependencies(componentFilePath, visitedImports)
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
		_, dependencies, err := importConfigFile(fileName)
		if err != nil {
			return err
		}
		for _, dep := range dependencies {
			externalDependencies[dep] = ""
		}
	}

	viteConfigImported := false

	for _, fileName := range []string{"vite.config.mts", "vite.config.js"} {
		viteConfigImported, dependencies, err := importConfigFile(fileName)
		if err != nil {
			return err
		}
		for _, dep := range dependencies {
			externalDependencies[dep] = ""
		}
		if viteConfigImported {
			break
		}
	}

	if !viteConfigImported {
		if err = makeViteConfig(); err != nil {
			return err
		}
	}

	// Mandatory packages
	for _, dep := range []string{"autoprefixer", "postcss", "tailwindcss", "vite", "react", "react-dom", "react-router-dom", "react-scripts"} {
		externalDependencies[dep] = ""
	}

	applyVersions(externalDependencies)

	if err = installAll(externalDependencies); err != nil {
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
    "private": true
}
`), os.ModePerm)
}

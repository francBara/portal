package build

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"portal/shared"
)

func BuildComponentPage(componentFilePath string, componentVar shared.UIVariable, mocks shared.PortalMocks) error {
	err := os.MkdirAll("component-preview/src/components", 0755)
	if err != nil {
		return err
	}

	visitedImports := make(map[string]struct{})

	if err = makePackage(); err != nil {
		return err
	}

	externalDependencies, err := handleDependencies(componentFilePath, mocks, visitedImports)
	if err != nil {
		return err
	}

	if err = makeEntryPoint(componentVar, componentFilePath); err != nil {
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
	for _, fileName := range []string{"postcss.config.mjs", "postcss.config.js"} {
		_, dependencies, err := importConfigFile(fileName)
		if err != nil {
			return err
		}
		for _, dep := range dependencies {
			externalDependencies[dep] = ""
		}
	}

	dependencies, err := importTailwindConfig()
	if err != nil {
		return err
	}

	viteConfigImported := false

	for _, fileName := range []string{"vite.config.mts", "vite.config.js"} {
		var viteDependencies []string
		viteConfigImported, viteDependencies, err = importConfigFile(fileName)
		if err != nil {
			return err
		}

		dependencies = append(dependencies, viteDependencies...)

		if viteConfigImported {
			break
		}
	}

	if !viteConfigImported {
		if err = makeViteConfig(); err != nil {
			return err
		}
	}

	// Collect packages
	for _, dep := range append(dependencies, []string{"autoprefixer", "postcss", "tailwindcss", "vite", "react", "react-dom", "react-router-dom", "react-scripts"}...) {
		externalDependencies[dep] = ""
	}

	applyVersions(externalDependencies)

	if err = installAll(externalDependencies); err != nil {
		return err
	}

	slog.Info("Built component preview")

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

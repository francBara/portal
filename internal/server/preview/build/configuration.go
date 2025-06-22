package build

import (
	"os"
	"path/filepath"
	"portal/internal/server/github"
	"portal/shared"
)

func importConfigFile(filePath string) (imported bool, dependencies []string, err error) {
	if !fileExists(filepath.Join(github.RepoFolderName, filePath)) {
		return false, []string{}, nil
	}

	dependencies, err = getComponentImports(filePath)
	if err != nil {
		return false, []string{}, err
	}

	err = copyFile(filepath.Join(github.RepoFolderName, filePath), filepath.Join("component-preview", filePath))
	if err != nil {
		return false, []string{}, err
	}

	return true, dependencies, nil
}

func importTailwindConfig() (dependencies []string, err error) {
	var imported bool

	for _, fileName := range []string{"tailwind.config.js", "tailwind.config.mjs"} {
		imported, dependencies, err = importConfigFile(fileName)
		if err != nil {
			return []string{}, err
		}

		if imported {
			rawContent, err := os.ReadFile(filepath.Join("component-preview", fileName))
			if err != nil {
				return []string{}, err
			}

			out, err := shared.ExecuteTool("addTailwindPath", map[string]string{
				"sourceCode": string(rawContent),
				"newPath":    "../app-preview/src/**/*.{js,jsx,ts,tsx}",
			})
			if err != nil {
				return []string{}, err
			}

			if err = os.WriteFile(filepath.Join("component-preview", fileName), out.Bytes(), os.ModePerm); err != nil {
				return []string{}, err
			}

			return dependencies, nil
		}
	}

	return dependencies, nil
}

func makeViteConfig() error {
	return os.WriteFile("component-preview/vite.config.js", []byte(`import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  root: '.',
  build: {
    outDir: 'dist',
  },
  envPrefix: "REACT_APP_",
});`), os.ModePerm)
}

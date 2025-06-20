package build

import (
	"log/slog"
	"os"
	"path/filepath"
	"portal/internal/server/github"
	"strings"
)

func importConfigFile(filePath string) (imported bool, err error) {
	if !fileExists(filepath.Join(github.RepoFolderName, filePath)) {
		return false, nil
	}

	imports, err := getComponentImports(filePath)
	if err != nil {
		return false, err
	}

	for _, importPath := range imports {
		if importPath[0] != '@' {
			importPath = strings.Split(importPath, "/")[0]
		}
		slog.Info("Installing package " + importPath)
		installPackage(importPath)
	}

	err = copyFile(filepath.Join(github.RepoFolderName, filePath), filepath.Join("component-preview", filePath))
	if err != nil {
		return false, err
	}

	return true, nil
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

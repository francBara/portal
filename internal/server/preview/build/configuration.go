package build

import (
	"os"
	"path/filepath"
	"portal/internal/server/github"
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

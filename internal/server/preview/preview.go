package preview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"portal/internal/server/github"
	"portal/internal/server/preview/build"
)

// installPackages calls "npm install" on the cloned repo.
func installPackages() error {
	cmd := exec.Command("npm", "install")

	cmd.Dir = github.RepoFolderName

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// startDevServer calls vite on the cloned repo, serving a development server.
func startDevServer() {
	cmd := exec.Command("vite", "--port", "3001", "--mode", "test")

	cmd.Dir = github.RepoFolderName

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		panic("error running vite server: " + err.Error())
	}
}

func buildComponentPreview(componentFilePath string) error {
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", github.RepoFolderName, componentFilePath))
	if err != nil {
		return err
	}

	var input bytes.Buffer
	err = json.NewEncoder(&input).Encode(map[string]any{
		"sourceCode": string(file),
	})
	if err != nil {
		return err
	}

	cmd := exec.Command("node", "tools/previewComponent.js")

	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Stdin = &input

	err = cmd.Run()
	if err != nil {
		return err
	}

	var result struct {
		Imports []string `json:"imports"`
	}

	err = json.NewDecoder(&out).Decode(&result)
	if err != nil {
		return err
	}

	//initComponentProject()

	handleImports(componentFilePath, result.Imports)

	return nil
}

func initComponentProject() error {
	build.CopyDir("assets", "component-preview")

	return nil
}

func handleImports(componentFilePath string, imports []string) error {
	for _, importPath := range imports {
		if importPath[0] == '.' {
			// Relative imports
			srcPath := filepath.Join(github.RepoFolderName, componentFilePath, importPath)
			destPath := filepath.Join("component-preview", filepath.Base(importPath))

			build.CopyFile(srcPath, destPath)
		} else {
			// Packages
			cmd := exec.Command("npm", "install", importPath)

			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ServePreview sets up the repo for a local dev server, starts the dev server and proxies it.
func ServePreview() {
	slog.Info("executing npm install...")
	err := installPackages()
	if err != nil {
		slog.Error("npm install", "error", err.Error())
		return
	}

	slog.Info("starting vite dev server...")
	go startDevServer()

	target, _ := url.Parse("http://localhost:3001")
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	slog.Info("Preview proxy ready")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		slog.Error("preview proxy", "error", err.Error())
	}
}

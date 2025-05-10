package preview

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"portal/internal/server/github"
	"portal/internal/server/utils"
)

func installPackages() error {
	cmd := exec.Command("npm", "install")

	cmd.Dir = github.RepoFolderName

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

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

func ServePreview() {
	variables, err := utils.LoadVariables()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	for _, ui := range variables.UI {
		generatePreview(fmt.Sprintf("%s/%s", github.RepoFolderName, ui.FilePath))
		slog.Info("generated preview", "file", ui.FilePath)
	}

	err = installPackages()
	if err != nil {
		slog.Error("npm install", "error", err.Error())
		return
	}

	fmt.Println("Installed preview npm packages")

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

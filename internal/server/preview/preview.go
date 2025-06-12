package preview

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
)

var previewRunning bool

// startDevServer calls vite on the cloned repo, serving a development server.
func startDevServer() {
	cmd := exec.Command("npx", "vite", "--port", "3001", "--mode", "test")

	cmd.Dir = "component-preview"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		panic("error running vite server: " + err.Error())
	}
}

// ServePreview sets up the repo for a local dev server, starts the dev server and proxies it.
func ServePreview() {
	if previewRunning {
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

	previewRunning = true

	slog.Info("Preview proxy ready")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		slog.Error("preview proxy", "error", err.Error())
	}

	previewRunning = false
}

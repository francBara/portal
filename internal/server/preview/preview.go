package preview

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"time"
)

var proxyRunning bool
var proxyMutex sync.Mutex

var devServerMutex sync.Mutex
var cancelFunction context.CancelFunc

// startDevServer calls vite on the cloned repo, serving a development server.
func startDevServer() {
	devServerMutex.Lock()
	defer devServerMutex.Unlock()

	var ctx context.Context

	ctx, cancelFunction = context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "npx", "vite", "--port", "3001", "--mode", "test")

	cmd.Dir = "component-preview"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()

	cancelFunction = nil
}

func serveProxy() {
	proxyMutex.Lock()
	defer proxyMutex.Unlock()

	proxyRunning = true

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

	proxyRunning = false
}

// ServePreview sets up the repo for a local dev server, starts the dev server and proxies it.
func ServePreview() {
	if cancelFunction != nil {
		cancelFunction()
		time.Sleep(1 * time.Second)
	}

	go startDevServer()

	if !proxyRunning {
		go serveProxy()
	}
}

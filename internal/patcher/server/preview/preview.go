package preview

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
)

const previewFolderName = "app-preview"

func cloneRepo(repoUrl string, branchName string) {
	info, err := os.Stat(previewFolderName)

	if err == nil && info.IsDir() {
		fmt.Println("Skipping git clone")
		return
	}

	cmd := exec.Command("git", "clone", "--recurse-submodules", "--branch", branchName, "--single-branch", repoUrl, previewFolderName)

	err = cmd.Run()
	if err != nil {
		panic("error cloning repo: " + err.Error())
	}
}

func installPackages() {
	cmd := exec.Command("npm", "install")

	cmd.Dir = previewFolderName

	err := cmd.Run()
	if err != nil {
		panic("error cloning repo: " + err.Error())
	}
}

func startDevServer() {
	cmd := exec.Command("vite", "--port", "3001", "--mode", "test")

	cmd.Dir = previewFolderName

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		panic("error running vite server: " + err.Error())
	}
}

func ServePreview(repoUrl string, branchName string) *httputil.ReverseProxy {
	cloneRepo(repoUrl, branchName)

	fmt.Println("Cloned preview repo")

	installPackages()

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

	fmt.Println("Serving proxy")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic("preview proxy server error: " + err.Error())
	}

	return proxy
}

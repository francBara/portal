package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"portal/internal/patcher"
	"portal/internal/patcher/generator"
	"portal/internal/patcher/server/auth"
	"portal/shared"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func loadVariables(path string) shared.PortalVariables {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var data shared.PortalVariables
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		panic(err)
	}
	file.Close()

	return data
}

func isAsset(path string) bool {
	return strings.HasSuffix(path, ".js") ||
		strings.HasSuffix(path, ".css") ||
		strings.HasSuffix(path, ".map") ||
		strings.HasSuffix(path, ".json") ||
		strings.HasSuffix(path, ".ico") ||
		strings.HasSuffix(path, ".png") ||
		strings.HasSuffix(path, ".jpg") ||
		strings.HasSuffix(path, ".jpeg") ||
		strings.HasPrefix(path, "/assets/")
}

func RunPatcher(port int, variablesPath string) {
	variables := loadVariables(variablesPath)

	configs, err := patcher.LoadConfigs()
	if err != nil {
		log.Fatalln("Could not load config file")
	}

	var github GithubStub
	github.Init(configs.RepoName, configs.RepoOwner, configs.RepoBranch, configs.Pac)

	r := chi.NewRouter()

	staticDir := "./static"

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		fullPath := filepath.Join(staticDir, requestPath)

		// Check if file exists
		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() {
			http.ServeFile(w, r, fullPath)
			return
		}

		// Block fallback for known asset paths
		if isAsset(requestPath) {
			http.NotFound(w, r)
			return
		}

		// Otherwise, fallback to index.html for Vue SPA
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	// API routes
	// Handles basic authentication
	r.Post("/api/signin", auth.Signin(configs.Users))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthenticateUser(configs.Users))

		// Populates and serves the main dashboard
		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user").(*auth.PortalUser)

			dashboard := generator.GenerateDashboard(variables, user.Name)
			fmt.Fprint(w, dashboard)
		})

		// Applies the update to the remote repo
		r.Post("/patch", PatcherController(variables, github, configs))
	})

	log.Printf("Starting server on http://localhost:%d...", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(port)), r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

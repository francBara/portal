package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"portal/internal/patcher"
	"portal/internal/patcher/generator"
	"portal/internal/patcher/server/auth"
	"portal/internal/patcher/server/preview"
	"portal/shared"
	"strconv"

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

func RunPatcher(port int, variablesPath string) {
	variables := loadVariables(variablesPath)

	configs, err := patcher.LoadConfigs()
	if err != nil {
		log.Fatalln("Could not load config file")
	}

	go preview.ServePreview("https://github.com/togiftit/togiftit-web", "demo/portal", configs.Pac)

	var github GithubStub
	github.Init(configs.RepoName, configs.RepoOwner, configs.RepoBranch, configs.Pac)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	})

	// Serves login page
	r.Get("/signin", auth.GetSigninPage)
	// Handles basic authentication
	r.Post("/signin", auth.Signin(configs.Users))

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

		r.Post("/preview/update", preview.UpdatePreview(variables))
	})

	log.Printf("Starting server on http://localhost:%d...", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(port)), r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

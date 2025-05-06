package server

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"portal/internal/server/auth"
	"portal/internal/server/controllers"
	"portal/internal/server/github"
	"portal/internal/server/preview"
	"portal/internal/server/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func RunServer(port int) {
	configs, err := utils.LoadConfigs()
	if err != nil {
		log.Fatalln("Could not load config file")
	}

	err = github.Init(configs.RepoName, configs.RepoOwner, configs.RepoBranch, configs.Pac)
	if err != nil {
		slog.Error("Error initializing github client", err.Error())
	}

	if github.GithubClient != nil && configs.ServePreview {
		go preview.ServePreview()
	}

	r := chi.NewRouter()

	// Handles basic authentication
	r.Post("/auth/signin", auth.Signin())

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthenticateUser())

		// Gets the current variables
		r.Get("/variables", controllers.GetVariables())

		// Applies the update to the remote repo
		r.Post("/patch", controllers.PushChanges(configs))

		// Updates the preview with new variables
		r.Post("/preview/update", preview.UpdatePreview())
	})

	log.Printf("Starting server on http://localhost:%d...", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(port)), r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

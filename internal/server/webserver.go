package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"portal/internal/server/auth"
	"portal/internal/server/controllers"
	"portal/internal/server/github"
	"portal/internal/server/globals"
	"portal/internal/server/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func RunServer(port int) {
	configs := utils.LoadConfigs()

	configs.Print()

	err := github.Init(configs)
	if err != nil {
		log.Fatal("Error initializing github client", " ", err.Error())
	}

	globals.LoadVariables()

	r := chi.NewRouter()

	r.Route("/api", func(api chi.Router) {
		api.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))

		// Handles basic authentication
		api.Post("/auth/signin", auth.Signin())

		// Protected routes
		api.Group(func(secureApi chi.Router) {
			secureApi.Use(auth.AuthenticateUser())

			// Gets the current variables
			secureApi.Get("/variables", controllers.GetVariables())

			// Applies the update to the remote repo
			secureApi.Post("/patch", controllers.PushChanges(configs))
		})
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/", http.FileServer(http.Dir("static")))
		path := filepath.Join("static", r.URL.Path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, "static/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	log.Printf("Starting server on http://localhost:%d...", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(port)), r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

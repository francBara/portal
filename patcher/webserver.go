package patcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"portal/parser"
	"portal/patcher/auth"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func loadVariables(path string) parser.PortalVariables {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var data parser.PortalVariables
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		panic(err)
	}
	file.Close()

	return data
}

const newBranch = "nontech"

func RunPatcher(port int, variablesPath string) {
	variables := loadVariables(variablesPath)

	configs, err := LoadConfigs()
	if err != nil {
		log.Fatalln("Could not load config file")
	}

	var github GithubStub
	github.Init(configs)

	r := chi.NewRouter()

	r.Get("/signin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./patcher/static/login.html")
	})

	r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		var user auth.LoginUser

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		token, err := user.Login(configs.Users)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   3600,
		})
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthenticateUser(configs.Users))

		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user").(*auth.PortalUser)

			// Serves the main dashboard
			dashboard := GenerateDashboard(variables, user.Name)
			fmt.Fprint(w, dashboard)
		})

		r.Post("/patch", func(w http.ResponseWriter, r *http.Request) {
			var update map[string]string

			err := json.NewDecoder(r.Body).Decode(&update)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			newVariables, err := variables.UpdateVariables(update)
			if err != nil {
				http.Error(w, "Could not update variables", http.StatusBadRequest)
				return
			}

			if configs.OpenPullRequest {
				github.CreateBranch(newBranch)
			}

			var updateBranch string
			if configs.OpenPullRequest {
				updateBranch = newBranch
			} else {
				updateBranch = github.RepoBranch
			}

			user := r.Context().Value("user").(*auth.PortalUser)

			for filePath, _ := range newVariables.FileHashes {
				fileContent, fileSha := github.GetRepoFile(filePath)

				newContent := PatchFile(fileContent, newVariables)

				github.UpdateFile(newContent, filePath, fileSha, updateBranch, "Eccoci qua", *user)
			}

			github.CreatePullRequest(newBranch, "Nuova pull request", "Una bellissima pull request per pull requestare")
		})

	})

	log.Printf("Starting server on http://localhost:%d...", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(port)), r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

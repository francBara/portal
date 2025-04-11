package patcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"portal/parser"
	"strconv"
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

type PatcherConfigs struct {
	RepoOwner  string
	RepoName   string
	RepoBranch string
	Pac        string
}

func loadConfigs() PatcherConfigs {
	file, err := os.Open("./patcher_config.json")
	if err != nil {
		panic(err)
	}

	var data PatcherConfigs
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
	configs := loadConfigs()

	user := PortalUser{
		Name:  "Marcolino",
		Email: "marcolino@gmail.com",
	}

	var github GithubStub
	github.Init(configs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dashboard := GenerateDashboard(variables)
		fmt.Fprint(w, dashboard)
	})

	http.HandleFunc("/patch", func(w http.ResponseWriter, r *http.Request) {
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

		github.CreateBranch(newBranch)

		for filePath, _ := range newVariables.FileHashes {
			fileContent, fileSha := github.GetRepoFile(filePath)

			newContent := PatchFile(fileContent, newVariables)

			github.UpdateFile(newContent, filePath, fileSha, newBranch, "Eccoci qua", user)
		}

		github.CreatePullRequest(newBranch, "Nuova pull request", "Una bellissima pull request per pull requestare")
	})

	log.Printf("Starting server on http://localhost:%d...", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(port)), nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

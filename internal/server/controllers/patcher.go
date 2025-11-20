package controllers

import (
	"encoding/json"
	"net/http"
	"portal/internal/patcher"
	"portal/internal/server/auth"
	"portal/internal/server/github"
	"portal/internal/server/globals"
	"portal/internal/server/utils"
	"portal/shared"
)

type PatcherPayload struct {
	Update        shared.AllVariables `json:"update"`
	CommitMessage string              `json:"commitMessage"`
}

func PushChanges(configs utils.PatcherConfigs) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload PatcherPayload

		github := github.GithubClient

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		user := r.Context().Value("user").(*auth.PortalUser)
		variables, err := globals.LoadVariables()
		if err != nil {
			http.Error(w, "Error loading server variables", http.StatusInternalServerError)
			return
		}

		//TODO: verify that payload.Update is coherent with variables
		for repoUrl, repoVars := range payload.Update {
			repo := github.GetRepo(repoUrl)
			var updateBranch string

			if configs.OpenPullRequest {
				updateBranch, err = repo.CreateBranch(github.Client)
				if err != nil {
					http.Error(w, "Error creating new branch", http.StatusBadRequest)
					return
				}
			} else {
				updateBranch = repo.Branch
			}

			for filePath, fileVars := range repoVars {
				fileContent, fileSha := repo.GetRepoFile(github.Client, filePath)

				if variables[repoUrl][filePath].Hash != fileSha {
					http.Error(w, "File has been modified since last pull", http.StatusBadRequest)
					return
				}

				language, err := shared.GetLanguageRegex(filePath)
				if err != nil {
					http.Error(w, "File has unsupported language", http.StatusBadRequest)
					return
				}

				newContent, err := patcher.PatchFile(fileContent, fileVars, language)
				if err != nil {
					http.Error(w, "Could not patch file", http.StatusInternalServerError)
				}

				repo.UpdateFile(github.Client, newContent, filePath, fileSha, updateBranch, payload.CommitMessage, *user)
			}

			if configs.OpenPullRequest {
				repo.CreatePullRequest(github.Client, updateBranch, "Portal", payload.CommitMessage)
			}
		}
	}
}

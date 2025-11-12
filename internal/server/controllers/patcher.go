package controllers

import (
	"encoding/json"
	"net/http"
	"portal/internal/patcher"
	"portal/internal/server/auth"
	"portal/internal/server/github"
	"portal/internal/server/utils"
	"portal/shared"
)

type PatcherPayload struct {
	Update        shared.AllVariables `json:"update"`
	BranchName    string              `json:"branchName"`
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

		var updateBranch string

		user := r.Context().Value("user").(*auth.PortalUser)

		for repoUrl, repoVars := range payload.Update {
			for filePath, fileVars := range repoVars {
				repo := github.GetRepo(repoUrl)

				if configs.OpenPullRequest {
					updateBranch = payload.BranchName
					err = repo.CreateBranch(github.Client)
					if err != nil {
						//TODO: If the error is "branch already exists", ignore error
						http.Error(w, "Branch already exists", http.StatusBadRequest)
						return
					}
				} else {
					updateBranch = repo.Branch
				}

				fileContent, fileSha := repo.GetRepoFile(github.Client, filePath)

				newContent, err := patcher.PatchFile(fileContent, fileVars)
				if err != nil {
					http.Error(w, "Could not patch file", http.StatusInternalServerError)
				}

				repo.UpdateFile(github.Client, newContent, filePath, fileSha, updateBranch, payload.CommitMessage, *user)

				if configs.OpenPullRequest {
					repo.CreatePullRequest(github.Client, payload.BranchName, "Portal", payload.CommitMessage)
				}
			}
		}
	}
}

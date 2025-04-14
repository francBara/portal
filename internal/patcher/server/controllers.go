package server

import (
	"encoding/json"
	"net/http"
	"portal/internal/patcher"
	"portal/internal/patcher/server/auth"
	"portal/shared"
)

type patcherPayload struct {
	Update        map[string]string `json:"update"`
	BranchName    string            `json:"branchName"`
	CommitMessage string            `json:"commitMessage"`
}

func PatcherController(variables shared.PortalVariables, github GithubStub, configs patcher.PatcherConfigs) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload patcherPayload

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		newVariables, err := variables.UpdateVariables(payload.Update)
		if err != nil {
			http.Error(w, "Could not update variables", http.StatusBadRequest)
			return
		}

		var updateBranch string

		if configs.OpenPullRequest {
			updateBranch = payload.BranchName
			err = github.CreateBranch(payload.BranchName)
			if err != nil {
				http.Error(w, "Branch already exists", http.StatusBadRequest)
				return
			}
		} else {
			updateBranch = github.RepoBranch
		}

		user := r.Context().Value("user").(*auth.PortalUser)

		for filePath := range newVariables.FileHashes {
			fileContent, fileSha := github.GetRepoFile(filePath)

			newContent := patcher.PatchFile(fileContent, newVariables)

			github.UpdateFile(newContent, filePath, fileSha, updateBranch, "Eccoci qua", *user)
		}

		if configs.OpenPullRequest {
			github.CreatePullRequest(payload.BranchName, "Portal", payload.CommitMessage)
		}
	}
}

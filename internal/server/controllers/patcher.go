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
	Update        shared.VariablesMap `json:"update"`
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

		variables, err := utils.LoadVariables()
		if err != nil {
			http.Error(w, "Could not load variables", http.StatusInternalServerError)
			return
		}

		newVariables, err := variables.GetPatch(payload.Update)
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

		for filePath, fileVars := range newVariables {
			fileContent, fileSha := github.GetRepoFile(filePath)

			newContent, err := patcher.PatchFile(fileContent, fileVars)
			if err != nil {
				http.Error(w, "Could not patch file", http.StatusInternalServerError)
			}

			github.UpdateFile(newContent, filePath, fileSha, updateBranch, "Eccoci qua", *user)
		}

		if configs.OpenPullRequest {
			github.CreatePullRequest(payload.BranchName, "Portal", payload.CommitMessage)
		}
	}
}

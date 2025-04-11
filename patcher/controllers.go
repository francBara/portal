package patcher

import (
	"encoding/json"
	"net/http"
	"portal/parser"
	"portal/patcher/auth"
)

type patcherPayload struct {
	Update        map[string]string `json:"update"`
	BranchName    string            `json:"branchName"`
	CommitMessage string            `json:"commitMessage"`
}

func PatcherController(variables parser.PortalVariables, github GithubStub, configs PatcherConfigs) func(w http.ResponseWriter, r *http.Request) {
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

		for filePath, _ := range newVariables.FileHashes {
			fileContent, fileSha := github.GetRepoFile(filePath)

			newContent := PatchFile(fileContent, newVariables)

			github.UpdateFile(newContent, filePath, fileSha, updateBranch, "Eccoci qua", *user)
		}

		if configs.OpenPullRequest {
			github.CreatePullRequest(payload.BranchName, "Portal", payload.CommitMessage)
		}
	}
}

package preview

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"portal/internal/patcher"
	"portal/internal/server/github"
	"portal/internal/server/utils"
	"portal/shared"
)

func UpdatePreview() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		varsUpdate, err := shared.JsonToVariablesMap(r.Body)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		variables, err := utils.LoadVariables()
		if err != nil {
			http.Error(w, "Could not load variables", http.StatusInternalServerError)
			return
		}

		newVariables, err := variables.GetPatch(varsUpdate)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Could not update variables", http.StatusInternalServerError)
			return
		}

		for filePath, fileVars := range newVariables {
			globalFilePath := fmt.Sprintf("%s/%s", github.RepoFolderName, filePath)

			rawFile, err := os.ReadFile(globalFilePath)
			if err != nil {
				slog.Error("Error reading file:", "error", err)
				return
			}

			newContent, err := patcher.PatchFile(string(rawFile), fileVars)
			if err != nil {
				http.Error(w, "Could not patch file", http.StatusInternalServerError)
				return
			}

			err = os.WriteFile(globalFilePath, []byte(newContent), 0644)
			if err != nil {
				panic(err)
			}
		}
	}
}

type previewStatus struct {
	IsPreviewAvailable bool `json:"isPreviewAvailable"`
}

func GetPreviewStatus(isPreviewAvailable bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(previewStatus{
			IsPreviewAvailable: isPreviewAvailable,
		})
	}
}

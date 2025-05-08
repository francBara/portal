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
)

func UpdatePreview() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var varsUpdate map[string]string

		err := json.NewDecoder(r.Body).Decode(&varsUpdate)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		variables, err := utils.LoadVariables()
		if err != nil {
			http.Error(w, "Could not load variables", http.StatusInternalServerError)
			return
		}

		newVariables, err := variables.UpdateVariables(varsUpdate)
		if err != nil {
			http.Error(w, "Could not update variables", http.StatusInternalServerError)
			return
		}

		for filePath := range newVariables.FileHashes {
			globalFilePath := fmt.Sprintf("%s/%s", github.RepoFolderName, filePath)

			rawFile, err := os.ReadFile(globalFilePath)
			if err != nil {
				slog.Error("Error reading file:", "error", err)
				return
			}

			newContent := patcher.PatchFile(string(rawFile), newVariables)

			err = os.WriteFile(globalFilePath, []byte(newContent), 0644)
			if err != nil {
				panic(err)
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

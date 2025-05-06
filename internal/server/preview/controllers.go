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

		newVariables, err := utils.LoadVariables().UpdateVariables(varsUpdate)
		if err != nil {
			panic("error in updating preview variables: " + err.Error())
		}

		for filePath := range newVariables.FileHashes {
			globalFilePath := fmt.Sprintf("%s/%s", github.RepoFolderName, filePath)

			rawFile, err := os.ReadFile(globalFilePath)
			if err != nil {
				slog.Error("Error reading file:", err)
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

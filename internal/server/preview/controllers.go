package preview

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"portal/internal/server/preview/build"
	"portal/internal/server/utils"
	"portal/shared"
)

var currentComponentPath string

func UpdatePreview() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if currentComponentPath == "" {
			http.Error(w, "No component was selected", http.StatusBadRequest)
			return
		}

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

		err = patchPreview(currentComponentPath, newVariables[currentComponentPath])
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Could not update preview", http.StatusInternalServerError)
			return
		}
	}
}

type buildComponentPayload struct {
	FilePath string `json:"filePath"`
}

func BuildComponentPreview() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload buildComponentPayload

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		currentComponentPath = payload.FilePath

		err = build.BuildComponentPage(payload.FilePath)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Could not build component preview", http.StatusInternalServerError)
			return
		}

		go ServePreview()
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

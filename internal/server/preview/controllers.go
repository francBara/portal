package preview

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"portal/internal/server/preview/build"
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

		err = patchPreview(newVariables)
		if err != nil {
			http.Error(w, "Could not update preview", http.StatusInternalServerError)
			return
		}
	}
}

type highlightComponentPayload struct {
	FilePath   string `json:"filePath"`
	UIVariable string `json:"varName"`
	NodeId     int    `json:"nodeId"`
}

func HighlightComponent() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload highlightComponentPayload

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		variables, err := utils.LoadVariables()
		if err != nil {
			http.Error(w, "Could not load variables", http.StatusInternalServerError)
			return
		}

		uiVar := variables[payload.FilePath].UI[payload.UIVariable]

		uiVar.HighlightedNode = payload.NodeId

		if fileVars, ok := variables[payload.FilePath]; ok {
			fileVars.UI[payload.UIVariable] = uiVar
		} else {
			http.Error(w, "File not found", http.StatusBadRequest)
			return
		}

		err = patchPreview(variables)
		if err != nil {
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

		err = build.BuildComponentPage(payload.FilePath)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Could not build component preview", http.StatusInternalServerError)
			return
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

package preview

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"portal/internal/server/preview/build"
	"portal/internal/server/utils"
	"portal/shared"
	"sync"
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

var mutex sync.Mutex

func BuildComponentPreview() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		var payload buildComponentPayload

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		if payload.FilePath == "" {
			http.Error(w, "empty component file path", http.StatusBadRequest)
			return
		}

		if payload.FilePath == currentComponentPath {
			w.WriteHeader(http.StatusOK)
			return
		}

		err = build.BuildComponentPage(payload.FilePath)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Could not build component preview", http.StatusInternalServerError)
			return
		}

		ServePreview()

		currentComponentPath = payload.FilePath
	}
}

type highlightNodePayload struct {
	FilePath string `json:"filePath"`
	NodeId   int    `json:"nodeId"`
	VarName  string `json:"varName"`
}

func HighlightNode() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		var payload highlightNodePayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		rawContent, err := os.ReadFile(filepath.Join("component-preview/src/components", payload.FilePath))
		if err != nil {
			http.Error(w, "File not found", http.StatusBadRequest)
			return
		}

		out, err := shared.ExecuteTool("highlightNode", map[string]any{
			"sourceCode": string(rawContent),
			"nodeId":     payload.NodeId,
			"rootName":   payload.VarName,
		})
		if err != nil {
			http.Error(w, "Could not highlight component", http.StatusInternalServerError)
			return
		}

		if err = os.WriteFile(filepath.Join("component-preview/src/components", payload.FilePath), out.Bytes(), os.ModePerm); err != nil {
			http.Error(w, "Could not update component", http.StatusInternalServerError)
		}
	}
}

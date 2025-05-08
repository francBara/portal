package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"portal/internal/server/utils"
)

func GetVariables() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		variables, err := utils.LoadVariables()
		if err != nil {
			slog.Error("GET api/variables", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(variables.ToMap())
	}
}

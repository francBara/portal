package controllers

import (
	"encoding/json"
	"net/http"
	"portal/internal/server/utils"
)

func GetVariables() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.LoadVariables().ToMap())
	}
}

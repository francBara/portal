package controllers

import (
	"encoding/json"
	"net/http"
	"portal/shared"
)

func GetVariables(variables shared.PortalVariables) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(variables)
	}
}

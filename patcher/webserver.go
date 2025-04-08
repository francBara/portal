package patcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"portal/parser"
	"strconv"
)

func loadVariables(path string) parser.PortalVariables {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var data parser.PortalVariables
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		panic(err)
	}
	file.Close()

	return data
}

func RunPatcher(port int, variablesPath string) {
	variables := loadVariables(variablesPath)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dashboard := GenerateDashboard(variables)
		fmt.Fprint(w, dashboard)
	})

	http.HandleFunc("/patch", func(w http.ResponseWriter, r *http.Request) {
		//client := GetGithubClient()
	})

	log.Printf("Starting server on http://localhost:%d ...", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(port)), nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}

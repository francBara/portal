package utils

import (
	"encoding/json"
	"os"
	"portal/internal/parser"
	"portal/internal/server/github"
	"portal/shared"
)

var variables *shared.PortalVariables

func LoadVariables() shared.PortalVariables {
	if variables != nil {
		return *variables
	}

	if github.GithubClient != nil {
		vars, err := parser.ParseProject(github.RepoFolderName, parser.ParseOptions{})
		if err != nil {
			panic(err)
		}
		variables = &vars
	} else {
		file, err := os.Open("variables.json")
		if err != nil {
			panic(err)
		}

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(variables); err != nil {
			panic(err)
		}
		file.Close()
	}

	return *variables
}

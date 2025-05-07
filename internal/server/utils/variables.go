package utils

import (
	"encoding/json"
	"log/slog"
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
		slog.Info("Parsing project")
		vars, err := parser.ParseProject(github.RepoFolderName, parser.ParseOptions{})
		if err != nil {
			panic(err)
		}
		slog.Info("Parsed project", "variables", vars.Length())

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

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

func LoadVariables() (shared.PortalVariables, error) {
	if variables != nil {
		return *variables, nil
	}

	if github.GithubClient != nil {
		slog.Info("Parsing project")
		vars, err := parser.ParseProject(github.RepoFolderName, parser.ParseOptions{})
		if err != nil {
			return shared.PortalVariables{}, err
		}
		slog.Info("Parsed project", "variables", vars.Length())

		variables = &vars
	} else {
		file, err := os.Open("variables.json")
		if err != nil {
			return shared.PortalVariables{}, err
		}

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(variables); err != nil {
			return shared.PortalVariables{}, err
		}
		file.Close()
	}

	return *variables, nil
}

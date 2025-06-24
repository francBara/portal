package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"portal/internal/parser"
	"portal/internal/server/github"
	"portal/shared"
)

var variables *shared.PortalVariables
var mocks *shared.PortalMocks

// LoadVariables loads and caches variables. If Github is set, variables are parsed in real time, otherwise are loaded from variables.json.
func LoadVariables() (shared.PortalVariables, error) {
	if variables != nil {
		return *variables, nil
	}

	if github.GithubClient != nil {
		slog.Info("Parsing project")
		vars, varMocks, err := parser.ParseProject(github.RepoFolderName, parser.ParseOptions{})
		if err != nil {
			return shared.PortalVariables{}, err
		}
		slog.Info("Parsed project", "variables", vars.Length())

		variables = &vars
		mocks = &varMocks
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

func LoadMocks() (shared.PortalMocks, error) {
	if mocks != nil {
		return *mocks, nil
	}

	if github.GithubClient != nil {
		slog.Info("Parsing project")
		vars, varMocks, err := parser.ParseProject(github.RepoFolderName, parser.ParseOptions{})
		if err != nil {
			return shared.PortalMocks{}, err
		}
		slog.Info("Parsed project", "variables", vars.Length())

		variables = &vars
		mocks = &varMocks
	} else {
		return shared.PortalMocks{}, fmt.Errorf("github client not set")
	}

	return *mocks, nil
}

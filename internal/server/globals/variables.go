package globals

import (
	"log/slog"
	"portal/internal/parser"
	"portal/internal/server/github"
	"portal/shared"
)

var variables *shared.AllVariables

// LoadVariables loads and caches variables. If Github is set, variables are parsed in real time, otherwise are loaded from variables.json.
func LoadVariables() (shared.AllVariables, error) {
	if variables != nil {
		return *variables, nil
	}

	slog.Info("Parsing repos")

	repoUrls := []string{}
	for _, repo := range github.GithubClient.Repos {
		repoUrls = append(repoUrls, repo.GetUrl())
	}

	vars, err := parser.ParseAll(repoUrls, parser.ParseOptions{})
	if err != nil {
		return shared.AllVariables{}, err
	}
	slog.Info("Parsed repos", "repos", len(github.GithubClient.Repos))

	variables = &vars

	return *variables, nil
}

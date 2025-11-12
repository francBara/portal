package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Repo struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
}

type PatcherConfigs struct {
	ReposNames      []string `json:"reposNames"`
	ReposOwners     []string `json:"reposOwners"`
	ReposBranches   []string `json:"reposBranches"`
	GithubUsername  string   `json:"githubUsername"`
	Pac             string   `json:"pac"`
	OpenPullRequest bool     `json:"openPullRequest"`
	ServePreview    bool     `json:"servePreview"`
}

func (config PatcherConfigs) Print() {
	if config.Pac != "" {
		config.Pac = "[REDACTED]"
	}

	jsonConfigs, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		slog.Error("Could not print config")
		return
	}

	fmt.Println(string(jsonConfigs))
}

var configFileCandidates = []struct {
	Name string
	Type string
}{
	{"config", "json"},
	{"config", "yaml"},
	{"config", "toml"},
}

func LoadConfigs() PatcherConfigs {
	godotenv.Load()

	viper.AutomaticEnv()
	viper.AddConfigPath(".")

	viper.SetDefault("repoBranch", "main")
	viper.SetDefault("openPullRequest", true)
	viper.SetDefault("servePreview", true)

	var config PatcherConfigs

	viper.BindEnv("githubUsername", "GITHUB_USERNAME")
	viper.BindEnv("reposNames", "REPOS_NAMES")
	viper.BindEnv("reposOwners", "REPOS_OWNERS")
	viper.BindEnv("reposBranches", "REPOS_BRANCHES")
	viper.BindEnv("pac", "PAC")
	viper.BindEnv("openPullRequest", "OPEN_PULL_REQUEST")
	viper.BindEnv("servePreview", "SERVE_PREVIEW")

	for _, candidate := range configFileCandidates {
		viper.SetConfigName(candidate.Name)
		viper.SetConfigType(candidate.Type)

		err := viper.ReadInConfig()
		if err == nil {
			fmt.Println("Loaded config:", viper.ConfigFileUsed())
			break
		}
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	if config.Pac == "" {
		log.Fatal("PAC not provided")
	}

	return config
}

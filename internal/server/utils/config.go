package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type PatcherConfigs struct {
	RepoOwner       string `json:"repoOwner"`
	GithubUsername  string `json:"githubUsername"`
	RepoName        string `json:"repoName"`
	RepoBranch      string `json:"repoBranch"`
	Pac             string `json:"pac"`
	OpenPullRequest bool   `json:"openPullRequest"`
	ServePreview    bool   `json:"servePreview"`
}

func (config PatcherConfigs) Print() {
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

	viper.BindEnv("repoOwner", "REPO_OWNER")
	viper.BindEnv("githubUsername", "GITHUB_USERNAME")
	viper.BindEnv("repoName", "REPO_NAME")
	viper.BindEnv("repoBranch", "REPO_BRANCH")
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

	if config.RepoOwner == "" {
		config.RepoOwner = config.GithubUsername
	}

	return config
}

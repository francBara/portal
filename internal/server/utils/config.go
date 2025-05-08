package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/spf13/viper"
)

type PatcherConfigs struct {
	RepoOwner       string `json:"repoOwner"`
	UserName        string `json:"userName"`
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

func LoadConfigs() (PatcherConfigs, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("portal")

	viper.AddConfigPath(".")

	viper.SetDefault("repoBranch", "main")
	viper.SetDefault("openPullRequest", true)
	viper.SetDefault("servePreview", true)

	var config PatcherConfigs

	for _, candidate := range configFileCandidates {
		viper.SetConfigName(candidate.Name)
		viper.SetConfigType(candidate.Type)

		err := viper.ReadInConfig()
		if err == nil {
			fmt.Println("Loaded config:", viper.ConfigFileUsed())

			err := viper.Unmarshal(&config)
			if err != nil {
				log.Fatalf("unable to decode into struct, %v", err)
			}

			if config.RepoOwner == "" {
				config.RepoOwner = config.UserName
			}

			return config, nil
		}
	}

	return config, errors.New("did not find config file")
}

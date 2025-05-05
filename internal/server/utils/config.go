package utils

import (
	"errors"
	"fmt"
	"log"
	"portal/internal/server/auth"

	"github.com/spf13/viper"
)

type PatcherConfigs struct {
	RepoOwner       string `json:"repoOwner"`
	RepoName        string `json:"repoName"`
	RepoBranch      string `json:"repoBranch"`
	Pac             string `json:"pac"`
	OpenPullRequest bool   `json:"openPullRequest"`
	ServePreview    bool   `json:"servePreview"`

	Users []auth.PortalUser `json:"users"`
}

var configFileCandidates = []struct {
	Name string
	Type string
}{
	{"patcher_config", "json"},
	{"patcher_config", "yaml"},
	{"config", "json"},
	{"config", "yaml"},
	{"app", "toml"},
}

func LoadConfigs() (PatcherConfigs, error) {
	viper.AutomaticEnv()
	viper.AddConfigPath(".")

	viper.SetDefault("repoBranch", "main")
	viper.SetDefault("openPullRequest", true)

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

			return config, nil
		}
	}

	return config, errors.New("did not find config file")
}

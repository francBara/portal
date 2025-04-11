package patcher

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type PatcherConfigs struct {
	RepoOwner  string `json:"repoOwner"`
	RepoName   string `json:"repoName"`
	RepoBranch string `json:"repoBranch"`
	Pac        string `json:"pac"`

	Users []PortalUser `json:"users"`
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

	viper.SetDefault("RepoBranch", "main")

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

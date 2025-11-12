package github

import (
	"context"
	"os"
	"portal/internal/server/utils"
	"strings"

	"github.com/google/go-github/v71/github"
	"golang.org/x/oauth2"
)

type GithubStub struct {
	Client   *github.Client
	UserName string
	Repos    []Repo
	Pac      string
}

var GithubClient *GithubStub

func Init(configs utils.PatcherConfigs) error {
	var stub GithubStub

	for i := 0; i < len(configs.ReposNames); i++ {
		stub.Repos = append(stub.Repos, Repo{
			Name:   strings.TrimSpace(configs.ReposNames[i]),
			Owner:  strings.TrimSpace(configs.ReposOwners[i]),
			Branch: strings.TrimSpace(configs.ReposBranches[i]),
		})
	}

	stub.UserName = configs.GithubUsername
	stub.Pac = configs.Pac

	os.Setenv("PORTAL_GH_TOKEN", configs.Pac)

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: configs.Pac},
	)
	tc := oauth2.NewClient(ctx, ts)

	stub.Client = github.NewClient(tc)

	for _, repo := range stub.Repos {
		err := repo.clone(stub)
		if err != nil {
			return err
		}
	}

	GithubClient = &stub

	return nil
}

// TODO: Replace repos array with map
func (stub GithubStub) GetRepo(repoUrl string) Repo {
	for _, repo := range stub.Repos {
		if repo.GetUrl() == repoUrl {
			return repo
		}
	}
	return stub.Repos[0]
}

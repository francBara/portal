package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v71/github"
	"golang.org/x/oauth2"
)

type GithubStub struct {
	Client     *github.Client
	RepoName   string
	RepoBranch string
	UserName   string
	RepoOwner  string
}

var GithubClient *GithubStub

const RepoFolderName = "app-preview"

func Init(repoName string, repoOwner string, userName string, repoBranch string, pac string) error {
	var githubClient GithubStub

	githubClient.RepoName = repoName
	githubClient.RepoOwner = repoOwner
	githubClient.UserName = userName
	githubClient.RepoBranch = repoBranch

	if pac != "" {
		os.Setenv("PORTAL_GH_TOKEN", pac)
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pac},
	)
	tc := oauth2.NewClient(ctx, ts)

	githubClient.Client = github.NewClient(tc)

	err := githubClient.cloneRepo(pac)
	if err != nil {
		return err
	}

	GithubClient = &githubClient

	return nil
}

func (stub GithubStub) GetRepoUrl() string {
	return fmt.Sprintf("https://github.com/%s/%s", stub.RepoOwner, stub.RepoName)
}

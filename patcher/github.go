package patcher

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v71/github"
	"golang.org/x/oauth2"
)

func GetGithubClient() *github.Client {
	ctx := context.Background()

	token := os.Getenv("PORTAL_GH_TOKEN")

	if token == "" {
		log.Fatalln("PORTAL_GH_TOKEN env variable is not set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

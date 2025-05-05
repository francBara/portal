package utils

import (
	"context"
	"log"
	"os"
	"portal/internal/server/auth"
	"time"

	"github.com/google/go-github/v71/github"
	"golang.org/x/oauth2"
)

type GithubStub struct {
	Client     *github.Client
	RepoName   string
	RepoBranch string
	RepoOwner  string
}

func (stub *GithubStub) Init(repoName string, repoOwner string, repoBranch string, pac string) {
	stub.RepoName = repoName
	stub.RepoOwner = repoOwner
	stub.RepoBranch = repoBranch

	if pac != "" {
		os.Setenv("PORTAL_GH_TOKEN", pac)
	}

	stub.getGithubClient()
}

func (stub *GithubStub) getGithubClient() {
	ctx := context.Background()

	token := os.Getenv("PORTAL_GH_TOKEN")

	if token == "" {
		log.Fatalln("PORTAL_GH_TOKEN env variable is not set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	stub.Client = github.NewClient(tc)
}

func (stub GithubStub) CreateBranch(branchName string) error {
	ctx := context.Background()

	baseRef, _, err := stub.Client.Git.GetRef(ctx, stub.RepoOwner, stub.RepoName, "refs/heads/"+stub.RepoBranch)
	if err != nil {
		log.Fatalf("Error getting base branch: %v", err)
	}

	newRef := &github.Reference{
		Ref: github.Ptr("refs/heads/" + branchName),
		Object: &github.GitObject{
			SHA: baseRef.Object.SHA,
		},
	}
	_, _, err = stub.Client.Git.CreateRef(ctx, stub.RepoOwner, stub.RepoName, newRef)
	if err != nil {
		return err
	}

	return nil
}

func (stub GithubStub) GetRepoFile(filePath string) (string, string) {
	ctx := context.Background()

	fileContent, _, _, err := stub.Client.Repositories.GetContents(ctx, stub.RepoOwner, stub.RepoName, filePath, &github.RepositoryContentGetOptions{
		Ref: stub.RepoBranch,
	})
	if err != nil {
		log.Fatalf("Error getting file contents: %v", err)
	}

	decodedContent, err := fileContent.GetContent()

	if err != nil {
		log.Fatalf("Error decoding file contents: %v", err)
	}

	return decodedContent, fileContent.GetSHA()
}

func (stub GithubStub) UpdateFile(newContent string, filePath string, oldFileSha string, branch string, commitMessage string, fromUser auth.PortalUser) {
	ctx := context.Background()

	options := &github.RepositoryContentFileOptions{
		Message: github.Ptr(commitMessage),
		Content: []byte(newContent),
		Branch:  github.Ptr(branch),
		SHA:     github.Ptr(oldFileSha),
		Committer: &github.CommitAuthor{
			Name:  github.Ptr(fromUser.Name),
			Email: github.Ptr(fromUser.Email),
			Date:  &github.Timestamp{Time: time.Now()},
		},
	}

	_, _, err := stub.Client.Repositories.UpdateFile(ctx, stub.RepoOwner, stub.RepoName, filePath, options)
	if err != nil {
		log.Fatalf("Error updating file: %v", err)
	}
}

func (stub GithubStub) CreatePullRequest(fromBranch string, title string, body string) {
	ctx := context.Background()

	pr := &github.NewPullRequest{
		Title: github.Ptr(title),
		Head:  github.Ptr(fromBranch),
		Base:  github.Ptr(stub.RepoBranch),
		Body:  github.Ptr(body),
	}

	_, _, err := stub.Client.PullRequests.Create(ctx, stub.RepoOwner, stub.RepoName, pr)
	if err != nil {
		log.Fatalf("Error creating PR: %v", err)
	}
}

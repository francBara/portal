package github

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"portal/internal/server/auth"
	"time"

	"github.com/google/go-github/v71/github"
)

type Repo struct {
	Name   string
	Owner  string
	Branch string
}

func (repo Repo) GetUrl() string {
	return fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
}

func (repo Repo) GetWebUrl() string {
	return fmt.Sprintf("https://github.com/%s/%s", repo.Owner, repo.Name)

}

func (repo Repo) Clone(stub GithubStub) error {
	slog.Info("Cloning", "repo", repo.GetUrl(), "user", stub.UserName)

	cred := fmt.Sprintf("https://%s:%s@github.com\n", stub.UserName, stub.Pac)
	err := os.WriteFile(os.Getenv("HOME")+"/.git-credentials", []byte(cred), 0600)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "config", "--global", "credential.helper", "store")
	err = cmd.Run()
	if err != nil {
		return err
	}

	os.RemoveAll(path.Join("repos", repo.GetUrl()))

	cmd = exec.Command("git", "clone", "--recurse-submodules", "--branch", repo.Branch, "--single-branch", repo.GetWebUrl(), fmt.Sprintf("%s/%s", "repos", repo.GetUrl()))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (repo Repo) CreateBranch(client *github.Client) (newBranchName string, err error) {
	ctx := context.Background()

	baseRef, _, err := client.Git.GetRef(ctx, repo.Owner, repo.Name, "refs/heads/"+repo.Branch)
	if err != nil {
		log.Fatalf("Error getting base branch: %v", err)
	}

	newBranchName = getNewBranchName()

	newRef := &github.Reference{
		Ref: github.Ptr(fmt.Sprintf("refs/heads/%s", newBranchName)),
		Object: &github.GitObject{
			SHA: baseRef.Object.SHA,
		},
	}

	_, _, err = client.Git.CreateRef(ctx, repo.Owner, repo.Name, newRef)
	if err != nil {
		return "", err
	}

	return newBranchName, nil
}

func (repo Repo) GetRepoFile(client *github.Client, filePath string) (string, string) {
	ctx := context.Background()

	fileContent, _, _, err := client.Repositories.GetContents(ctx, repo.Owner, repo.Name, filePath, &github.RepositoryContentGetOptions{
		Ref: repo.Branch,
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

func (repo Repo) UpdateFile(client *github.Client, newContent string, filePath string, oldFileSha string, branch string, commitMessage string, fromUser auth.PortalUser) {
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

	_, _, err := client.Repositories.UpdateFile(ctx, repo.Owner, repo.Name, filePath, options)
	if err != nil {
		log.Fatalf("Error updating file: %v", err)
	}
}

func (repo Repo) CreatePullRequest(client *github.Client, fromBranch string, title string, body string) {
	ctx := context.Background()

	pr := &github.NewPullRequest{
		Title: github.Ptr(title),
		Head:  github.Ptr(fromBranch),
		Base:  github.Ptr(repo.Branch),
		Body:  github.Ptr(body),
	}

	_, _, err := client.PullRequests.Create(ctx, repo.Owner, repo.Name, pr)
	if err != nil {
		log.Fatalf("Error creating PR: %v", err)
	}
}

func getNewBranchName() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	randomChars := make([]byte, 6)
	for i := range randomChars {
		randomChars[i] = charset[rand.Intn(len(charset))]
	}

	return fmt.Sprintf("portal/%s", randomChars)
}

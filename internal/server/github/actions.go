package github

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"portal/internal/server/auth"
	"time"

	"github.com/google/go-github/v71/github"
)

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

func (stub GithubStub) cloneRepo(pac string) error {
	info, err := os.Stat(RepoFolderName)

	if err == nil && info.IsDir() {
		slog.Info("Skipping git clone")
		return nil
	}

	slog.Info("Cloning", stub.RepoName, stub.RepoBranch)

	cred := fmt.Sprintf("https://%s:%s@github.com\n", stub.UserName, pac)
	err = os.WriteFile(os.Getenv("HOME")+"/.git-credentials", []byte(cred), 0600)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "config", "--global", "credential.helper", "store")
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "clone", "--recurse-submodules", "--branch", stub.RepoBranch, "--single-branch", stub.GetRepoUrl(), RepoFolderName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

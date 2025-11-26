package controllers

import (
	"net/http"
	"portal/internal/server/github"
)

func SyncRemote() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, repo := range github.GithubClient.Repos {
			repo.Clone(*github.GithubClient)
		}
	}
}

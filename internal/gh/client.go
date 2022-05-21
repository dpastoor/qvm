package gh

import (
	"context"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

func NewClient(token string) *github.Client {
	if token == "" {
		return github.NewClient(nil)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

package gh

import (
	"context"
	"time"

	"github.com/google/go-github/v44/github"
	log "github.com/sirupsen/logrus"
)

func GetRelease(client *github.Client, release string) (*github.RepositoryRelease, error) {
	rel, resp, err := client.Repositories.GetReleaseByTag(
		context.Background(),
		"quarto-dev",
		"quarto-cli",
		release,
	)
	log.WithField("resp", resp).Trace("get-latest-release")
	if err != nil {
		return nil, err
	}
	return rel, err
}

func GetLatestRelease(client *github.Client) (*github.RepositoryRelease, error) {
	rel, resp, err := client.Repositories.GetLatestRelease(
		context.Background(),
		"quarto-dev",
		"quarto-cli",
	)
	log.WithField("resp", resp).Trace("get-latest-release")
	return rel, err
}

func GetReleases(client *github.Client, n int) ([]*github.RepositoryRelease, error) {
	// max of 50 per page
	perPage := 50
	remaining := n - perPage
	if n < 50 {
		remaining = 0
		perPage = n
	}
	var releases []*github.RepositoryRelease
	opts := &github.ListOptions{PerPage: perPage}
	for {
		start := time.Now()
		log.Tracef("perpage: %d, remaining: %d", opts.PerPage, remaining)
		rel, resp, err := client.Repositories.ListReleases(
			context.Background(),
			"quarto-dev",
			"quarto-cli",
			opts,
		)
		if err != nil {
			return releases, err
		}
		releases = append(releases, rel...)
		log.Tracef("repository release paginator: %s, page: %d", time.Since(start), resp.NextPage)
		if remaining <= 0 || resp.NextPage == 0 {
			break
		}
		if remaining <= perPage {
			opts.PerPage = remaining
			remaining = 0
		} else {
			remaining -= perPage
		}
		opts.Page = resp.NextPage
	}
	return releases, nil
}

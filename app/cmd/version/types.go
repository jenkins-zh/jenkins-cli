package cmd

import (
	"github.com/google/go-github/v29/github"
	"net/http"
)

// SelfUpgradeOption is the option for self upgrade command
type SelfUpgradeOption struct {
	ShowProgress bool

	GitHubClient *github.Client
	RoundTripper http.RoundTripper
}

package version

import (
	"github.com/google/go-github/v29/github"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"net/http"
)

// PrintOption is the version option
type PrintOption struct {
	Changelog  bool
	ShowLatest bool

	JenkinsClient    common.JenkinsClient
	JenkinsConfigMgr common.JenkinsConfigMgr
}

// SelfUpgradeOption is the option for self upgrade command
type SelfUpgradeOption struct {
	ShowProgress bool

	GitHubClient *github.Client
	RoundTripper http.RoundTripper
}

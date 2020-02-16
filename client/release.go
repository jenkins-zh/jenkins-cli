package client

import (
	"context"
	"github.com/google/go-github/v29/github"
)

// GitHubReleaseClient is the client of jcli github
type GitHubReleaseClient struct {
	Client *github.Client
}

// ReleaseAsset is the asset from GitHub release
type ReleaseAsset struct {
	TagName string
	Body    string
}

// Init init the GitHub client
func (g *GitHubReleaseClient) Init() {
	g.Client = github.NewClient(nil)
}

// GetLatestJCLIAsset returns the latest jcli asset
func (g *GitHubReleaseClient) GetLatestJCLIAsset() (*ReleaseAsset, error) {
	return g.GetLatestReleaseAsset("jenkins-zh", "jenkins-cli")
}

// GetLatestReleaseAsset returns the latest release asset
func (g *GitHubReleaseClient) GetLatestReleaseAsset(owner, repo string) (ra *ReleaseAsset, err error) {
	ctx := context.Background()

	var release *github.RepositoryRelease
	if release, _, err = g.Client.Repositories.GetLatestRelease(ctx, owner, repo); err == nil {
		ra = &ReleaseAsset{
			TagName: release.GetTagName(),
			Body:    release.GetBody(),
		}
	}
	return
}

// GetJCLIAsset returns the asset from a tag name
func (g *GitHubReleaseClient) GetJCLIAsset(tagName string) (*ReleaseAsset, error) {
	return g.GetReleaseAssetByTagName("jenkins-zh", "jenkins-cli", tagName)
}

// GetReleaseAssetByTagName returns the release asset by tag name
func (g *GitHubReleaseClient) GetReleaseAssetByTagName(owner, repo, tagName string) (ra *ReleaseAsset, err error) {
	ctx := context.Background()

	opt := &github.ListOptions{
		PerPage: 99999,
	}

	var releaseList []*github.RepositoryRelease
	if releaseList, _, err = g.Client.Repositories.ListReleases(ctx, owner, repo, opt); err == nil {
		for _, item := range releaseList {
			if item.GetTagName() == tagName {
				ra = &ReleaseAsset{
					TagName: item.GetTagName(),
					Body:    item.GetBody(),
				}
				break
			}
		}
	}
	return
}

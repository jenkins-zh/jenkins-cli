package client_test

import (
	jClient "github.com/jenkins-zh/jenkins-cli/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	ghClient := jClient.GitHubReleaseClient{}

	assert.Nil(t, ghClient.Client)
	ghClient.Init()
	assert.NotNil(t, ghClient.Client)
}

func TestGetLatestReleaseAsset(t *testing.T) {
	client, teardown := jClient.PrepareForGetLatestReleaseAsset() //setup()
	defer teardown()

	ghClient := jClient.GitHubReleaseClient{
		Client: client,
	}
	asset, err := ghClient.GetLatestReleaseAsset("o", "r")

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}

func TestGetLatestJCLIAsset(t *testing.T) {
	client, teardown := jClient.PrepareForGetLatestJCLIAsset() //setup()
	defer teardown()

	ghClient := jClient.GitHubReleaseClient{
		Client: client,
	}
	asset, err := ghClient.GetLatestJCLIAsset()

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}

func TestGetJCLIAsset(t *testing.T) {
	client, teardown := jClient.PrepareForGetJCLIAsset("tagName") //setup()
	defer teardown()

	ghClient := jClient.GitHubReleaseClient{
		Client: client,
	}
	asset, err := ghClient.GetJCLIAsset("tagName")

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}

func TestGetReleaseAssetByTagName(t *testing.T) {
	client, teardown := jClient.PrepareForGetReleaseAssetByTagName() //setup()
	defer teardown()

	ghClient := jClient.GitHubReleaseClient{
		Client: client,
	}
	asset, err := ghClient.GetReleaseAssetByTagName("jenkins-zh", "jenkins-cli", "tagName")

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}

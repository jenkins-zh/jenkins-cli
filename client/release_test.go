package client_test

import (
	"fmt"
	"github.com/google/go-github/v29/github"
	jClient "github.com/jenkins-zh/jenkins-cli/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	ghClient := jClient.GitHubReleaseClient{}

	assert.Nil(t, ghClient.Client)
	ghClient.Init()
	assert.NotNil(t, ghClient.Client)
}

func TestGetLatestReleaseAsset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":3, "body":"body", "tag_name":"tagName"}`)
	})

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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/jenkins-zh/jenkins-cli/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":3, "body":"body", "tag_name":"tagName"}`)
	})

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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/jenkins-zh/jenkins-cli/releases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":3, "body":"body", "tag_name":"tagName"}]`)
	})

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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/jenkins-zh/jenkins-cli/releases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":3, "body":"body", "tag_name":"tagName"}]`)
	})

	ghClient := jClient.GitHubReleaseClient{
		Client: client,
	}
	asset, err := ghClient.GetReleaseAssetByTagName("jenkins-zh", "jenkins-cli", "tagName")

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath = "/api-v3"
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
// this was copied from https://github.com/google/go-github
func setup() (client *github.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client = github.NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url

	return client, mux, server.URL, server.Close
}

// this was copied from https://github.com/google/go-github
func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

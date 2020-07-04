package client

import (
	"fmt"
	"github.com/google/go-github/v29/github"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// PrepareForGetJCLIAsset only for test
func PrepareForGetJCLIAsset(ver string) (client *github.Client, teardown func()) {
	var mux *http.ServeMux

	client, mux, _, teardown = setup()

	mux.HandleFunc("/repos/jenkins-zh/jenkins-cli/releases", func(w http.ResponseWriter, r *http.Request) {
		//testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, fmt.Sprintf(`[{"id":3, "body":"body", "tag_name":"%s"}]`, ver))
	})
	return
}

// PrepareForGetReleaseAssetByTagName only for test
func PrepareForGetReleaseAssetByTagName() (client *github.Client, teardown func()) {
	var mux *http.ServeMux

	client, mux, _, teardown = setup()

	mux.HandleFunc("/repos/jenkins-zh/jenkins-cli/releases", func(w http.ResponseWriter, r *http.Request) {
		//testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"id":3, "body":"body", "tag_name":"tagName"}]`)
	})
	return
}

// PrepareForGetLatestJCLIAsset only for test
func PrepareForGetLatestJCLIAsset() (client *github.Client, teardown func()) {
	var mux *http.ServeMux

	client, mux, _, teardown = setup()

	mux.HandleFunc("/repos/jenkins-zh/jenkins-cli/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		//testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"id":3, "body":"body", "tag_name":"tagName"}`)
	})
	return
}

// PrepareForGetLatestReleaseAsset only for test
func PrepareForGetLatestReleaseAsset() (client *github.Client, teardown func()) {
	var mux *http.ServeMux

	client, mux, _, teardown = setup()

	mux.HandleFunc("/repos/o/r/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		//testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"id":3, "body":"body", "tag_name":"tagName"}`)
	})
	return
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

	// uncomment here once it gets useful
	//apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	//	fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
	//	fmt.Fprintln(os.Stderr)
	//	fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
	//	fmt.Fprintln(os.Stderr)
	//	fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
	//	fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
	//	http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	//})

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
//func testMethod(t *testing.T, r *http.Request, want string) {
//	t.Helper()
//	if got := r.Method; got != want {
//		t.Errorf("Request method: %v, want %v", got, want)
//	}
//}

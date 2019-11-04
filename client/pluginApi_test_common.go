package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

// PrepareShowTrend only for test
func PrepareShowTrend(roundTripper *mhttp.MockRoundTripper, keyword string) (
	response *http.Response) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://plugins.jenkins.io/api/plugin/%s", keyword), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
		{"name":"fake","version": "0.1.8","url": "http://updates.jenkins-ci.org/download/plugins/hugo/0.1.8/hugo.hpi",
		"stats": {"installations":[{"total":1512},{"total":3472},{"total":4385},{"total":3981}]}}
		`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	return
}

// PrepareOnePluginInfo only for test
func PrepareOnePluginInfo(roundTripper *mhttp.MockRoundTripper, pluginName string) {
	PrepareShowTrend(roundTripper, pluginName)
}

// PrepareOnePluginWithDep only for test
func PrepareOnePluginWithDep(roundTripper *mhttp.MockRoundTripper, pluginName string) {
	response := PrepareShowTrend(roundTripper, pluginName)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{"name":"fake","version": "0.1.8","url": "http://updates.jenkins-ci.org/download/plugins/hugo/0.1.8/hugo.hpi",
		"dependencies":[{"name":"fake-1","optional":false,"version":"2.4"}]}
		`))

	fake1 := PrepareShowTrend(roundTripper, "fake-1")
	fake1.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{"name":"fake-1","version": "0.1.8","url": "http://updates.jenkins-ci.org/download/plugins/hugo/0.1.8/hugo.hpi"}
		`))
}

// PrepareOnePluginWithOptionalDep only for test
func PrepareOnePluginWithOptionalDep(roundTripper *mhttp.MockRoundTripper, pluginName string) {
	response := PrepareShowTrend(roundTripper, pluginName)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(`
		{"name":"fake","version": "0.1.8","url": "http://updates.jenkins-ci.org/download/plugins/hugo/0.1.8/hugo.hpi",
		"dependencies":[{"name":"fake-1","optional":true,"version":"2.4"}]}
		`))
}

// PrepareDownloadPlugin only for test
func PrepareDownloadPlugin(roundTripper *mhttp.MockRoundTripper) {
	request, _ := http.NewRequest("GET",
		"http://updates.jenkins-ci.org/download/plugins/hugo/0.1.8/hugo.hpi", nil)
	response := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
}

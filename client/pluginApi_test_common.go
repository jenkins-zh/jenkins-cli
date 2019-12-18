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
func PrepareDownloadPlugin(roundTripper *mhttp.MockRoundTripper) (response *http.Response) {
	request, _ := http.NewRequest("GET",
		"http://updates.jenkins-ci.org/download/plugins/hugo/0.1.8/hugo.hpi", nil)
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	return
}

// PrepareCheckUpdate only for test
func PrepareCheckUpdate(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string) {
	api := fmt.Sprintf("%s/pluginManager/checkUpdatesServer", rootURL)
	request, _ := http.NewRequest("POST", api, nil)
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareShowPlugins only for test
func PrepareShowPlugins(roundTripper *mhttp.MockRoundTripper, keyword string) (
	response *http.Response) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://plugins.jenkins.io/api/plugins/?q=%s&page=1&limit=1000", keyword), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
		"limit":1000,"page":1,"pages"":1,"total":1,
		"plugins":[{"name":"fake","version": "0.1.8","url": "http://updates.jenkins-ci.org/download/plugins/hugo/0.1.8/hugo.hpi",
		"stats": {"installations":[{"total":1512},{"total":3472},{"total":4385},{"total":3981}]},
		"securityWarnings":[{"versions":[{"firstVersion":null,"lastVersion":"0.1.8"}],"id":"SECURITY-659",
		"message":"XML External Entity (XXE) processing vulnerability","url":"https://jenkins.io/security/advisory/2018-02-05/","active":true}]}]
		`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	return
}

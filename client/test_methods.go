package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

// PrepareForEmptyAvaiablePluginList only for test
func PrepareForEmptyAvaiablePluginList(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/pluginManager/plugins", rootURL), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"status": "ok",
			"data": []
		}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	return
}

// PrepareForOneAvaiablePlugin only for test
func PrepareForOneAvaiablePlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/pluginManager/plugins", rootURL), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"status": "ok",
			"data": [{
				"name": "fake",
				"title": "fake"
			}]
		}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	return
}

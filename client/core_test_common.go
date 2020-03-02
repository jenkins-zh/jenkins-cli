package client

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"io/ioutil"
	"net/http"
)

// PrepareRestart only for test
func PrepareRestart(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, statusCode int) {
	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/safeRestart", rootURL), nil)
	response := PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	response.StatusCode = statusCode
	return
}

// PrepareRestartDirectly only for test
func PrepareRestartDirectly(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, statusCode int) {
	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/restart", rootURL), nil)
	response := PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	response.StatusCode = statusCode
	return
}

// PrepareForGetIdentity only for test
func PrepareForGetIdentity(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/instance", rootURL), nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
{"fingerprint":"fingerprint","publicKey":"publicKey","systemMessage":"systemMessage"}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)

	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
}

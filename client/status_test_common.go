package client

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"io/ioutil"
	"net/http"
)

//PrepareGetStatus only for test
func PrepareGetStatus(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/json", rootURL), nil)
	response := &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"nodeName":"master"}`)),
	}
	response.Header.Add("X-Jenkins", "version")
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)

	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
	return
}

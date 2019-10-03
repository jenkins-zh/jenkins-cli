package client

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

// PrepareGetUser only for test
func PrepareGetUser(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd string) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/user/%s/api/json", rootURL, user), nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"fullName":"admin"}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
	return
}
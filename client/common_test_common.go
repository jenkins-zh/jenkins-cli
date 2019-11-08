package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

// PrepareForGetIssuer only for test
func PrepareForGetIssuer(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s%s", rootURL, "/crumbIssuer/api/json"), nil)
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"CrumbRequestField":"CrumbRequestField","Crumb":"Crumb"}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
	return
}

// PrepareForGetIssuerWith500 only for test
func PrepareForGetIssuerWith500(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string) {
	_, response := PrepareForGetIssuer(roundTripper, rootURL, user, password)
	response.StatusCode = 500
}

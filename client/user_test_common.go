package client

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"strings"
	"net/url"
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

// PrepareCreateUser only for test
func PrepareCreateUser(roundTripper *mhttp.MockRoundTripper, rootURL,
	user, passwd, targetUserName string) {
	payload, _ := genSimpleUserAsPayload(targetUserName)
	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/securityRealm/createAccountByAdmin", rootURL), payload)
	PrepareCommonPost(request, roundTripper, user, passwd, rootURL)
}

// PrepareCreateToken only for test
func PrepareCreateToken(roundTripper *mhttp.MockRoundTripper, rootURL,
	user, passwd, newTokenName string) {
	formData := url.Values{}
	formData.Add("newTokenName", newTokenName)
	payload := strings.NewReader(formData.Encode())

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/user/%s/descriptorByName/jenkins.security.ApiTokenProperty/generateNewToken", rootURL, user), payload)
	response := PrepareCommonPost(request, roundTripper, user, passwd, rootURL)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(`
	{"status":"ok"}
	`))
}

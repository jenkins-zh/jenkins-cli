package client

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// PrepareGetUser only for test
func PrepareGetUser(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd string) (
	response *http.Response) {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user/%s/api/json", rootURL, user), nil)
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"fullName":"admin","description":"fake-description"}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
	return
}

// PrepareCreateUser only for test
func PrepareCreateUser(roundTripper *mhttp.MockRoundTripper, rootURL,
	user, passwd, targetUserName string) (response *http.Response) {
	payload, _ := genSimpleUserAsPayload(targetUserName, "fakePass")
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/securityRealm/createAccountByAdmin", rootURL), payload)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	response = PrepareCommonPost(request, "", roundTripper, user, passwd, rootURL)
	return
}

// PrepareCreateToken only for test
func PrepareCreateToken(roundTripper *mhttp.MockRoundTripper, rootURL,
	user, passwd, newTokenName, targetUser string) (response *http.Response) {
	formData := url.Values{}
	formData.Add("newTokenName", newTokenName)
	payload := strings.NewReader(formData.Encode())

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/user/%s/descriptorByName/jenkins.security.ApiTokenProperty/generateNewToken", rootURL, targetUser), payload)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	response = PrepareCommonPost(request, `{"status":"ok"}`, roundTripper, user, passwd, rootURL)
	return
}

// PrepareForEditUserDesc only for test
func PrepareForEditUserDesc(roundTripper *mhttp.MockRoundTripper, rootURL, userName, description, user, passwd string) {
	formData := url.Values{}
	formData.Add("description", description)
	payload := strings.NewReader(formData.Encode())

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/user/%s/submitDescription", rootURL, userName), payload)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	PrepareCommonPost(request, "", roundTripper, user, passwd, rootURL)
	return
}

// PrepareForDeleteUser only for test
func PrepareForDeleteUser(roundTripper *mhttp.MockRoundTripper, rootURL, userName, user, passwd string) (
	response *http.Response) {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/securityRealm/user/%s/doDelete", rootURL, userName), nil)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	response = PrepareCommonPost(request, "", roundTripper, user, passwd, rootURL)
	return
}

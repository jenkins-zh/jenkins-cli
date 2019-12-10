package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

// PrepareForGetCredentialList only for test
func PrepareForGetCredentialList(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, store string) {
	api := fmt.Sprintf("%s/credentials/store/%s/domain/_/api/json?pretty=true&depth=1", rootURL, store)
	request, _ := http.NewRequest("GET", api, nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(PrepareForCredentialListJson())),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
}

// PrepareForDeleteCredential only for test
func PrepareForDeleteCredential(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, store, id string) {
	api := fmt.Sprintf("%s/credentials/store/%s/domain/_/credential/%s/doDelete", rootURL, store, id)
	request, _ := http.NewRequest("POST", api, nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
}

// PrepareForCreateUsernamePasswordCredential only for test
func PrepareForCreateUsernamePasswordCredential(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, store string, cred UsernamePasswordCredential) {
	//api := fmt.Sprintf("%s/credentials/store/%s/domain/_/createCredentials", rootURL, store)
	//
	//formData := url.Values{}
	//formData.Add("json", fmt.Sprintf(`{"credentials": %s}`, credential))
	//payload := strings.NewReader(formData.Encode())
	//
	//request, _ := http.NewRequest("POST", api, nil)
	//response := &http.Response{
	//	StatusCode: 200,
	//	Request:    request,
	//	Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
	//}
	//roundTripper.EXPECT().
	//	RoundTrip(request).Return(response, nil)
	//if user != "" && password != "" {
	//	request.SetBasicAuth(user, password)
	//}
}

// PrepareForCreateSecretCredential only for test
func PrepareForCreateSecretCredential(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, store string) {

}

// PrepareForCredentialListJson only for test
func PrepareForCredentialListJson() string {
	return `{
  "_class" : "com.cloudbees.plugins.credentials.CredentialsStoreAction$DomainWrapper",
  "credentials" : [
    {
      "description" : "",
      "displayName" : "displayName",
      "fingerprint" : {
      },
      "fullName" : "system/_/19c27487-acca-4a39-9889-9ddd500388f3",
      "id" : "19c27487-acca-4a39-9889-9ddd500388f3",
      "typeName" : "Username with password"
    }
  ],
  "description" : "Credentials that should be available irrespective of domain specification to requirements matching.",
  "displayName" : "全局凭据 (unrestricted)",
  "fullDisplayName" : "系统 » 全局凭据 (unrestricted)",
  "fullName" : "system/_",
  "global" : true,
  "urlName" : "_"
}`
}

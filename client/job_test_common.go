package client

import (
	"bytes"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"github.com/jenkins-zh/jenkins-cli/util"
)

// PrepareForGetJobInputActions only for test
func PrepareForGetJobInputActions(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd, jobName string, buildID int) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/job/%s/%d/wfapi/pendingInputActions", rootURL, jobName, buildID), nil)
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
[{"id":"Eff7d5dba32b4da32d9a67a519434d3f","proceedText":"继续","message":"message","inputs":[],
"proceedUrl":"/job/test/5/wfapi/inputSubmit?inputId=Eff7d5dba32b4da32d9a67a519434d3f",
"abortUrl":"/job/test/5/input/Eff7d5dba32b4da32d9a67a519434d3f/abort","redirectApprovalUrl":"/job/test/5/input/"}]`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
	return
}

// PrepareForSubmitInput only for test
func PrepareForSubmitInput(roundTripper *mhttp.MockRoundTripper, rootURL, jobPath, user, passwd string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s%s/%d/input/%s/abort?json={\"parameter\":[]}", rootURL, jobPath, 1, "Eff7d5dba32b4da32d9a67a519434d3f"), nil)
	PrepareCommonPost(request, "", roundTripper, user, passwd, rootURL)
	return
}

// PrepareForSubmitProcessInput only for test
func PrepareForSubmitProcessInput(roundTripper *mhttp.MockRoundTripper, rootURL, jobPath, user, passwd string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s%s/%d/input/%s/proceed?json={\"parameter\":[]}", rootURL, jobPath, 1, "Eff7d5dba32b4da32d9a67a519434d3f"), nil)
	PrepareCommonPost(request, "", roundTripper, user, passwd, rootURL)
	return
}

// PrepareForBuildWithNoParams only for test
func PrepareForBuildWithNoParams(roundTripper *mhttp.MockRoundTripper, rootURL, jobName, user, passwd string) (
	request *http.Request, response *http.Response) {
	formData := url.Values{"json": {`{"parameter": []}`}}
	payload := strings.NewReader(formData.Encode())
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s/job/%s/build", rootURL, jobName), payload)
	request.Header.Add(util.ContentType, util.ApplicationForm)
	response = PrepareCommonPost(request, "", roundTripper, user, passwd, rootURL)
	response.StatusCode = 201
	return
}

// PrepareForBuildWithParams only for test
func PrepareForBuildWithParams(roundTripper *mhttp.MockRoundTripper, rootURL, jobName, user, passwd string) (
	request *http.Request, response *http.Response) {
	formData := url.Values{"json": {`{"parameter": {"Description":"","name":"name","Type":"","value":"value","DefaultParameterValue":{"Description":"","Value":null}}}`}}
	payload := strings.NewReader(formData.Encode())
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s/job/%s/build", rootURL, jobName), payload)
	request.Header.Add(util.ContentType, util.ApplicationForm)
	response = PrepareCommonPost(request, "", roundTripper, user, passwd, rootURL)
	response.StatusCode = 201
	return
}

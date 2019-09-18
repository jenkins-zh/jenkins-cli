package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
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
func PrepareForSubmitInput(roundTripper *mhttp.MockRoundTripper, rootURL, actionURL, user, passwd string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s%s", rootURL, actionURL), nil)
	PrepareCommonPost(request, roundTripper, user, passwd, rootURL)
	return
}

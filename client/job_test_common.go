package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"github.com/jenkins-zh/jenkins-cli/util"
)

// PrepareForGetJobInputActions only for test
func PrepareForGetJobInputActions(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, jobName string, buildID int) (
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

	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
	return
}

// PrepareForSubmitInput only for test
func PrepareForSubmitInput(roundTripper *mhttp.MockRoundTripper, rootURL, jobPath, user, password string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s%s/%d/input/%s/abort?json={\"parameter\":[]}", rootURL, jobPath, 1, "Eff7d5dba32b4da32d9a67a519434d3f"), nil)
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	return
}

// PrepareForSubmitProcessInput only for test
func PrepareForSubmitProcessInput(roundTripper *mhttp.MockRoundTripper, rootURL, jobPath, user, password string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s%s/%d/input/%s/proceed?json={\"parameter\":[]}", rootURL, jobPath, 1, "Eff7d5dba32b4da32d9a67a519434d3f"), nil)
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	return
}

// PrepareForBuildWithNoParams only for test
func PrepareForBuildWithNoParams(roundTripper *mhttp.MockRoundTripper, rootURL, jobName, user, password string) (
	request *http.Request, response *http.Response) {
	formData := url.Values{"json": {`{"parameter": []}`}}
	payload := strings.NewReader(formData.Encode())
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s/job/%s/build", rootURL, jobName), payload)
	request.Header.Add(util.ContentType, util.ApplicationForm)
	response = PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	response.StatusCode = 201
	return
}

// PrepareForBuildWithParams only for test
func PrepareForBuildWithParams(roundTripper *mhttp.MockRoundTripper, rootURL, jobName, user, password string) (
	request *http.Request, response *http.Response) {
	formData := url.Values{"json": {`{"parameter": {"Description":"","name":"name","Type":"","value":"value","DefaultParameterValue":{"Description":"","Value":null}}}`}}
	payload := strings.NewReader(formData.Encode())
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s/job/%s/build", rootURL, jobName), payload)
	request.Header.Add(util.ContentType, util.ApplicationForm)
	response = PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	response.StatusCode = 201
	return
}

// PrepareForGetJob only for test
func PrepareForGetJob(roundTripper *mhttp.MockRoundTripper, rootURL, jobName, user, password string) (
	response *http.Response) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/job/%s/api/json", rootURL, jobName), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
  "name" : "%s",
  "builds" : [
    {
      "number" : 1,
      "url" : "http://localhost:8080/job/we/1/"
    },
    {
      "number" : 2,
      "url" : "http://localhost:8080/job/we/2/"
    }]
				}`, jobName))),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
	return
}

// PrepareForGetJobWithParams only for test
func PrepareForGetJobWithParams(roundTripper *mhttp.MockRoundTripper, rootURL, jobName, user, password string) {
	response := PrepareForGetJob(roundTripper, rootURL, jobName, user, password)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
  "name" : "%s",
  "builds" : [
    {
      "number" : 1,
      "url" : "http://localhost:8080/job/we/1/"
    },
    {
      "number" : 2,
      "url" : "http://localhost:8080/job/we/2/"
    }],
  "property" : [
    {
      "_class" : "io.alauda.jenkins.devops.sync.WorkflowJobProperty"
    },
    {
      "parameterDefinitions" : [
        {
          "defaultParameterValue" : {
            "name" : "name",
            "value" : "jake"
          },
          "description" : "",
          "name" : "name",
          "type" : "StringParameterDefinition"
        }
      ]
    }
  ]
}`, jobName)))
}

// PrepareForGetBuild only for test
func PrepareForGetBuild(roundTripper *mhttp.MockRoundTripper, rootURL, jobName string, buildID int, user, password string) {
	api := ""
	if buildID == -1 {
		api = fmt.Sprintf("%s/job/%s/lastBuild/api/json", rootURL, jobName)
	} else {
		api = fmt.Sprintf("%s/job/%s/%d/api/json", rootURL, jobName, buildID)
	}
	request, _ := http.NewRequest("GET", api, nil)
	response := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"displayName":"fake"}
				`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
}

// PrepareForJobLog only for test
func PrepareForJobLog(roundTripper *mhttp.MockRoundTripper, rootURL, jobName string, buildID int, user, password string) {
	var api string
	if buildID == -1 {
		api = fmt.Sprintf("%s/job/%s/lastBuild/logText/progressiveText?start=%d", rootURL, jobName, 0)
	} else {
		api = fmt.Sprintf("%s/job/%s/%d/logText/progressiveText?start=%d", rootURL, jobName, buildID, 0)
	}
	request, _ := http.NewRequest("GET", api, nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Header: map[string][]string{
			"X-More-Data": []string{"false"},
			"X-Text-Size": []string{"8"},
		},
		Body: ioutil.NopCloser(bytes.NewBufferString("fake log")),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
}

// PrepareOneItem only for test
func PrepareOneItem(roundTripper *mhttp.MockRoundTripper, rootURL, name, kind, user, token string) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/items/list?name=%s&type=%s&start=%d&limit=%d",
		rootURL, name, kind, 0, 50), nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`[{"name":"fake","displayName":"fake","description":null,"type":"WorkflowJob","shortURL":"job/fake/","url":"job/fake/"}]`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && token != "" {
		request.SetBasicAuth(user, token)
	}
}

// PrepareEmptyItems only for test
func PrepareEmptyItems(roundTripper *mhttp.MockRoundTripper, rootURL, name, kind, user, token string) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/items/list?name=%s&type=%s&start=%d&limit=%d",
		rootURL, name, kind, 0, 50), nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	if user != "" && token != "" {
		request.SetBasicAuth(user, token)
	}
}

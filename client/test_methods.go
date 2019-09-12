package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

// PrepareForEmptyAvaiablePluginList only for test
func PrepareForEmptyAvaiablePluginList(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/pluginManager/plugins", rootURL), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"status": "ok",
			"data": []
		}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	return
}

// PrepareForOneAvaiablePlugin only for test
func PrepareForOneAvaiablePlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/pluginManager/plugins", rootURL), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"status": "ok",
			"data": [{
				"name": "fake",
				"title": "fake"
			}]
		}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	return
}

// PrepareForEmptyInstalledPluginList only for test
func PrepareForEmptyInstalledPluginList(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/pluginManager/api/json?depth=1", rootURL), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"plugins": []
			}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	return
}

// PrepareForOneInstalledPlugin only for test
func PrepareForOneInstalledPlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, response = PrepareForEmptyInstalledPluginList(roundTripper, rootURL)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(`{
			"plugins": [{
				"shortName": "fake",
				"version": "1.0",
				"hasUpdate": true,
				"enable": true,
				"active": true
			}]
		}`))
	return
}

// PrepareFor500InstalledPluginList only for test
func PrepareFor500InstalledPluginList(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, response = PrepareForEmptyInstalledPluginList(roundTripper, rootURL)
	response.StatusCode = 500
	return
}

// PrepareForUploadPlugin only for test
func PrepareForUploadPlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response, requestCrumb *http.Request, responseCrumb *http.Response) {
	tmpfile, _ := ioutil.TempFile("", "example")

	bytesBuffer := &bytes.Buffer{}
	writer := multipart.NewWriter(bytesBuffer)
	part, _ := writer.CreateFormFile("@name", filepath.Base(tmpfile.Name()))

	io.Copy(part, tmpfile)

	request, _ = http.NewRequest("POST", fmt.Sprintf("%s/pluginManager/uploadPlugin", rootURL), nil)
	request.Header.Add("CrumbRequestField", "Crumb")
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	// common crumb request
	requestCrumb, responseCrumb = RequestCrumb(roundTripper, rootURL)
	return
}

// PrepareForUninstallPlugin only for test
func PrepareForUninstallPlugin(roundTripper *mhttp.MockRoundTripper, rootURL, pluginName string) (
	request *http.Request, response *http.Response, requestCrumb *http.Request, responseCrumb *http.Response) {
	request, _ = http.NewRequest("POST", fmt.Sprintf("%s/pluginManager/plugin/%s/doUninstall", rootURL, pluginName), nil)
	request.Header.Add("CrumbRequestField", "Crumb")
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	// common crumb request
	requestCrumb, responseCrumb = RequestCrumb(roundTripper, rootURL)
	return
}

// PrepareForUninstallPluginWith500 only for test
func PrepareForUninstallPluginWith500(roundTripper *mhttp.MockRoundTripper, rootURL, pluginName string) (
	request *http.Request, response *http.Response, requestCrumb *http.Request, responseCrumb *http.Response) {
	request, response, requestCrumb, responseCrumb = PrepareForUninstallPlugin(roundTripper, rootURL, pluginName)
	response.StatusCode = 500
	return
}

// PrepareCancelQueue only for test
func PrepareCancelQueue(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd string) {
	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/queue/cancelItem?id=1", rootURL), nil)
	request.Header.Add("CrumbRequestField", "Crumb")
	response := &http.Response{
		StatusCode: 200,
		Header:     map[string][]string{},
		Proto:      "HTTP/1.1",
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	requestCrumb, _ := RequestCrumb(roundTripper, rootURL)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
		requestCrumb.SetBasicAuth(user, passwd)
	}
}

// PrepareGetQueue only for test
func PrepareGetQueue(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd string) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/queue/api/json", rootURL), nil)
	response := &http.Response{
		StatusCode: 200,
		Header:     map[string][]string{},
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
		{
			"_class" : "hudson.modexl.Queue",
			"discoverableItems" : [],
			"items" : [
			  {
				"actions" : [],
				"blocked" : false,
				"buildable" : true,
				"id" : 62,
				"inQueueSince" : 1567753826770,
				"params" : "",
				"stuck" : true,
				"task" : {
				  "_class" : "org.jenkinsci.plugins.workflow.support.steps.ExecutorStepExecution$PlaceholderTask"
				},
				"url" : "queue/item/62/",
				"why" : "等待下一个可用的执行器",
				"buildableStartMilliseconds" : 1567753826770,
				"pending" : false
			  }
			]
		  }`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
}

// RequestCrumb only for the test case
func RequestCrumb(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	requestCrumb *http.Request, responseCrumb *http.Response) {
	requestCrumb, _ = http.NewRequest("GET", fmt.Sprintf("%s%s", rootURL, "/crumbIssuer/api/json"), nil)
	responseCrumb = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    requestCrumb,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
		{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
		`)),
	}
	roundTripper.EXPECT().
		RoundTrip(requestCrumb).Return(responseCrumb, nil)
	return
}

// PrepareForInstallPlugin only for test
func PrepareForInstallPlugin(roundTripper *mhttp.MockRoundTripper, rootURL, pluginName, user, passwd string) {
	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/pluginManager/install?plugin.%s=", rootURL, pluginName), nil)
	request.Header.Add("CrumbRequestField", "Crumb")
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	// common crumb request
	requestCrumb, _ := RequestCrumb(roundTripper, rootURL)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
		requestCrumb.SetBasicAuth(user, passwd)
	}
	return
}

// PrepareForPipelineJob only for test
func PrepareForPipelineJob(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/job/test/restFul", rootURL), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"type":null,"displayName":null,"script":"script","sandbox":true}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
	return
}

// PrepareForUpdatePipelineJob only for test
func PrepareForUpdatePipelineJob(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd string) {
	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/job/test/restFul/update", rootURL), nil)
	request.Header.Add("CrumbRequestField", "Crumb")
	response := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	// common crumb request
	requestCrumb, _ := RequestCrumb(roundTripper, rootURL)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
		requestCrumb.SetBasicAuth(user, passwd)
	}
	return
}

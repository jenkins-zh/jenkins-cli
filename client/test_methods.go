package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
)

// PrepareForEmptyAvaiablePluginList only for test
func PrepareForEmptyAvaiablePluginList(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/pluginManager/plugins", rootURL), nil)
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"status": "ok",
			"data": []
		}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)
	return
}

// PrepareForOneAvaiablePlugin only for test
func PrepareForOneAvaiablePlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, response = PrepareForEmptyAvaiablePluginList(roundTripper, rootURL)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(`{
			"status": "ok",
			"data": [{
					"name": "fake",
					"title": "fake"
			}]
		}`))
	return
}

// PrepareForManyAvaiablePlugin only for test
func PrepareForManyAvaiablePlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, response = PrepareForEmptyAvaiablePluginList(roundTripper, rootURL)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(`{
			"status": "ok",
			"data": [
				{
					"name": "fake-ocean",
					"title": "fake-ocean"
				},
				{
					"name": "fake-ln",
					"title": "fake-ln"
				},
				{
					"name": "fake-is",
					"title": "fake-is"
				},
				{
					"name": "fake-oa",
					"title": "fake-oa"
				},
				{
					"name": "fake-open",
					"title": "fake-open"
				},
				{
					"name": "fake",
					"title": "fake"
				}
			]
		}`))
	return
}

// PrepareForEmptyInstalledPluginList only for test
func PrepareForEmptyInstalledPluginList(roundTripper *mhttp.MockRoundTripper, rootURL string, depth int) (
	request *http.Request, response *http.Response) {
	if depth > 1 {
		request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/pluginManager/api/json?depth=%d", rootURL, depth), nil)
	} else {
		request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/pluginManager/api/json?depth=1", rootURL), nil)
	}
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"plugins": []
			}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)
	return
}

// PrepareForOneInstalledPlugin only for test
func PrepareForOneInstalledPlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, response = PrepareForOneInstalledPluginWithPluginName(
		roundTripper, rootURL, "fake")
	return
}

// PrepareForManyInstalledPlugins only for test
func PrepareForManyInstalledPlugins(roundTripper *mhttp.MockRoundTripper, rootURL string, depth int) (
	request *http.Request, response *http.Response) {
	if depth > 1 {
		request, response = PrepareForEmptyInstalledPluginList(roundTripper, rootURL, depth)
	} else {
		request, response = PrepareForEmptyInstalledPluginList(roundTripper, rootURL, 1)
	}
	response.Body = ioutil.NopCloser(bytes.NewBufferString(`{
			"plugins": [
				{
					"shortName": "fake-ocean",
					"version": "1.18.111",
					"hasUpdate": false,
					"enable": true,
					"active": true,
                    "dependencies":[{"optional":true,"shortName":"fake-ln","version":"1.19"}]
				},
				{
					"shortName": "fake-ln",
					"version": "1.18.1",
					"hasUpdate": true,
					"enable": true,
					"active": true,
                    "dependencies":[{"optional":true,"shortName":"fake-is","version":"1.18.121-2.0"}]
				},
				{
					"shortName": "fake-is",
					"version": "1.18.131-2.0",
					"hasUpdate": true,
					"enable": true,
					"active": true,
                    "dependencies":[]
				},
				{
					"shortName": "fake",
					"version": "1.0",
					"hasUpdate": true,
					"enable": true,
					"active": true,
                    "dependencies":[]
				}
			]
		}`))
	return
}

// PrepareFor500InstalledPluginList only for test
func PrepareFor500InstalledPluginList(roundTripper *mhttp.MockRoundTripper, rootURL string, depth int) (
	request *http.Request, response *http.Response) {
	if depth > 1 {
		request, response = PrepareForEmptyInstalledPluginList(roundTripper, rootURL, depth)
	} else {
		request, response = PrepareForEmptyInstalledPluginList(roundTripper, rootURL, 1)
	}
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

	request, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/pluginManager/uploadPlugin", rootURL), nil)
	request.Header.Add("CrumbRequestField", "Crumb")
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	// common crumb request
	requestCrumb, responseCrumb = PrepareForGetIssuer(roundTripper, rootURL, "", "")
	return
}

// PrepareForUninstallPlugin only for test
func PrepareForUninstallPlugin(roundTripper *mhttp.MockRoundTripper, rootURL, pluginName string) (
	request *http.Request, response *http.Response, requestCrumb *http.Request, responseCrumb *http.Response) {
	request, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/pluginManager/plugin/%s/doUninstall", rootURL, pluginName), nil)
	request.Header.Add("CrumbRequestField", "Crumb")
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	// common crumb request
	requestCrumb, responseCrumb = PrepareForGetIssuer(roundTripper, rootURL, "", "")
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
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/queue/cancelItem?id=1", rootURL), nil)
	request.Header.Add("CrumbRequestField", "Crumb")
	response := &http.Response{
		StatusCode: 200,
		Header:     map[string][]string{},
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)
	PrepareForGetIssuer(roundTripper, rootURL, user, passwd)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
}

// PrepareGetQueue only for test
func PrepareGetQueue(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd string) {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/queue/api/json", rootURL), nil)
	response := &http.Response{
		StatusCode: 200,
		Header:     map[string][]string{},
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
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
}

// PrepareForRequestUpdateCenter only for the test case
func PrepareForRequestUpdateCenter(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	requestCenter *http.Request, responseCenter *http.Response) {
	requestCenter, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/updateCenter/site/default/api/json?pretty=true&depth=2", rootURL), nil)
	responseCenter = &http.Response{
		StatusCode: 200,
		Request:    requestCenter,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
		{
			"_class": "hudson.model.UpdateSite",
			"connectionCheckUrl": "http://www.google.com/",
			"dataTimestamp": 1567999067717,
			"hasUpdates": true,
			"id": "default",
			"updates": [{
				"name": "fake-ocean",
				"sourceId": "default",
				"requiredCore": "2.138.4",
				"version": "1.19.011",
				"title": "fake-ocean",
				"sourceId": "default",
				"installed": {
					"active": true,
					"backupVersion": "1.17.011",
					"hasUpdate": true,
					"version": "1.18.111"
				}
			},{
				"name": "fake-ln",
				"sourceId": "default",
				"requiredCore": "2.138.4",
				"version": "1.19.011",
				"title": "fake-ln",
				"sourceId": "default",
				"installed": {
					"active": true,
					"hasUpdate": true,
					"version": "1.18.1"
				}
			},{
				"name": "fake-is",
				"sourceId": "default",
				"requiredCore": "2.138.4",
				"version": "1.19.1",
				"title": "fake-is",
				"sourceId": "default",
				"installed": {
					"active": true,
					"backupVersion": "1.17.011",
					"hasUpdate": true,
					"version": "1.18.111"
				}
			}
			],
			"availables": [{
				"name": "fake-oa",
				"sourceId": "default",
				"requiredCore": "2.138.4",
				"version": "1.13.011",
				"title": "fake-oa",
				"installed": null
			},{
				"name": "fake-open",
				"sourceId": "default",
				"requiredCore": "2.138.4",
				"version": "1.13.0",
				"title": "fake-open",
				"installed": null
			}
			],
			"url": "https://updates.jenkins.io/update-center.json"
		}
		`)),
	}
	roundTripper.EXPECT().RoundTrip(NewRequestMatcher(requestCenter)).Return(responseCenter, nil)
	return
}

// PrepareForNoAvailablePlugins only for the test case
func PrepareForNoAvailablePlugins(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	requestCenter *http.Request, responseCenter *http.Response) {
	requestCenter, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/updateCenter/site/default/api/json?pretty=true&depth=2", rootURL), nil)
	responseCenter = &http.Response{
		StatusCode: 200,
		Request:    requestCenter,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
		{
			"_class": "hudson.model.UpdateSite",
			"connectionCheckUrl": "http://www.google.com/",
			"dataTimestamp": 1567999067717,
			"hasUpdates": true,
			"id": "default",
			"updates": [
			],
			"availables": [
			],
			"url": "https://updates.jenkins.io/update-center.json"
		}
		`)),
	}
	roundTripper.EXPECT().RoundTrip(NewRequestMatcher(requestCenter)).Return(responseCenter, nil)
	return
}

// PrepareForRequest500UpdateCenter only for the test case
func PrepareForRequest500UpdateCenter(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	requestCenter *http.Request, responseCenter *http.Response) {
	requestCenter, responseCenter = PrepareForNoAvailablePlugins(roundTripper, rootURL)
	responseCenter.StatusCode = 500
	return
}

// PrepareForInstallPlugin only for test
func PrepareForInstallPlugin(roundTripper *mhttp.MockRoundTripper, rootURL, pluginName, user, passwd string) {
	PrepareForInstallPluginWithCode(roundTripper, 200, rootURL, pluginName, user, passwd)
}

// PrepareForInstallPluginWithVersion only for test
func PrepareForInstallPluginWithVersion(roundTripper *mhttp.MockRoundTripper, rootURL, pluginName, version, user, passwd string) {
	PrepareForInstallPluginWithCode(roundTripper, 200, rootURL, pluginName+"@"+version, user, passwd)
}

// PrepareForInstallPluginWithCode only for test
func PrepareForInstallPluginWithCode(roundTripper *mhttp.MockRoundTripper,
	statusCode int, rootURL, pluginName, user, passwd string) (response *http.Response) {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/pluginManager/install?plugin.%s=", rootURL, pluginName), nil)
	request.Header.Add("CrumbRequestField", "Crumb")
	response = &http.Response{
		StatusCode: statusCode,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	// common crumb request
	requestCrumb, _ := PrepareForGetIssuer(roundTripper, rootURL, user, passwd)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
		requestCrumb.SetBasicAuth(user, passwd)
	}
	return
}

// PrepareForPipelineJob only for test
func PrepareForPipelineJob(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd string) (
	request *http.Request, response *http.Response) {
	request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/job/test/restFul", rootURL), nil)
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"type":null,"displayName":null,"script":"script","sandbox":true}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
	return
}

// PrepareForUpdatePipelineJob only for test
func PrepareForUpdatePipelineJob(roundTripper *mhttp.MockRoundTripper, rootURL, script, user, password string) {
	formData := url.Values{}
	formData.Add("script", script)
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/job/test/restFul/update?%s", rootURL, formData.Encode()), nil)
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	return
}

// PrepareForCreatePipelineJob only for test
func PrepareForCreatePipelineJob(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, jobPayload CreateJobPayload) {
	playLoadData, _ := json.Marshal(jobPayload)
	formData := url.Values{
		"json": {string(playLoadData)},
		"name": {jobPayload.Name},
		"mode": {jobPayload.Mode},
		"from": {jobPayload.From},
	}
	payload := strings.NewReader(formData.Encode())

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/view/all/createItem", rootURL), payload)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	return
}

// PrepareCommonPost only for test
func PrepareCommonPost(request *http.Request, responseBody string, roundTripper *mhttp.MockRoundTripper, user, passwd, rootURL string) (
	response *http.Response) {
	request.Header.Add("CrumbRequestField", "Crumb")
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}
	roundTripper.EXPECT().
		RoundTrip(NewVerboseRequestMatcher(request).WithBody().WithQuery()).Return(response, nil)

	// common crumb request
	PrepareForGetIssuer(roundTripper, rootURL, user, passwd)
	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
	return
}

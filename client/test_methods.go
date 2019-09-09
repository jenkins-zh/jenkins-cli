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
	request *http.Request, response *http.Response, requestCenter *http.Request, responseCenter *http.Response, requestAvliable *http.Request, responseAvliable *http.Response) {
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
	requestCenter, responseCenter = RequestUpdateCenter(roundTripper, rootURL)
	requestAvliable, responseAvliable = PrepareForOneInstalledPlugin(roundTripper, rootURL)
	return
}

// PrepareForOneAvaiablePlugin only for test
func PrepareForOneAvaiablePlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response, requestCenter *http.Request, responseCenter *http.Response, requestAvliable *http.Request, responseAvliable *http.Response) {
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
	requestCenter, responseCenter = RequestUpdateCenter(roundTripper, rootURL)
	requestAvliable, responseAvliable = PrepareForOneInstalledPlugin(roundTripper, rootURL)
	return
}

// PrepareForManyAvaiablePlugin only for test
func PrepareForManyAvaiablePlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response, requestCenter *http.Request, responseCenter *http.Response, requestAvliable *http.Request, responseAvliable *http.Response) {
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/pluginManager/plugins", rootURL), nil)
	response = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`{
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
		}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
	requestCenter, responseCenter = RequestUpdateCenter(roundTripper, rootURL)
	requestAvliable, responseAvliable = PrepareForOneInstalledPlugin(roundTripper, rootURL)
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

// RequestUpdateCenter only for the test case
func RequestUpdateCenter(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	requestCenter *http.Request, responseCenter *http.Response) {
	requestCenter, _ = http.NewRequest("GET", fmt.Sprintf("%s/updateCenter/site/default/api/json?pretty=true&depth=2", rootURL), nil)
	responseCenter = &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
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
	roundTripper.EXPECT().RoundTrip(requestCenter).Return(responseCenter, nil)
	return
}

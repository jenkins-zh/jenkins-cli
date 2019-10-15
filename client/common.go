package client

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/jenkins-zh/jenkins-cli/util"
)

// JenkinsCore core informations of Jenkins
type JenkinsCore struct {
	JenkinsCrumb
	URL       string
	UserName  string
	Token     string
	Proxy     string
	ProxyAuth string

	Debug        bool
	Output       io.Writer
	RoundTripper http.RoundTripper
}

// JenkinsCrumb crumb for Jenkins
type JenkinsCrumb struct {
	CrumbRequestField string
	Crumb             string
}

// GetClient get the default http Jenkins client
func (j *JenkinsCore) GetClient() (client *http.Client) {
	var roundTripper http.RoundTripper
	if j.RoundTripper != nil {
		roundTripper = j.RoundTripper
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		if err := util.SetProxy(j.Proxy, j.ProxyAuth, tr); err != nil {
			log.Fatal(err)
		}
		roundTripper = tr
	}
	client = &http.Client{Transport: roundTripper}
	return
}

// ProxyHandle takes care of the proxy setting
func (j *JenkinsCore) ProxyHandle(request *http.Request) {
	if j.ProxyAuth != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(j.ProxyAuth))
		request.Header.Add("Proxy-Authorization", basicAuth)
	}
}

// AuthHandle takes care of the auth
func (j *JenkinsCore) AuthHandle(request *http.Request) (err error) {
	if j.UserName != "" && j.Token != "" {
		request.SetBasicAuth(j.UserName, j.Token)
	}

	// not add the User-Agent for tests
	if j.RoundTripper == nil {
		request.Header.Set("User-Agent", app.GetCombinedVersion())
	}

	j.ProxyHandle(request)

	// all post request to Jenkins must be has the crumb
	if request.Method == "POST" {
		err = j.CrumbHandle(request)
	}
	return
}

// CrumbHandle handle crum with http request
func (j *JenkinsCore) CrumbHandle(request *http.Request) error {
	if c, err := j.GetCrumb(); err == nil && c != nil {
		// cannot get the crumb could be a normal situation
		j.CrumbRequestField = c.CrumbRequestField
		j.Crumb = c.Crumb
		request.Header.Add(j.CrumbRequestField, j.Crumb)
	} else {
		return err
	}

	return nil
}

// GetCrumb get the crumb from Jenkins
func (j *JenkinsCore) GetCrumb() (crumbIssuer *JenkinsCrumb, err error) {
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = j.Request("GET", "/crumbIssuer/api/json", nil, nil); err == nil {
		if statusCode == 200 {
			err = json.Unmarshal(data, &crumbIssuer)
		} else if statusCode == 404 {
			// return 404 if Jenkins does no have crumb
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
		}
	}

	return
}

// RequestWithData requests the api and parse the data into an interface
func (j *JenkinsCore) RequestWithData(method, api string, headers map[string]string,
	payload io.Reader, successCode int, obj interface{}) (err error) {
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = j.Request(method, api, headers, payload); err == nil {
		if statusCode == successCode {
			err = json.Unmarshal(data, obj)
		} else {
			err = j.ErrorHandle(statusCode, data)
		}
	}
	return
}

// RequestWithoutData requests the api without handling data
func (j *JenkinsCore) RequestWithoutData(method, api string, headers map[string]string,
	payload io.Reader, successCode int) (statusCode int, err error) {
	var (
		data []byte
	)

	if statusCode, data, err = j.Request(method, api, headers, payload); err == nil &&
		statusCode != successCode {
		err = j.ErrorHandle(statusCode, data)
	}
	return
}

// ErrorHandle handles the error cases
func (j *JenkinsCore) ErrorHandle(statusCode int, data []byte) (err error) {
	if statusCode >= 400 && statusCode < 500 {
		err = j.PermissionError(statusCode)
	} else {
		err = fmt.Errorf("unexpected status code: %d", statusCode)
	}
	if j.Debug {
		ioutil.WriteFile("debug.html", data, 0664)
	}
	return
}

// PermissionError handles the no permission
func (j *JenkinsCore) PermissionError(statusCode int) (err error) {
	if statusCode == 404 {
		err = fmt.Errorf("Not found resources")
	} else {
		err = fmt.Errorf("The current user no permission")
	}
	return
}

// RequestWithResponseHeader make a common request
func (j *JenkinsCore) RequestWithResponseHeader(method, api string, headers map[string]string, payload io.Reader, obj interface{}) (
	response *http.Response, err error){
	response, err = j.RequestWithResponse(method, api, headers, payload)
	if err != nil {
		return
	}

	var data []byte
	if response.StatusCode == 200 {
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			err = json.Unmarshal(data, obj)
		}
	}
	return
}

// RequestWithResponse make a common request
func (j *JenkinsCore) RequestWithResponse(method, api string, headers map[string]string, payload io.Reader) (
	response *http.Response, err error) {
	var (
		req *http.Request
	)

	if req, err = http.NewRequest(method, fmt.Sprintf("%s%s", j.URL, api), payload); err != nil {
		return
	}
	if err = j.AuthHandle(req); err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := j.GetClient()
	return client.Do(req)
}

// Request make a common request
func (j *JenkinsCore) Request(method, api string, headers map[string]string, payload io.Reader) (
	statusCode int, data []byte, err error) {
	var (
		req      *http.Request
		response *http.Response
	)

	if req, err = http.NewRequest(method, fmt.Sprintf("%s%s", j.URL, api), payload); err != nil {
		return
	}
	if err = j.AuthHandle(req); err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := j.GetClient()
	if response, err = client.Do(req); err == nil {
		statusCode = response.StatusCode
		data, err = ioutil.ReadAll(response.Body)
	}
	return
}

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
	"strconv"
	"time"

	"go.uber.org/zap"
	"moul.io/http2curl"

	"github.com/jenkins-zh/jenkins-cli/util"
	ext "github.com/linuxsuren/cobra-extension/version"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
)

// language is for global Accept Language
var language string

// SetLanguage set the language
func SetLanguage(lan string) {
	language = lan
}

// JenkinsCore core information of Jenkins
type JenkinsCore struct {
	JenkinsCrumb
	Timeout            time.Duration
	URL                string
	InsecureSkipVerify bool
	UserName           string
	Token              string
	Proxy              string
	ProxyAuth          string

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
			TLSClientConfig: &tls.Config{InsecureSkipVerify: j.InsecureSkipVerify},
		}
		if err := httpdownloader.SetProxy(j.Proxy, j.ProxyAuth, tr); err != nil {
			log.Fatal(err)
		}
		roundTripper = tr
	}

	// make sure have a default timeout here
	if j.Timeout <= 0 {
		j.Timeout = 15
	}

	client = &http.Client{
		Transport: roundTripper,
		Timeout:   j.Timeout * time.Second,
	}
	return
}

// ProxyHandle takes care of the proxy setting
func (j *JenkinsCore) ProxyHandle(request *http.Request) {
	if j.ProxyAuth != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(j.ProxyAuth))
		logger.Debug("setting proxy for HTTP request", zap.String("header", basicAuth))
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
		request.Header.Set("User-Agent", ext.GetCombinedVersion())
	}

	j.ProxyHandle(request)

	// all post request to Jenkins must be has the crumb
	if request.Method == http.MethodPost {
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

	if statusCode, data, err = j.Request(http.MethodGet, "/crumbIssuer/api/json", nil, nil); err == nil {
		if statusCode == 200 {
			err = json.Unmarshal(data, &crumbIssuer)
		} else if statusCode == 404 {
			// return 404 if Jenkins does no have crumb
			//err = fmt.Errorf("crumb is disabled")
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

// RequestWithoutData requests the api without handling data
func (j *JenkinsCore) RequestWithDataResponse(method, api string, headers map[string]string,
	payload io.Reader, successCode int) (JenkinsBuildState, error) {
	var (
		data  []byte
		state JenkinsBuildState
	)

	if response, err := j.RequestWithResponse(method, api, headers, payload); err == nil {
		statusCode := response.StatusCode
		data, _ = ioutil.ReadAll(response.Body)
		if statusCode == successCode {
			state.BodyData = data
			state.StatusCode = response.StatusCode
			if len(response.Header.Get("Location")) > 0 {
				locationSlice := util.ArraySplitAndDeleteEmpty(response.Header.Get("Location"), "/")
				queueId := locationSlice[len(locationSlice)-1]
				if len(queueId) > 0 {
					if state.QueueId, err = strconv.ParseInt(queueId, 10, 64); err != nil {
						logger.Error("request job run queue error", zap.String("queue id", queueId))
						return state, err
					}
				}
			}
		} else {
			err = j.ErrorHandle(statusCode, data)
			return state, err
		}
	} else {
		return state, err
	}

	return state, nil
}

// ErrorHandle handles the error cases
func (j *JenkinsCore) ErrorHandle(statusCode int, data []byte) (err error) {
	if statusCode >= 400 && statusCode < 500 {
		err = j.PermissionError(statusCode)
	} else {
		err = fmt.Errorf("unexpected status code: %d", statusCode)
	}

	logger.Debug("get response", zap.String("data", string(data)))
	return
}

// PermissionError handles the no permission
func (j *JenkinsCore) PermissionError(statusCode int) (err error) {
	switch statusCode {
	case 400:
		err = fmt.Errorf("bad request, code %d", statusCode)
	case 404:
		err = fmt.Errorf("not found resources")
	default:
		err = fmt.Errorf("the current user has not permission, code %d", statusCode)
	}
	return
}

// RequestWithResponseHeader make a common request
func (j *JenkinsCore) RequestWithResponseHeader(method, api string, headers map[string]string, payload io.Reader, obj interface{}) (
	response *http.Response, err error) {
	response, err = j.RequestWithResponse(method, api, headers, payload)

	if err == nil && obj != nil && response.StatusCode == 200 {
		var data []byte
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

	if curlCmd, curlErr := http2curl.GetCurlCommand(req); curlErr == nil {
		logger.Debug("HTTP request as curl", zap.String("cmd", curlCmd.String()))
	}
	return client.Do(req)
}

// Request make a common request
func (j *JenkinsCore) Request(method, api string, headers map[string]string, payload io.Reader) (
	statusCode int, data []byte, err error) {
	var (
		req        *http.Request
		response   *http.Response
		requestURL string
	)

	if requestURL, err = util.URLJoinAsString(j.URL, api); err != nil {
		err = fmt.Errorf("cannot parse the URL of Jenkins, error is %v", err)
		return
	}

	logger.Debug("send HTTP request", zap.String("URL", requestURL), zap.String("method", method))
	if req, err = http.NewRequest(method, requestURL, payload); err != nil {
		return
	}
	if language != "" {
		req.Header.Set("Accept-Language", language)
	}
	if err = j.AuthHandle(req); err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if curlCmd, curlErr := http2curl.GetCurlCommand(req); curlErr == nil {
		logger.Debug("HTTP request as curl", zap.String("cmd", curlCmd.String()))
	}

	client := j.GetClient()
	if response, err = client.Do(req); err == nil {
		statusCode = response.StatusCode
		data, err = ioutil.ReadAll(response.Body)
	}
	return
}

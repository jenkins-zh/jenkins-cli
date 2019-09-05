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
	"net/url"
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

func (j *JenkinsCore) GetClient() (client *http.Client) {
	var roundTripper http.RoundTripper
	if j.RoundTripper != nil {
		roundTripper = j.RoundTripper
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		if j.Proxy != "" {
			if proxyURL, err := url.Parse(j.Proxy); err == nil {
				tr.Proxy = http.ProxyURL(proxyURL)
			} else {
				log.Fatal(err)
			}

			if j.ProxyAuth != "" {
				basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(j.ProxyAuth))
				tr.ProxyConnectHeader = http.Header{}
				tr.ProxyConnectHeader.Add("Proxy-Authorization", basicAuth)
			}
		}
		roundTripper = tr
	}
	client = &http.Client{Transport: roundTripper}
	return
}

func (j *JenkinsCore) ProxyHandle(request *http.Request) {
	if j.ProxyAuth != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(j.ProxyAuth))
		request.Header.Add("Proxy-Authorization", basicAuth)
	}
}

func (j *JenkinsCore) AuthHandle(request *http.Request) (err error) {
	if j.UserName != "" && j.Token != "" {
		request.SetBasicAuth(j.UserName, j.Token)
	}

	j.ProxyHandle(request)

	if request.Method == "POST" {
		err = j.CrumbHandle(request)
	}
	return
}

// CrumbHandle handle crum with http request
func (j *JenkinsCore) CrumbHandle(request *http.Request) error {
	if c, err := j.GetCrumb(); err == nil && c != nil {
		// cannot get the crumb could be a noraml situation
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
			json.Unmarshal(data, &crumbIssuer)
		} else if statusCode == 404 {
			// return 404 if Jenkins does no have crumb
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
		}
	}

	return
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
	j.AuthHandle(req)

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

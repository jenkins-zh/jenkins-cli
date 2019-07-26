package client

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type JenkinsCore struct {
	JenkinsCrumb
	URL       string
	UserName  string
	Token     string
	Proxy     string
	ProxyAuth string

	Debug bool
}

type JenkinsCrumb struct {
	CrumbRequestField string
	Crumb             string
}

func (j *JenkinsCore) GetClient() (client *http.Client) {
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
	client = &http.Client{Transport: tr}
	return
}

func (j *JenkinsCore) ProxyHandle(request *http.Request) {
	if j.ProxyAuth != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(j.ProxyAuth))
		request.Header.Add("Proxy-Authorization", basicAuth)
	}
}

func (j *JenkinsCore) AuthHandle(request *http.Request) (err error) {
	request.SetBasicAuth(j.UserName, j.Token)
	j.ProxyHandle(request)

	if request.Method == "POST" {
		err = j.CrumbHandle(request)
	}
	return
}

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

func (j *JenkinsCore) GetCrumb() (*JenkinsCrumb, error) {
	api := fmt.Sprintf("%s/crumbIssuer/api/json", j.URL)

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}
	j.AuthHandle(req)

	var crumbIssuer JenkinsCrumb
	client := j.GetClient()
	if response, err := client.Do(req); err == nil {
		if data, err := ioutil.ReadAll(response.Body); err == nil {
			if response.StatusCode == 200 {
				json.Unmarshal(data, &crumbIssuer)
			} else if response.StatusCode == 404 {
				return nil, err
			} else {
				log.Printf("Unexpected status code: %d.", response.StatusCode)
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
	return &crumbIssuer, nil
}

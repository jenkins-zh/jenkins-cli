package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JenkinsCore struct {
	JenkinsCrumb
	URL      string
	UserName string
	Token    string
}

type JenkinsCrumb struct {
	CrumbRequestField string
	Crumb             string
}

func (j *JenkinsCore) AuthHandle(request *http.Request) {
	request.SetBasicAuth(j.UserName, j.Token)
}

func (j *JenkinsCore) CrumbHandle(request *http.Request) error {
	if c, err := j.GetCrumb(); err == nil {
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if response, err := client.Do(req); err == nil {
		if data, err := ioutil.ReadAll(response.Body); err == nil {
			if response.StatusCode == 200 {
				json.Unmarshal(data, &crumbIssuer)
			} else {
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

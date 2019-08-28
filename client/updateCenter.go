package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/util"
)

// UpdateCenterManager manages the UpdateCenter
type UpdateCenterManager struct {
	JenkinsCore
}

// UpdateCenter represents the update center of Jenkins
type UpdateCenter struct {
	Availables                   []Plugin
	Jobs                         []InstallationJob
	RestartRequiredForCompletion bool
	Sites                        []CenterSite
}

// UpdateCenterJob represents the job for updateCenter which execute a task
type UpdateCenterJob struct {
	ErrorMessage string
	ID           int `json:"id"`
	Type         string
}

type InstallationJob struct {
	UpdateCenterJob

	Name   string
	Status InstallationJobStatus
}

type InstallationJobStatus struct {
	Success bool
	Type    string
}

// CenterSite represents the site of update center
type CenterSite struct {
	ConnectionCheckURL string `json:"connectionCheckUrl"`
	HasUpdates         bool
	ID                 string `json:"id"`
	URL                string `json:"url"`
}

type InstallStates struct {
	Data   InstallStatesData
	Status string
}

type InstallStatesData struct {
	Jobs  InstallStatesJob
	State string
}

type InstallStatesJob struct {
	InstallStatus   string
	Name            string
	RequiresRestart string
	Title           string
	Version         string
}

func (u *UpdateCenterManager) Status() (status *UpdateCenter, err error) {
	api := fmt.Sprintf("%s/updateCenter/api/json?pretty=false&depth=1", u.URL)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("GET", api, nil)
	if err == nil {
		u.AuthHandle(req)
	} else {
		return
	}

	client := u.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			if err == nil {
				status = &UpdateCenter{}
				err = json.Unmarshal(data, status)
			}
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

// Upgrade the Jenkins core
func (u *UpdateCenterManager) Upgrade() (err error) {
	api := fmt.Sprintf("%s/updateCenter/upgrade", u.URL)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("POST", api, nil)
	if err == nil {
		if err = u.AuthHandle(req); err != nil {
			log.Fatal(err)
		}
	} else {
		return
	}

	client := u.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code != 200 {
			fmt.Println("status code", code)
		}
		if err == nil && u.Debug && len(data) > 0 {
			ioutil.WriteFile("debug.html", data, 0664)
		}
	} else {
		log.Fatal(err)
	}
	return
}

// DownloadJenkins download Jenkins
func (u *UpdateCenterManager) DownloadJenkins(lts bool, output string) (err error) {
	var url string
	if lts {
		url = "http://mirrors.jenkins.io/war-stable/latest/jenkins.war"
	} else {
		url = "http://mirrors.jenkins.io/war/latest/jenkins.war"
	}

	downloader := util.HTTPDownloader{
		TargetFilePath: output,
		URL:            url,
		ShowProgress:   true,
	}
	err = downloader.DownloadFile()
	return
}

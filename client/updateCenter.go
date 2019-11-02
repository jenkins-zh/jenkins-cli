package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/util"
)

// UpdateCenterManager manages the UpdateCenter
type UpdateCenterManager struct {
	MirrorSite string

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

// InstallationJob represents the installation job
type InstallationJob struct {
	UpdateCenterJob

	Name   string
	Status InstallationJobStatus
}

// InstallationJobStatus represents the installation job status
type InstallationJobStatus struct {
	Success bool
	Type    string
}

// CenterSite represents the site of update center
type CenterSite struct {
	AvailablesPlugins  []CenterPlugin `json:"availables"`
	ConnectionCheckURL string         `json:"connectionCheckUrl"`
	DataTimestamp      int64          `json:"dataTimestamp"`
	HasUpdates         bool           `json:"hasUpdates"`
	ID                 string         `json:"id"`
	UpdatePlugins      []CenterPlugin `json:"updates"`
	URL                string         `json:"url"`
}

// InstallStates is the installation states
type InstallStates struct {
	Data   InstallStatesData
	Status string
}

// InstallStatesData is the installation state data
type InstallStatesData struct {
	Jobs  InstallStatesJob
	State string
}

// InstallStatesJob is the installation state job
type InstallStatesJob struct {
	InstallStatus   string
	Name            string
	RequiresRestart string
	Title           string
	Version         string
}

// CenterPlugin represents the all plugin from UpdateCenter
type CenterPlugin struct {
	CompatibleWithInstalledVersion bool
	Excerpt                        string
	Installed                      InstalledPlugin
	MinimumJavaVersion             string
	Name                           string
	RequiredCore                   string
	SourceID                       string
	Title                          string
	URL                            string
	Version                        string
	Wiki                           string
}

// Status returns the status of Jenkins
func (u *UpdateCenterManager) Status() (status *UpdateCenter, err error) {
	err = u.RequestWithData("GET", "/updateCenter/api/json?pretty=false&depth=1", nil, nil, 200, &status)
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
func (u *UpdateCenterManager) DownloadJenkins(lts, showProgress bool, output string) (err error) {
	var url string
	if lts {
		url = fmt.Sprintf("%s/war-stable/latest/jenkins.war", strings.TrimRight(u.MirrorSite, "/"))
	} else {
		url = fmt.Sprintf("%s/war/latest/jenkins.war", strings.TrimRight(u.MirrorSite, "/"))
	}

	downloader := util.HTTPDownloader{
		RoundTripper:   u.RoundTripper,
		TargetFilePath: output,
		URL:            url,
		ShowProgress:   showProgress,
	}
	err = downloader.DownloadFile()
	return
}

// GetSite is get Available Plugins and Updated Plugins from UpdateCenter
func (u *UpdateCenterManager) GetSite() (site *CenterSite, err error) {
	err = u.RequestWithData("GET", "/updateCenter/site/default/api/json?pretty=true&depth=2", nil, nil, 200, &site)
	return
}

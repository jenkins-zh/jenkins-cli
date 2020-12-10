package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
)

// UpdateCenterManager manages the UpdateCenter
type UpdateCenterManager struct {
	JenkinsCore

	MirrorSite string

	LTS     bool
	Version string
	Output  string

	Formula string

	ShowProgress bool
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
	err = u.RequestWithData(http.MethodGet, "/updateCenter/api/json?pretty=false&depth=1", nil, nil, 200, &status)
	return
}

// Upgrade the Jenkins core
func (u *UpdateCenterManager) Upgrade() (err error) {
	_, err = u.RequestWithoutData(http.MethodPost, "/updateCenter/upgrade",
		nil, nil, 200)
	return
}

// DownloadJenkins download Jenkins
func (u *UpdateCenterManager) DownloadJenkins() (err error) {
	showProgress, output := u.ShowProgress, u.Output
	warURL := u.GetJenkinsWarURL()

	downloader := httpdownloader.HTTPDownloader{
		RoundTripper:   u.RoundTripper,
		TargetFilePath: output,
		URL:            warURL,
		ShowProgress:   showProgress,
	}
	err = downloader.DownloadFile()
	return
}

// GetJenkinsWarURL returns a URL of Jenkins war file
func (u *UpdateCenterManager) GetJenkinsWarURL() (warURL string) {
	version := u.Version
	if version == "" {
		version = "latest"
	}

	if u.Formula != "" {
		warURL = fmt.Sprintf("https://dl.bintray.com/jenkins-zh/generic/jenkins/%s/jenkins-%s.war", version, u.Formula)
	} else if u.LTS {
		warURL = fmt.Sprintf("%s/war-stable/%s/jenkins.war", strings.TrimRight(u.MirrorSite, "/"), version)
	} else {
		warURL = fmt.Sprintf("%s/war/%s/jenkins.war", strings.TrimRight(u.MirrorSite, "/"), version)
	}
	return
}

// GetSite is get Available Plugins and Updated Plugins from UpdateCenter
func (u *UpdateCenterManager) GetSite() (site *CenterSite, err error) {
	err = u.RequestWithData(http.MethodGet, "/updateCenter/site/default/api/json?pretty=true&depth=2", nil, nil, 200, &site)
	return
}

// ChangeUpdateCenterSite updates the update center address
func (u *UpdateCenterManager) ChangeUpdateCenterSite(name, updateCenterURL string) (err error) {
	formData := url.Values{}
	formData.Add("site", updateCenterURL)
	payload := strings.NewReader(formData.Encode())

	api := "/pluginManager/siteConfigure"
	_, err = u.RequestWithoutData(http.MethodPost, api,
		map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 200)
	return
}

// SetMirrorCertificate take the mirror certificate file or not
func (u *UpdateCenterManager) SetMirrorCertificate(enable bool) (err error) {
	api := "/update-center-mirror/use"
	if !enable {
		api = "/update-center-mirror/remove"
	}

	_, err = u.RequestWithoutData(http.MethodPost, api,
		map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, nil, 200)
	return
}

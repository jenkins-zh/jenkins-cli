package client

import (
	"net/http"
)

// AgentLabel represents the label of Jenkins agent
type AgentLabel struct {
	Name string
}

// View represents the view of Jenkins
type View struct {
	Name string
	URL  string
}

// JenkinsStatus holds the status of Jenkins
type JenkinsStatus struct {
	AssignedLabels  []AgentLabel
	Description     string
	Jobs            []Job
	Mode            string
	NodeDescription string
	NodeName        string
	NumExecutors    int
	PrimaryView     View
	QuietingDown    bool
	SlaveAgentPort  int
	UseCrumbs       bool
	UseSecurity     bool
	Views           []View
	Version         string
}

// JenkinsStatusClient use to connect with Jenkins status
type JenkinsStatusClient struct {
	JenkinsCore
}

// Get returns status of Jenkins
func (q *JenkinsStatusClient) Get() (status *JenkinsStatus, err error) {
	status = &JenkinsStatus{}
	var response *http.Response
	response, err = q.RequestWithResponseHeader(http.MethodGet, "/api/json", nil, nil, status)
	if err == nil {
		if ver, ok := response.Header["X-Jenkins"]; ok && len(ver) > 0 {
			status.Version = ver[0]
		}
	}
	return
}

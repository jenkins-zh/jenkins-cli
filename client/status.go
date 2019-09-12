package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	api := fmt.Sprintf("%s/api/json", q.URL)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("GET", api, nil)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}

	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			status = &JenkinsStatus{}
			if ver, ok := response.Header["X-Jenkins"]; ok && len(ver) > 0 {
				status.Version = ver[0]
			}
			err = json.Unmarshal(data, status)
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

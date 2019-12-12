package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ComputerClient is client for operate computers
type ComputerClient struct {
	JenkinsCore
}

// List get the computer list
func (c *ComputerClient) List() (computers ComputerList, err error) {
	err = c.RequestWithData("GET", "/computer/api/json",
		nil, nil, 200, &computers)
	return
}

// Launch starts up a agent
func (c *ComputerClient) Launch(name string) (err error) {
	api := fmt.Sprintf("/computer/%s/launchSlaveAgent", name)
	_, err = c.RequestWithoutData("POST", api, nil, nil, 200)
	return
}

// GetLog fetch the log a computer
func (c *ComputerClient) GetLog(name string) (log string, err error) {
	var response *http.Response
	api := fmt.Sprintf("/computer/%s/logText/progressiveText", name)
	if response, err = c.RequestWithResponse("GET", api, nil, nil); err == nil {
		statusCode := response.StatusCode
		if statusCode != 200 {
			err = fmt.Errorf("unexpected status code %d", statusCode)
			return
		}

		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			log = string(data)
		}
	}
	return
}

// Computer is the agent of Jenkins
type Computer struct {
	AssignedLabels      []ComputerLabel
	Description         string
	DisplayName         string
	Idle                bool
	JnlpAgent           bool
	LaunchSupported     bool
	ManualLaunchAllowed bool
	NumExecutors        int
	Offline             bool
	OfflineCause        OfflineCause
	OfflineCauseReason  string
	TemporarilyOffline  bool
}

// OfflineCause is the cause of computer offline
type OfflineCause struct {
	Timestamp   int64
	Description string
}

// ComputerList represents the list of computer from API
type ComputerList struct {
	busyExecutors  int
	Computer       []Computer
	TotalExecutors int
}

// ComputerLabel represents the label of a computer
type ComputerLabel struct {
	Name string
}

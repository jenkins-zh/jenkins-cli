package client

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
	OfflineCause        string
	OfflineCauseReason  string
	TemporarilyOffline  bool
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

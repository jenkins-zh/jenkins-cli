package client

import (
	"fmt"

	"github.com/jenkins-zh/jenkins-cli/util"
)

// QueueClient is the client of queue
type QueueClient struct {
	JenkinsCore
}

// Get returns the job queue
func (q *QueueClient) Get() (status *JobQueue, err error) {
	err = q.RequestWithData("GET", "/queue/api/json", nil, nil, 200, &status)
	return
}

// Cancel will cancel a job from the queue
func (q *QueueClient) Cancel(id int) (err error) {
	api := fmt.Sprintf("/queue/cancelItem?id=%d", id)
	header := make(map[string]string)
	header["Content-Type"] = util.APP_FORM
	var statusCode int
	if statusCode, err = q.RequestWithoutData("POST", api, header, nil, 302); err != nil && statusCode == 200 {
		err = nil
	}
	return
}

// JobQueue represent the job queue
type JobQueue struct {
	Items []QueueItem
}

// QueueItem is the item of job queue
type QueueItem struct {
	Blocked                    bool
	Buildable                  bool
	ID                         int
	Params                     string
	Pending                    bool
	Stuck                      bool
	URL                        string
	Why                        string
	BuildableStartMilliseconds int64
	InQueueSince               int64
	Actions                    []CauseAction
}

// CauseAction is the collection of causes
type CauseAction struct {
	Causes []Cause
}

// Cause represent the reason why job is triggered
type Cause struct {
	UpstreamURL      string
	UpstreamProject  string
	UpstreamBuild    int
	ShortDescription string
}

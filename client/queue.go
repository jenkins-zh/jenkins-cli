package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type QueueClient struct {
	JenkinsCore
}

func (q *QueueClient) Get() (status *JobQueue, err error) {
	api := fmt.Sprintf("%s/queue/api/json", q.URL)
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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			if err == nil {
				status = &JobQueue{}
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

type JobQueue struct {
	Items []QueueItem
}

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

type CauseAction struct {
	Causes []Cause
}

type Cause struct {
	UpstreamUrl      string
	UpstreamProject  string
	UpstreamBuild    int
	ShortDescription string
}

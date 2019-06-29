package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type JobClient struct {
	JenkinsCore
}

// Search find a set of jobs by name
func (q *JobClient) Search(keyword string) (status *SearchResult, err error) {
	api := fmt.Sprintf("%s/search/suggest?query=%s", q.URL, keyword)
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
			if err == nil {
				status = &SearchResult{}
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

func (q *JobClient) Build(jobName string) (err error) {
	jobItems := strings.Split(jobName, " ")
	path := ""
	for _, item := range jobItems {
		path = fmt.Sprintf("%s/job/%s", path, item)
	}

	api := fmt.Sprintf("%s/%s/build", q.URL, path)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("POST", api, nil)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}

	if err = q.CrumbHandle(req); err != nil {
		log.Fatal(err)
	}

	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			fmt.Println("build successfully")
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

// Log get the log of a job
func (q *JobClient) Log(jobName string, start int64) (jobLog JobLog, err error) {
	jobItems := strings.Split(jobName, " ")
	path := ""
	for _, item := range jobItems {
		path = fmt.Sprintf("%s/job/%s", path, item)
	}

	api := fmt.Sprintf("%s/%s/lastBuild/logText/progressiveText?start=%d", q.URL, path, start)
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
	jobLog = JobLog{
		HasMore:   false,
		Text:      "",
		NextStart: int64(0),
	}
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			jobLog.Text = string(data)

			if response.Header != nil {
				jobLog.HasMore = strings.ToLower(response.Header.Get("X-More-Data")) == "true"
				jobLog.NextStart, _ = strconv.ParseInt(response.Header.Get("X-Text-Size"), 10, 64)
			}
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

type JobLog struct {
	HasMore   bool
	NextStart int64
	Text      string
}

type SearchResult struct {
	Suggestions []SearchResultItem
}

type SearchResultItem struct {
	Name string
}

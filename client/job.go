package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
		if code == 201 { // Jenkins will send redirect by this api
			fmt.Println("build successfully")
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *JobClient) GetBuild(jobName string, id int) (job *JobBuild, err error) {
	jobItems := strings.Split(jobName, " ")
	path := ""
	for _, item := range jobItems {
		path = fmt.Sprintf("%s/job/%s", path, item)
	}

	var api string
	if id == -1 {
		api = fmt.Sprintf("%s/%s/lastBuild/api/json", q.URL, path)
	} else {
		api = fmt.Sprintf("%s/%s/%d/api/json", q.URL, path, id)
	}
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
			job = &JobBuild{}
			err = json.Unmarshal(data, job)
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *JobClient) BuildWithParams(jobName string, parameters []ParameterDefinition) (err error) {
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

	var paramJSON []byte

	if len(parameters) == 1 {
		paramJSON, err = json.Marshal(parameters[0])
	} else {
		paramJSON, err = json.Marshal(parameters)
	}

	formData := url.Values{"json": {fmt.Sprintf("{\"parameter\": %s}", string(paramJSON))}}
	payload := strings.NewReader(formData.Encode())
	req, err = http.NewRequest("POST", api, payload)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}

	if err = q.CrumbHandle(req); err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 201 { // Jenkins will send redirect by this api
			fmt.Println("build successfully")
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *JobClient) GetJob(name string) (job *Job, err error) {
	jobItems := strings.Split(name, " ")
	path := ""
	for _, item := range jobItems {
		path = fmt.Sprintf("%s/job/%s", path, item)
	}

	api := fmt.Sprintf("%s/%s/api/json", q.URL, path)
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
			job = &Job{}
			err = json.Unmarshal(data, job)
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *JobClient) GetJobTypeCategories() (jobCategories []JobCategory, err error) {
	api := fmt.Sprintf("%s/view/all/itemCategories?depth=3", q.URL)
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
			type innerJobCategories struct {
				Categories []JobCategory
			}
			result := &innerJobCategories{}
			err = json.Unmarshal(data, result)
			jobCategories = result.Categories
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *JobClient) UpdatePipeline(name, script string) (err error) {
	jobItems := strings.Split(name, " ")
	path := ""
	for i, item := range jobItems {
		if i == 0 {
			path = fmt.Sprintf("job/%s", item)
		} else {
			path = fmt.Sprintf("%s/job/%s", path, item)
		}
	}

	api := fmt.Sprintf("%s/%s/restFul/update", q.URL, path)
	var (
		req      *http.Request
		response *http.Response
	)

	formData := url.Values{"script": {script}}
	payload := strings.NewReader(formData.Encode())
	req, err = http.NewRequest("POST", api, payload)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}

	if err = q.CrumbHandle(req); err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			fmt.Println("updated")
		} else {
			fmt.Println("code", code)
			log.Fatal(string(data))
		}
	} else {
		fmt.Println("request is error")
		log.Fatal(err)
	}
	return
}

func (q *JobClient) GetPipeline(name string) (pipeline *Pipeline, err error) {
	jobItems := strings.Split(name, " ")
	path := ""
	for _, item := range jobItems {
		path = fmt.Sprintf("%s/job/%s", path, item)
	}

	api := fmt.Sprintf("%s/%s/restFul", q.URL, path)
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
			pipeline = &Pipeline{}
			err = json.Unmarshal(data, pipeline)
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *JobClient) GetHistory(name string) (builds []JobBuild, err error) {
	var job *Job
	if job, err = q.GetJob(name); err == nil {
		builds = job.Builds
	}
	return
}

// Log get the log of a job
func (q *JobClient) Log(jobName string, history int, start int64) (jobLog JobLog, err error) {
	jobItems := strings.Split(jobName, " ")
	path := ""
	for _, item := range jobItems {
		path = fmt.Sprintf("%s/job/%s", path, item)
	}

	var api string
	if history == -1 {
		api = fmt.Sprintf("%s/%s/lastBuild/logText/progressiveText?start=%d", q.URL, path, start)
	} else {
		api = fmt.Sprintf("%s/%s/%d/logText/progressiveText?start=%d", q.URL, path, history, start)
	}
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

func (q *JobClient) Create(jobName string, jobType string) (err error) {
	api := fmt.Sprintf("%s/view/all/createItem", q.URL)
	var (
		req      *http.Request
		response *http.Response
	)

	type playLoad struct {
		Name string `json:"name"`
		Mode string `json:"mode"`
		From string
	}

	playLoadObj := &playLoad{
		Name: jobName,
		Mode: jobType,
		From: "",
	}

	playLoadData, _ := json.Marshal(playLoadObj)

	formData := url.Values{
		"json": {string(playLoadData)},
		"name": {jobName},
		"mode": {jobType},
	}
	payload := strings.NewReader(formData.Encode())

	req, err = http.NewRequest("POST", api, payload)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 302 || code == 200 { // Jenkins will send redirect by this api
			fmt.Println("create successfully")
		} else {
			fmt.Printf("status code: %d\n", code)
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *JobClient) Delete(jobName string) (err error) {
	api := fmt.Sprintf("%s/job/%s/doDelete", q.URL, jobName)
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
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 302 || code == 200 { // Jenkins will send redirect by this api
			fmt.Println("delete successfully")
		} else {
			fmt.Printf("status code: %d\n", code)
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

type Job struct {
	Type            string `json:"_class"`
	Builds          []JobBuild
	Color           string
	ConcurrentBuild bool
	Name            string
	NextBuildNumber int
	URL             string
	Buildable       bool

	Property []ParametersDefinitionProperty
}

type ParametersDefinitionProperty struct {
	ParameterDefinitions []ParameterDefinition
}

type ParameterDefinition struct {
	Description           string
	Name                  string `json:"name"`
	Type                  string
	Value                 string `json:"value"`
	DefaultParameterValue DefaultParameterValue
}

type DefaultParameterValue struct {
	Description string
	Value       string
}

type SimpleJobBuild struct {
	Number int
	URL    string
}

type JobBuild struct {
	SimpleJobBuild
	Building          bool
	Description       string
	DisplayName       string
	Duration          int64
	EstimatedDuration int64
	FullDisplayName   string
	ID                string
	KeepLog           bool
	QueueID           int
	Result            string
	Timestamp         int64
	PreviousBuild     SimpleJobBuild
	NextBuild         SimpleJobBuild
}

type Pipeline struct {
	Script  string
	Sandbox bool
}

type JobCategory struct {
	Description string
	ID          string
	Items       []JobCategoryItem
	MinToShow   int
	Name        string
	Order       int
}

type JobCategoryItem struct {
	Description string
	DisplayName string
	Order       int
	Class       string
}

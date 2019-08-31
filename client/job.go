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

	"github.com/jenkins-zh/jenkins-cli/util"
)

// JobClient is client for operate jobs
type JobClient struct {
	JenkinsCore
}

// Search find a set of jobs by name
func (q *JobClient) Search(keyword string) (status *SearchResult, err error) {
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = q.Request("GET", fmt.Sprintf("/search/suggest?query=%s", keyword), nil, nil); err == nil {
		if statusCode == 200 {
			json.Unmarshal(data, &status)
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
		}
	}
	return
}

// Build trigger a job
func (q *JobClient) Build(jobName string) (err error) {
	path := parseJobPath(jobName)

	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = q.Request("POST", fmt.Sprintf("%s/build", path), nil, nil); err == nil {
		if statusCode == 201 {
			fmt.Println("build successfully")
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
			if q.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
	}
	return
}

// GetBuild get build information of a job
func (q *JobClient) GetBuild(jobName string, id int) (job *JobBuild, err error) {
	path := parseJobPath(jobName)
	var api string
	if id == -1 {
		api = fmt.Sprintf("%s/lastBuild/api/json", path)
	} else {
		api = fmt.Sprintf("%s/%d/api/json", path, id)
	}

	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = q.Request("GET", api, nil, nil); err == nil {
		if statusCode == 200 {
			job = &JobBuild{}
			err = json.Unmarshal(data, job)
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
			if q.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
	}
	return
}

// BuildWithParams build a job which has params
func (q *JobClient) BuildWithParams(jobName string, parameters []ParameterDefinition) (err error) {
	path := parseJobPath(jobName)
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
	req.Header.Add(util.CONTENT_TYPE, util.APP_FORM)
	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 201 { // Jenkins will send redirect by this api
			fmt.Println("build successfully")
		} else {
			fmt.Println("Status code", code)
			if q.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
	} else {
		log.Fatal(err)
	}
	return
}

// StopJob stops a job build
func (q *JobClient) StopJob(jobName string, num int) (err error) {
	path := parseJobPath(jobName)
	api := fmt.Sprintf("%s/%d/stop", path, num)
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = q.Request("POST", api, nil, nil); err == nil {
		if statusCode == 200 {
			fmt.Println("stoped successfully")
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
			if q.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
	}
	return
}

// GetJob returns the job info
func (q *JobClient) GetJob(name string) (job *Job, err error) {
	path := parseJobPath(name)
	api := fmt.Sprintf("%s/api/json", path)
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = q.Request("GET", api, nil, nil); err == nil {
		if statusCode == 200 {
			job = &Job{}
			err = json.Unmarshal(data, job)
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
			if q.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
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
	path := parseJobPath(name)
	api := fmt.Sprintf("%s/%s/wfapisu/update", q.URL, path)
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
	req.Header.Add(util.CONTENT_TYPE, util.APP_FORM)
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
	path := parseJobPath(name)
	api := fmt.Sprintf("%s/%s/wfapisu/script", q.URL, path)
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

		for i, build := range builds {
			api := fmt.Sprintf("%sapi/json", build.URL)
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
					err = json.Unmarshal(data, &build)
					builds[i] = build
				} else {
					log.Fatal(string(data))
				}
			} else {
				log.Fatal(err)
			}
		}
	}
	return
}

// Log get the log of a job
func (q *JobClient) Log(jobName string, history int, start int64) (jobLog JobLog, err error) {
	path := parseJobPath(jobName)
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
	req.Header.Add(util.CONTENT_TYPE, util.APP_FORM)

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
	req.Header.Add(util.CONTENT_TYPE, util.APP_FORM)

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

// parseJobPath leads with slash
func parseJobPath(jobName string) (path string) {
	jobItems := strings.Split(jobName, " ")
	path = ""
	for _, item := range jobItems {
		path = fmt.Sprintf("%s/job/%s", path, item)
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
	Value       interface{}
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

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
func (q *JobClient) Search(keyword string, max int) (status *SearchResult, err error) {
	err = q.RequestWithData("GET", fmt.Sprintf("/search/suggest?query=%s&max=%d", keyword, max), nil, nil, 200, &status)
	return
}

// Build trigger a job
func (q *JobClient) Build(jobName string) (err error) {
	path := parseJobPath(jobName)
	_, err = q.RequestWithoutData("POST", fmt.Sprintf("%s/build", path), nil, nil, 201)
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

	err = q.RequestWithData("GET", api, nil, nil, 200, &job)
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
	req.Header.Add(util.ContentType, util.ApplicationForm)
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

	_, err = q.RequestWithoutData("POST", api, nil, nil, 200)
	return
}

// GetJob returns the job info
func (q *JobClient) GetJob(name string) (job *Job, err error) {
	path := parseJobPath(name)
	api := fmt.Sprintf("%s/api/json", path)

	err = q.RequestWithData("GET", api, nil, nil, 200, &job)
	return
}

// GetJobTypeCategories returns all categories of jobs
func (q *JobClient) GetJobTypeCategories() (jobCategories []JobCategory, err error) {
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = q.Request("GET", "/view/all/itemCategories?depth=3", nil, nil); err == nil {
		if statusCode == 200 {
			type innerJobCategories struct {
				Categories []JobCategory
			}
			result := &innerJobCategories{}
			err = json.Unmarshal(data, result)
			jobCategories = result.Categories
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
			if q.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
	}
	return
}

// GetPipeline return the pipeline object
func (q *JobClient) GetPipeline(name string) (pipeline *Pipeline, err error) {
	path := parseJobPath(name)
	api := fmt.Sprintf("%s/restFul", path)
	err = q.RequestWithData("GET", api, nil, nil, 200, &pipeline)
	return
}

// UpdatePipeline updates the pipeline script
func (q *JobClient) UpdatePipeline(name, script string) (err error) {
	path := parseJobPath(name)
	api := fmt.Sprintf("%s/restFul/update", path)

	formData := url.Values{"script": {script}}
	payload := strings.NewReader(formData.Encode())
	_, err = q.RequestWithoutData("POST", api, map[string]string{util.ContentType: util.ApplicationForm}, payload, 200)
	return
}

// GetHistory returns the build history of a job
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

// Create can create a job
func (q *JobClient) Create(jobName string, jobType string) (err error) {
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

	var code int
	code, err = q.RequestWithoutData("POST", "/view/all/createItem", map[string]string{util.ContentType: util.ApplicationForm}, payload, 200)
	if code == 302 {
		err = nil
	}
	return
}

// Delete will delete a job by name
func (q *JobClient) Delete(jobName string) (err error) {
	var (
		statusCode int
		data       []byte
	)

	api := fmt.Sprintf("/job/%s/doDelete", jobName)
	header := map[string]string{
		util.ContentType: util.ApplicationForm,
	}

	if statusCode, data, err = q.Request("POST", api, header, nil); err == nil {
		if statusCode == 200 || statusCode == 302 {
			fmt.Println("delete successfully")
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
			if q.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
	}
	return
}

// GetJobInputActions returns the all pending actions
func (q *JobClient) GetJobInputActions(jobName string, buildID int) (actions []JobInputItem, err error) {
	path := parseJobPath(jobName)
	err = q.RequestWithData("GET", fmt.Sprintf("%s/%d/wfapi/pendingInputActions", path, buildID), nil, nil, 200, &actions)
	return
}

// jenkinsInputParametersRequest represents the parameters for the Jenkins input request
type jenkinsInputParametersRequest struct {
	Parameter []jenkinsPipelineParameter `json:"parameter"`
}

type jenkinsPipelineParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (q *JobClient) JobInputSubmitTest(jobName, inputID string, buildID int, abort bool, params map[string]string) (err error) {
	jobPath := parseJobPath(jobName)
	// /job/pipeline/1/input/inputid/proceed
	// /job/pipeline/1/input/inputid/abort
	var api string
	if abort {
		api = fmt.Sprintf("%s/%d/input/%s/abort", jobPath, buildID, inputID)
	} else {
		api = fmt.Sprintf("%s/%d/input/%s/proceed", jobPath, buildID, inputID)
	}

	request := jenkinsInputParametersRequest{
		Parameter: make([]jenkinsPipelineParameter, 0),
	}

	for k, v := range params {
		request.Parameter = append(request.Parameter, jenkinsPipelineParameter{
			Name:  k,
			Value: v,
		})
	}

	var paramData []byte
	// var payload *strings.Reader
	// if paramData, err = json.Marshal(request); err == nil {
	// 	payload = strings.NewReader(url.Values{"json": {string(paramData)}}.Encode())
	// }
	paramData, _ = json.Marshal(request)

	api = fmt.Sprintf("%s?json=%s", api, string(paramData))

	// header := map[string]string{
	// 	"Content-Type": "application/x-www-form-urlencoded",
	// }

	// if len(params) == 0 {
	// 	_, err = q.RequestWithoutData(http.MethodPost, api, header, nil, 200)
	// 	payload = nil
	// } else {
	// 	_, err = q.RequestWithoutData(http.MethodPost, api, header, payload, 200)
	// }
	// fmt.Println(payload)
	fmt.Println(api)

	_, err = q.RequestWithoutData("POST", api, nil, nil, 200)

	return
}

// JobInputSubmit submit the params
func (q *JobClient) JobInputSubmit(submitURL string, params map[string]string) (err error) {
	paramDefs := []ParameterDefinition{}

	for k, v := range params {
		paramDefs = append(paramDefs, ParameterDefinition{
			Name:  k,
			Value: v,
		})
	}

	paramsMap := map[string][]ParameterDefinition{}
	paramsMap["parameter"] = paramDefs

	json, _ := json.Marshal(paramsMap)
	formData := url.Values{
		"json": []string{string(json)},
	}
	payload := strings.NewReader(formData.Encode())
	fmt.Println(formData.Encode())
	fmt.Println(submitURL)

	_, err = q.RequestWithoutData("POST", submitURL, nil, payload, 200)
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

// JobLog holds the log text
type JobLog struct {
	HasMore   bool
	NextStart int64
	Text      string
}

// SearchResult holds the result items
type SearchResult struct {
	Suggestions []SearchResultItem
}

// SearchResultItem hold the result item
type SearchResultItem struct {
	Name string
}

// Job represents a job
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

// ParametersDefinitionProperty holds the param definition property
type ParametersDefinitionProperty struct {
	ParameterDefinitions []ParameterDefinition
}

// ParameterDefinition holds the parameter definition
type ParameterDefinition struct {
	Description           string
	Name                  string `json:"name"`
	Type                  string
	Value                 string `json:"value"`
	DefaultParameterValue DefaultParameterValue
}

// DefaultParameterValue represents the default value for param
type DefaultParameterValue struct {
	Description string
	Value       interface{}
}

// SimpleJobBuild represents a simple job build
type SimpleJobBuild struct {
	Number int
	URL    string
}

// JobBuild represents a job build
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

// Pipeline represents a pipeline
type Pipeline struct {
	Script  string
	Sandbox bool
}

// JobCategory represents a job category
type JobCategory struct {
	Description string
	ID          string
	Items       []JobCategoryItem
	MinToShow   int
	Name        string
	Order       int
}

// JobCategoryItem represents a job category item
type JobCategoryItem struct {
	Description string
	DisplayName string
	Order       int
	Class       string
}

// JobInputItem represents a job input action
type JobInputItem struct {
	ID                  string
	AbortURL            string
	Message             string
	ProceedText         string
	ProceedURL          string
	RedirectApprovalURL string
}

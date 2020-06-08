package client

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"moul.io/http2curl"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/util"
)

// JobClient is client for operate jobs
type JobClient struct {
	JenkinsCore

	Parent string
}

// Search find a set of jobs by name
func (q *JobClient) Search(name, kind string, start, limit int) (items []JenkinsItem, err error) {
	err = q.RequestWithData("GET", fmt.Sprintf("/items/list?name=%s&type=%s&start=%d&limit=%d&parent=%s",
		name, kind, start, limit, q.Parent),
		nil, nil, 200, &items)
	return
}

// Build trigger a job
func (q *JobClient) Build(jobName string) (err error) {
	path := ParseJobPath(jobName)
	_, err = q.RequestWithoutData("POST", fmt.Sprintf("%s/build", path), nil, nil, 201)
	return
}

// GetBuild get build information of a job
func (q *JobClient) GetBuild(jobName string, id int) (job *JobBuild, err error) {
	path := ParseJobPath(jobName)
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
	path := ParseJobPath(jobName)
	api := fmt.Sprintf("%s/build", path)

	var paramJSON []byte
	if len(parameters) == 1 {
		paramJSON, err = json.Marshal(parameters[0])
	} else {
		paramJSON, err = json.Marshal(parameters)
	}

	if err == nil {
		formData := url.Values{"json": {fmt.Sprintf("{\"parameter\": %s}", string(paramJSON))}}
		payload := strings.NewReader(formData.Encode())

		_, err = q.RequestWithoutData("POST", api,
			map[string]string{util.ContentType: util.ApplicationForm}, payload, 201)
	}
	return
}

// DisableJob disable a job
func (q *JobClient) DisableJob(jobName string) (err error) {
	path := ParseJobPath(jobName)
	api := fmt.Sprintf("%s/disable", path)

	_, err = q.RequestWithoutData("POST", api, nil, nil, 200)
	return
}

// EnableJob disable a job
func (q *JobClient) EnableJob(jobName string) (err error) {
	path := ParseJobPath(jobName)
	api := fmt.Sprintf("%s/enable", path)

	_, err = q.RequestWithoutData("POST", api, nil, nil, 200)
	return
}

// StopJob stops a job build
func (q *JobClient) StopJob(jobName string, num int) (err error) {
	path := ParseJobPath(jobName)

	var api string
	if num <= 0 {
		api = fmt.Sprintf("%s/lastBuild/stop", path)
	} else {
		api = fmt.Sprintf("%s/%d/stop", path, num)
	}

	_, err = q.RequestWithoutData("POST", api, nil, nil, 200)
	return
}

// GetJob returns the job info
func (q *JobClient) GetJob(name string) (job *Job, err error) {
	path := ParseJobPath(name)
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
		}
	}
	return
}

// GetPipeline return the pipeline object
func (q *JobClient) GetPipeline(name string) (pipeline *Pipeline, err error) {
	path := ParseJobPath(name)
	api := fmt.Sprintf("%s/restFul", path)
	err = q.RequestWithData("GET", api, nil, nil, 200, &pipeline)
	return
}

// UpdatePipeline updates the pipeline script
func (q *JobClient) UpdatePipeline(name, script string) (err error) {
	formData := url.Values{}
	formData.Add("script", script)

	path := ParseJobPath(name)
	api := fmt.Sprintf("%s/restFul/update?%s", path, formData.Encode())

	_, err = q.RequestWithoutData("POST", api, nil, nil, 200)
	return
}

// GetHistory returns the build history of a job
func (q *JobClient) GetHistory(name string) (builds []*JobBuild, err error) {
	var job *Job
	if job, err = q.GetJob(name); err == nil {
		buildList := job.Builds // only contains basic info

		var build *JobBuild
		for _, buildItem := range buildList {
			build, err = q.GetBuild(name, buildItem.Number)
			if err != nil {
				break
			}
			builds = append(builds, build)
		}
	}
	return
}

// Log get the log of a job
func (q *JobClient) Log(jobName string, history int, start int64) (jobLog JobLog, err error) {
	path := ParseJobPath(jobName)
	var api string
	if history == -1 {
		api = fmt.Sprintf("%s%s/lastBuild/logText/progressiveText?start=%d", q.URL, path, start)
	} else {
		api = fmt.Sprintf("%s%s/%d/logText/progressiveText?start=%d", q.URL, path, history, start)
	}
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("GET", api, nil)
	if err == nil {
		err = q.AuthHandle(req)
	}
	if err != nil {
		return
	}

	client := q.GetClient()
	jobLog = JobLog{
		HasMore:   false,
		Text:      "",
		NextStart: int64(0),
	}

	if curlCmd, curlErr := http2curl.GetCurlCommand(req); curlErr == nil {
		logger.Debug("HTTP request as curl", zap.String("cmd", curlCmd.String()))
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
		}
	}
	return
}

// CreateJobPayload the payload for creating a job
type CreateJobPayload struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
	From string `json:"from"`
}

// Create can create a job
func (q *JobClient) Create(jobPayload CreateJobPayload) (err error) {
	playLoadData, _ := json.Marshal(jobPayload)
	formData := url.Values{
		"json": {string(playLoadData)},
		"name": {jobPayload.Name},
		"mode": {jobPayload.Mode},
		"from": {jobPayload.From},
	}
	payload := strings.NewReader(formData.Encode())

	var code int
	code, err = q.RequestWithoutData("POST", "/view/all/createItem",
		map[string]string{util.ContentType: util.ApplicationForm}, payload, 200)
	if code == 302 {
		err = nil
	}
	return
}

// Delete will delete a job by name
func (q *JobClient) Delete(jobName string) (err error) {
	var (
		statusCode int
	)

	jobName = ParseJobPath(jobName)
	api := fmt.Sprintf("%s/doDelete", jobName)
	header := map[string]string{
		util.ContentType: util.ApplicationForm,
	}

	if statusCode, _, err = q.Request("POST", api, header, nil); err == nil {
		if statusCode != 200 && statusCode != 302 {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
		}
	}
	return
}

// GetJobInputActions returns the all pending actions
func (q *JobClient) GetJobInputActions(jobName string, buildID int) (actions []JobInputItem, err error) {
	path := ParseJobPath(jobName)
	err = q.RequestWithData("GET", fmt.Sprintf("%s/%d/wfapi/pendingInputActions", path, buildID), nil, nil, 200, &actions)
	return
}

// JenkinsInputParametersRequest represents the parameters for the Jenkins input request
type JenkinsInputParametersRequest struct {
	Parameter []ParameterDefinition `json:"parameter"`
}

// JobInputSubmit submit the pending input request
func (q *JobClient) JobInputSubmit(jobName, inputID string, buildID int, abort bool, params map[string]string) (err error) {
	jobPath := ParseJobPath(jobName)
	var api string
	if abort {
		api = fmt.Sprintf("%s/%d/input/%s/abort", jobPath, buildID, inputID)
	} else {
		api = fmt.Sprintf("%s/%d/input/%s/proceed", jobPath, buildID, inputID)
	}

	request := JenkinsInputParametersRequest{
		Parameter: make([]ParameterDefinition, 0),
	}

	for k, v := range params {
		request.Parameter = append(request.Parameter, ParameterDefinition{
			Name:  k,
			Value: v,
		})
	}

	paramData, _ := json.Marshal(request)

	api = fmt.Sprintf("%s?json=%s", api, string(paramData))
	_, err = q.RequestWithoutData("POST", api, nil, nil, 200)

	return
}

// ParseJobPath leads with slash
func ParseJobPath(jobName string) (path string) {
	path = jobName
	if jobName == "" || strings.HasPrefix(jobName, "/job/") ||
		strings.HasPrefix(jobName, "job/") {
		return
	}

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

// JenkinsItem represents the item of Jenkins
type JenkinsItem struct {
	Name        string
	DisplayName string
	URL         string
	Description string
	Type        string

	/** comes from Job */
	Buildable bool
	Building  bool
	InQueue   bool

	/** comes from ParameterizedJob */
	Parameterized bool
	Disabled      bool
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
	Inputs              []ParameterDefinition
}

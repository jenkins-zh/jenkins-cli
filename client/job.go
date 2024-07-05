package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	httpdownloader "github.com/linuxsuren/http-downloader/pkg"

	"go.uber.org/zap"
	"moul.io/http2curl"
)

const (
	// StringParameterDefinition is the definition for string parameter
	StringParameterDefinition = "StringParameterDefinition"
	// FileParameterDefinition is the definition for file parameter
	FileParameterDefinition = "FileParameterDefinition"
	// QueueWaitDefinition is the definition for file queue state wait define
	QueueWaitDefinition = "hudson.model.Queue$WaitingItem"
)

// JobClient is client for operate jobs
type JobClient struct {
	JenkinsCore

	Parent string
}

// Search find a set of jobs by name
func (q *JobClient) Search(name, kind string, start, limit int) (items []JenkinsItem, err error) {
	err = q.RequestWithData(http.MethodGet, fmt.Sprintf("/items/list?name=%s&type=%s&start=%d&limit=%d&parent=%s",
		name, kind, start, limit, q.Parent),
		nil, nil, 200, &items)
	return
}

// SearchViaBlue searches jobs via the BlueOcean API
func (q *JobClient) SearchViaBlue(name string, start, limit int) (items []JenkinsItem, err error) {
	api := fmt.Sprintf("/blue/rest/search/?q=pipeline:*%s*;type:pipeline;organization:jenkins;excludedFromFlattening=jenkins.branch.MultiBranchProject,com.cloudbees.hudson.plugins.folder.AbstractFolder&filter=no-folders&start=%d&limit=%d",
		name, start, limit)
	err = q.RequestWithData(http.MethodGet, api,
		nil, nil, 200, &items)
	return
}

// Build trigger a job
func (q *JobClient) Build(jobName string) (err error) {
	path := ParseJobPath(jobName)
	_, err = q.RequestWithoutData(http.MethodPost, fmt.Sprintf("%s/build", path), nil, nil, 201)
	return
}

// IdentityBuild is the build which carry the identity cause
type IdentityBuild struct {
	Build JobBuild
	Cause IdentityCause
}

// IdentityCause carray a identity cause
type IdentityCause struct {
	UUID             string `json:"uuid"`
	ShortDescription string `json:"shortDescription"`
	Message          string
}

// BuildAndReturn trigger a job then returns the build info
func (q *JobClient) BuildAndReturn(jobName, cause string, timeout, delay int) (build IdentityBuild, err error) {
	path := ParseJobPath(jobName)

	api := fmt.Sprintf("%s/restFul/build?1=1", path)
	if timeout >= 0 {
		api += fmt.Sprintf("&timeout=%d", timeout)
	}
	if delay >= 0 {
		api += fmt.Sprintf("&delay=%d", delay)
	}
	if cause != "" {
		api += fmt.Sprintf("&identifyCause=%s", cause)
	}

	err = q.RequestWithData(http.MethodPost, api, nil, nil, 200, &build)
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

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	hasFileParam := false
	stringParameters := make([]ParameterDefinition, 0, len(parameters))
	for _, parameter := range parameters {
		if parameter.Type == FileParameterDefinition {
			hasFileParam = true
			var file *os.File
			file, err = os.Open(parameter.Filepath)
			if err != nil {
				return err
			}
			defer file.Close()

			var fWriter io.Writer
			fWriter, err = writer.CreateFormFile(parameter.Filepath, filepath.Base(parameter.Filepath))
			if err != nil {
				return err
			}
			_, err = io.Copy(fWriter, file)
		} else {
			stringParameters = append(stringParameters, parameter)
		}
	}

	var paramJSON []byte
	if len(stringParameters) == 1 {
		paramJSON, err = json.Marshal(stringParameters[0])
	} else {
		paramJSON, err = json.Marshal(stringParameters)
	}
	if err != nil {
		return
	}

	if hasFileParam {
		if err = writer.WriteField("json", fmt.Sprintf("{\"parameter\": %s}", string(paramJSON))); err != nil {
			return
		}

		if err = writer.Close(); err != nil {
			return
		}

		_, err = q.RequestWithoutData(http.MethodPost, api,
			map[string]string{httpdownloader.ContentType: writer.FormDataContentType()}, body, 201)
	} else {
		formData := url.Values{"json": {fmt.Sprintf("{\"parameter\": %s}", string(paramJSON))}}
		payload := strings.NewReader(formData.Encode())

		_, err = q.RequestWithDataResponse(http.MethodPost, api,
			map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 201)
	}
	return
}

// DisableJob disable a job
func (q *JobClient) DisableJob(jobName string) (err error) {
	path := ParseJobPath(jobName)
	api := fmt.Sprintf("%s/disable", path)

	_, err = q.RequestWithoutData(http.MethodPost, api, nil, nil, 200)
	return
}

// EnableJob disable a job
func (q *JobClient) EnableJob(jobName string) (err error) {
	path := ParseJobPath(jobName)
	api := fmt.Sprintf("%s/enable", path)

	_, err = q.RequestWithoutData(http.MethodPost, api, nil, nil, 200)
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

	_, err = q.RequestWithoutData(http.MethodPost, api, nil, nil, 200)
	return
}

// GetJob returns the job info
func (q *JobClient) GetJob(name string) (job *Job, err error) {
	path := ParseJobPath(name)
	api := fmt.Sprintf("%s/api/json", path)

	err = q.RequestWithData(http.MethodGet, api, nil, nil, 200, &job)
	return
}

// AddParameters add parameters to a Pipeline
func (q *JobClient) AddParameters(name, parameters string) (err error) {
	path := ParseJobPath(name)
	api := fmt.Sprintf("%s/restFul/addParameter", path)

	formData := url.Values{
		"params": {parameters},
	}
	payload := strings.NewReader(formData.Encode())
	_, err = q.RequestWithoutData(http.MethodPost, api, map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 200)
	return
}

// RemoveParameters add parameters to a Pipeline
func (q *JobClient) RemoveParameters(name, parameters string) (err error) {
	path := ParseJobPath(name)
	api := fmt.Sprintf("%s/restFul/removeParameter?params=%s", path, parameters)

	_, err = q.RequestWithoutData(http.MethodPost, api, nil, nil, 200)
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

	_, err = q.RequestWithoutData(http.MethodPost, api, nil, nil, 200)
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

// DeleteHistory returns the build history of a job
func (q *JobClient) DeleteHistory(jobName string, num int) (err error) {
	path := ParseJobPath(jobName)
	api := fmt.Sprintf("%s/%d/doDelete", path, num)
	_, err = q.RequestWithoutData(http.MethodPost, api, nil, nil, 200)
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
	return q.CreateJobInFolder(jobPayload, "")
}

// CreateJobInFolder creates a job in a specific folder and create folder first if the folder does not exist
func (q *JobClient) CreateJobInFolder(jobPayload CreateJobPayload, path string) (err error) {
	// create a job in path
	playLoadData, _ := json.Marshal(jobPayload)
	formData := url.Values{
		"json": {string(playLoadData)},
		"name": {jobPayload.Name},
		"mode": {jobPayload.Mode},
		"from": {jobPayload.From},
	}
	payload := strings.NewReader(formData.Encode())
	path = ParseJobPath(path)
	api := fmt.Sprintf("/view/all%s/createItem", path)
	var code int
	code, err = q.RequestWithoutData(http.MethodPost, api,
		map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 200)
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
		httpdownloader.ContentType: httpdownloader.ApplicationForm,
	}

	if statusCode, _, err = q.Request(http.MethodPost, api, header, nil); err == nil {
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
	_, err = q.RequestWithoutData(http.MethodPost, api, nil, nil, 200)

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

// BuildWithParamsGetResponse get params request response with run id...
func (q *JobClient) BuildWithParamsGetResponse(jobName string, parameters []ParameterDefinition, options JobCmdOptionsCommon) (resp JenkinsBuildState, err error) {
	path := ParseJobPath(jobName)
	api := fmt.Sprintf("%s/buildWithParameters?", path)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	hasFileParam := false
	stringParameters := make([]ParameterDefinition, 0, len(parameters))
	formData := url.Values{}
	for _, parameter := range parameters {
		if parameter.Type == FileParameterDefinition {
			hasFileParam = true
			var file *os.File
			file, err = os.Open(parameter.Filepath)
			if err != nil {
				return
			}
			defer file.Close()

			var fWriter io.Writer
			fWriter, err = writer.CreateFormFile(parameter.Filepath, filepath.Base(parameter.Filepath))
			if err != nil {
				return
			}
			_, err = io.Copy(fWriter, file)
		} else {
			stringParameters = append(stringParameters, parameter)
			formData.Set(parameter.Name, parameter.Value)
		}
	}

	var paramJSON []byte
	if len(stringParameters) == 1 {
		paramJSON, err = json.Marshal(stringParameters[0])
	} else {
		paramJSON, err = json.Marshal(stringParameters)
	}
	if err != nil {
		return
	}

	if hasFileParam {
		if err = writer.WriteField("json", fmt.Sprintf("{\"parameter\": %s}", string(paramJSON))); err != nil {
			return
		}

		if err = writer.Close(); err != nil {
			return
		}

		_, err = q.RequestWithoutData(http.MethodPost, api,
			map[string]string{httpdownloader.ContentType: writer.FormDataContentType()}, body, 201)
	} else {
		payload := strings.NewReader(formData.Encode())
		var jobRespState JenkinsBuildState
		var queueResp JenkinsBuildExecutable

		jobRespState, err = q.RequestWithDataResponse(http.MethodPost, api,
			map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 201)
		logger.Info("Build job trigger response msg...",
			zap.Int("statusCode", jobRespState.StatusCode),
			zap.Int64("queueId", jobRespState.QueueId),
			zap.Bool("isWaitForRunID", options.Wait))

		// if wait will query runId
		if options.Wait {
			if jobRespState.QueueId > 0 {
				if queueResp, err = q.GetBuildQueueIdResponseWait(jobRespState.QueueId, options.WaitInterval); err == nil {
					jobRespState.RunId = queueResp.Executable.Number
					jobRespState.BuildUrl = queueResp.Executable.URL
				}
				resp = jobRespState
				logger.Info("Build job state msg",
					zap.Int64("runId", jobRespState.RunId),
					zap.String("buildUrl", jobRespState.BuildUrl),
					zap.Int("statusCode", jobRespState.StatusCode),
				)
			}
		}
	}
	return
}

// GetBuildQueueIdResponse get queue api by queue id
func (q *JobClient) GetBuildQueueIdResponseWait(queueId int64, interval int) (resp JenkinsBuildExecutable, err error) {
	if queueId > 0 {
		for {
			time.Sleep(time.Duration(interval) * time.Second)
			logger.Info("Waiting seconds for query run id by queue id...", zap.Int("interval", interval), zap.Int64("queueId", queueId))
			if err = q.RequestWithData(http.MethodGet, GetQueueApi(queueId), nil, nil, 200, &resp); err != nil {
				logger.Error("failed to get queue item", zap.Error(err))
				return
			}
			if !resp.isWaitItem() {
				break
			}
		}
	}
	return
}

// GetQueue api uri
func GetQueueApi(queueId int64) string {
	return fmt.Sprintf("queue/item/%d/api/json", queueId)
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

	/** comes from blueOcean */
	FullName     string
	WeatherScore int
	Parameters   []ParameterDefinition
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
	Filepath              string `json:"file"`
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

// Jenkins build state for job
type JenkinsBuildState struct {
	StatusCode int    `json:"-"`
	RunId      int64  `json:"run_id,omitempty"`
	BuildUrl   string `json:"build_url,omitempty"`
	QueueId    int64  `json:"queue_id,omitempty"`
	BodyData   []byte `json:"-"`
}

// Jenkins build executable for response
type JenkinsBuildExecutable struct {
	Class      string                       `json:"_class,omitempty"`
	Executable JenkinsBuildExecutableInline `json:"executable,omitempty"`
}

// Jenkins build executable inline for response
type JenkinsBuildExecutableInline struct {
	Class  string `json:"_class,omitempty"`
	Number int64  `json:"number,omitempty"`
	URL    string `json:"url,omitempty"`
}

func (j *JenkinsBuildExecutable) isWaitItem() bool {
	return j.Class == QueueWaitDefinition
}

type JobCmdOptionsCommon struct {
	Wait         bool
	WaitTime     int
	LogConsole   bool
	WaitInterval int
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

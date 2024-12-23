package client

import (
	"fmt"
	"net/http"
	"strings"
)

// Artifact represents the artifacts from Jenkins build
type Artifact struct {
	ID   string
	Name string
	Path string
	URL  string
	Size int64
}

// JobWithArtifacts is the artifacts from a job
type JobWithArtifacts struct {
	Artifacts []JobArtifact `json:"artifacts"`
}

// GetArtifacts gets the artifacts from the JobWithArtifacts object
func (j JobWithArtifacts) GetArtifacts() (artifacts []Artifact) {
	for _, a := range j.Artifacts {
		artifacts = append(artifacts, Artifact{
			ID:   a.FileName,
			Name: a.FileName,
			Path: a.RelativePath,
		})
	}
	return
}

// JobArtifact represents the artifact object
type JobArtifact struct {
	RelativePath string `json:"relativePath"`
	FileName     string `json:"fileName"`
}

// ArtifactClient is client for getting the artifacts
type ArtifactClient struct {
	JenkinsCore
}

// List get the list of artifacts from a build
func (q *ArtifactClient) List(jobName string, buildID int) (artifacts []Artifact, err error) {
	path := ParseJobPath(jobName)
	var api string
	var oldAPI string
	if buildID < 1 {
		api = fmt.Sprintf("%s/lastBuild/wfapi/artifacts", path)
		oldAPI = fmt.Sprintf("%s/lastBuild/api/json", path)
	} else {
		api = fmt.Sprintf("%s/%d/wfapi/artifacts", path, buildID)
		oldAPI = fmt.Sprintf("%s/%d/api/json", path, buildID)
	}
	err = q.RequestWithData(http.MethodGet, api, nil, nil, 200, &artifacts)
	if err != nil {
		job := JobWithArtifacts{}
		if err = q.RequestWithData(http.MethodGet, oldAPI, nil, nil, 200, &job); err == nil {
			artifacts = job.GetArtifacts()

			for i := 0; i < len(artifacts); i++ {
				if artifacts[i].URL == "" {
					artifacts[i].URL = strings.ReplaceAll(oldAPI, "api/json", "artifact/") + artifacts[i].Path
				}
			}
		}
	}
	return
}

package client

import (
	"fmt"
)

// Artifact represents the artifacts from Jenkins build
type Artifact struct {
	ID   string
	Name string
	Path string
	URL  string
	Size int64
}

// ArtifactClient is client for getting the artifacts
type ArtifactClient struct {
	JenkinsCore
}

// List get the list of artifacts from a build
func (q *ArtifactClient) List(jobName string, buildID int) (artifacts []Artifact, err error) {
	path := ParseJobPath(jobName)
	var api string
	if buildID < 1 {
		api = fmt.Sprintf("%s/lastBuild/wfapi/artifacts", path)
	} else {
		api = fmt.Sprintf("%s/%d/wfapi/artifacts", path, buildID)
	}
	err = q.RequestWithData("GET", api, nil, nil, 200, &artifacts)
	return
}

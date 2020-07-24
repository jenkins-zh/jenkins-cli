package jenkins_plugin

import "encoding/xml"

type POM struct {
	XMLName    xml.Name  `xml:"project"`
	POMParent POMParent `xml:"parent"`
}

// POMParent is the versioning of maven
type POMParent struct {
	XMLName xml.Name `xml:"parent"`
	GAV GAV   `xml:"release"`
}

type GAV struct {
	GroupID  string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version string   `xml:"version"`
}
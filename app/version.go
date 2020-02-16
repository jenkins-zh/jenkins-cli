package app

import (
	"fmt"
)

var (
	version string
	commit  string
)

// GetVersion returns the version
func GetVersion() string {
	return version
}

// SetVersion is only for the test purpose
func SetVersion(ver string) {
	version = ver
}

// GetCommit returns the commit id
func GetCommit() string {
	return commit
}

// GetCombinedVersion returns the version and commit id
func GetCombinedVersion() string {
	return fmt.Sprintf("jcli; %s; %s", GetVersion(), GetCommit())
}

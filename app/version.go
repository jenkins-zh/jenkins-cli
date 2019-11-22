package app

import (
	"encoding/base64"
	"fmt"
)

var (
	version   string
	commit    string
	changelog string
)

// GetVersion returns the version
func GetVersion() string {
	return version
}

// GetCommit returns the commit id
func GetCommit() string {
	return commit
}

// GetCombinedVersion returns the version and commit id
func GetCombinedVersion() string {
	return fmt.Sprintf("jcli; %s; %s", GetVersion(), GetCommit())
}

// GetChangeLog returns the change log of release
func GetChangeLog() (log string, err error) {
	data := make([]byte, len(changelog))
	var length int
	length, err = base64.StdEncoding.Decode(data, []byte(changelog))
	fmt.Println(length)
	log = string(data[:length])
	return
}

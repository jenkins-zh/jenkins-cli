package test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestCrumb(t *testing.T) {
	cmd := exec.Command("jcli", "crumb", "--url", GetJenkinsURL())
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "Jenkins-Crumb=")
}

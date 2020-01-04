package test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestCrumb(t *testing.T) {
	cmd := exec.Command("jcli", "crumb", "--url", GetJenkinsURL())
	data, err := cmd.CombinedOutput()
	assert.NotNil(t, err)
	assert.Contains(t, string(data), "Error: crumb is disabled")
}

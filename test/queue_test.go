package test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestListQueue(t *testing.T) {
	cmd := exec.Command("jcli", "queue", "list", "--url", GetJenkinsURL())
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	assert.Contains(t, string(data), "ID Why URL")
}

package test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestListQueue(t *testing.T) {
	cmd := exec.Command("jcli", "queue", "list", "--url", "http://localhost:8080")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	assert.Contains(t, string(data), "ID Why URL")
}

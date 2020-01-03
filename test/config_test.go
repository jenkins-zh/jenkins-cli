package test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestConfigList(t *testing.T) {
	cmd := exec.Command("jcli", "config", "list")
	_, err := cmd.CombinedOutput()
	assert.Nil(t, err)
}

func TestConfigGenerate(t *testing.T) {
	cmd := exec.Command("jcli", "config", "generate", "-i=false")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "jenkins_servers")
}

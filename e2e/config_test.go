package e2e

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestConfigList(t *testing.T) {
	cmd := exec.Command("jcli", "config", "list")
	_, err := cmd.CombinedOutput()
	assert.NotNil(t, err)
}

func TestConfigGenerate(t *testing.T) {
	cmd := exec.Command("jcli", "config", "generate", "-i=false")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "jenkins_servers")
}

func TestShowCurrentConfig(t *testing.T) {
	cmd := exec.Command("jcli", "config")
	data, err := cmd.CombinedOutput()
	assert.NotNil(t, err)
	assert.Contains(t, string(data), "Error: no config file found or no current setting")
}

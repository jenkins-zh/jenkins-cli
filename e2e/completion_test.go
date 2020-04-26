package e2e

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestBashCompletion(t *testing.T) {
	cmd := exec.Command("jcli", "completion")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "bash completion for jcli")

	// with options
	cmd = exec.Command("jcli", "completion", "--type", "bash")
	data, err = cmd.CombinedOutput()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "bash completion for jcli")
}

func TestZshCompletion(t *testing.T) {
	cmd := exec.Command("jcli", "completion", "--type", "zsh")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "#compdef _jcli jcli")
}

func TestPowerShellCompletion(t *testing.T) {
	cmd := exec.Command("jcli", "completion", "--type", "powerShell")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "using namespace System.Management.Automation")
}

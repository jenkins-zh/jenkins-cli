package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestShutdown(t *testing.T) {
	cmd := exec.Command("jcli", "shutdown", "--url", GetJenkinsURL())
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err, fmt.Sprintf("failed in shutdown Jenkins, output is %s", string(data)))
}

package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestListQueue(t *testing.T) {
	cmd := exec.Command("jcli", "queue", "list", "--url", GetJenkinsURL())
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err, fmt.Sprintf("failed in cmd queue list, output is %s", string(data)))

	assert.Contains(t, string(data), "ID Why URL")
}

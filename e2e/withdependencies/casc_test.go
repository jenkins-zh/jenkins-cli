package withdependencies

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestCascExport(t *testing.T) {
	cmd := exec.Command("jcli", "casc", "export", "--url", GetJenkinsURL())
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "adminAddress")
}

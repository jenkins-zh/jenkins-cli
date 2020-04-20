package withdependencies

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchJobs(t *testing.T) {
	cmd := exec.Command("jcli", "job", "search", "--url", GetJenkinsURL())
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	fmt.Println(string(data))
}

func TestCreateJob(t *testing.T) {
	cmd := exec.Command("jcli", "job", "create", "fake",
		"--type", "com.cloudbees.hudson.plugins.folder.Folder", "--url", GetJenkinsURL())
	_, err := cmd.CombinedOutput()
	assert.Nil(t, err)
}

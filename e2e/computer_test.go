package e2e

import (
	"fmt"
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListComputers(t *testing.T) {
	cmd := exec.Command("jcli", "computer", "list", "--url", GetJenkinsURL())
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)
	fmt.Println(string(data))

	cmd = exec.Command("jcli", "computer", "create", "go", "--url", GetJenkinsURL())
	data, err = cmd.CombinedOutput()
	assert.Nil(t, err)
	fmt.Println(string(data))

	// test agent commands with docker mode
	if containerIsReady() {
		cmd = exec.Command("jcli", "computer", "launch", "go", "--agent-type", "golang", "-m", "docker", "--url", GetJenkinsURL())
		RunAndWait(cmd, func(reader io.ReadCloser) {
			WaitAgentRunningUp(reader)
		})

		cmd = exec.Command("jcli", "computer", "delete", "go", "--url", GetJenkinsURL())
		data, err = cmd.CombinedOutput()
		assert.Nil(t, err)
		fmt.Println(string(data))
	}
}

func containerIsReady() bool {
	var err error
	if _, err = exec.LookPath("docker"); err == nil {
		// only run these tests when the docker exists
		_, err = exec.Command("docker", "ps").CombinedOutput()
	}
	return err == nil
}

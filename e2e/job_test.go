package e2e

import (
	"fmt"
	"math"
	"math/rand"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListJobType(t *testing.T) {
	cmd := exec.Command("jcli", "job", "type", "--url", GetJenkinsURL())
	fmt.Println(cmd.String())
	data, err := cmd.CombinedOutput()
	fmt.Println(string(data))
	assert.Nil(t, err)

	rand.Seed(math.MaxInt8)
	name := fmt.Sprintf("%d", rand.Int())
	cmd = exec.Command("jcli", "job", "create", name, "--type", "hudson.model.FreeStyleProject", "--url", GetJenkinsURL(), "--logger-level", "debug")
	data, err = cmd.CombinedOutput()
	fmt.Println(string(data))
	assert.Nil(t, err)

	cmd = exec.Command("jcli", "job", "build", name, "-b", "--url", GetJenkinsURL())
	data, err = cmd.CombinedOutput()
	fmt.Println(string(data))
	assert.Nil(t, err)
	time.Sleep(time.Second * 6)

	cmd = exec.Command("jcli", "job", "history", name, "-d", "1", "--url", GetJenkinsURL())
	data, err = cmd.CombinedOutput()
	fmt.Println(string(data))
	assert.Nil(t, err)
}

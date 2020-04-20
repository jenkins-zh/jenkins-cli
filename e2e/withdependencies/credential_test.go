package withdependencies

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCredentials(t *testing.T) {
	cmd := exec.Command("jcli", "credential", "list", "--url", GetJenkinsURL())
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	fmt.Println(string(data))
}

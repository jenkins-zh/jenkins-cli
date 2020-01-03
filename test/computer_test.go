package test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListComputers(t *testing.T) {
	cmd := exec.Command("jcli", "computer", "list", "--url", "http://localhost:8080")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	fmt.Println(string(data))
}

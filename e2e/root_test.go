package e2e_test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestRoot(t *testing.T) {
	cmd := exec.Command("jcli")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)
	assert.Contains(t, string(data), "Jenkins CLI (jcli) manage your Jenkins")
}

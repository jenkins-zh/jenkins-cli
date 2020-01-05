package test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestDoc(t *testing.T) {
	tempDir := os.TempDir()
	defer os.Remove(tempDir)

	cmd := exec.Command("jcli", "doc", tempDir)
	_, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	_, err = os.Stat(path.Join(tempDir, "jcli.md"))
	assert.True(t, os.IsExist(err))
}

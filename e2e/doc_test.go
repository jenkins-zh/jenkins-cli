package e2e

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoc(t *testing.T) {
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("jcli", "doc", tempDir)
	_, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	_, err = os.Stat(path.Join(tempDir, "jcli.md"))
	assert.Nil(t, err)
}

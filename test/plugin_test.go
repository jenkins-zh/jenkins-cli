package test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPlugins(t *testing.T) {
	cmd := exec.Command("jcli", "plugin", "list", "--url", "http://localhost:8080")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	fmt.Println(string(data))
}

func TestSearchPlugins(t *testing.T) {
	cmd := exec.Command("jcli", "plugin", "search", "localization-zh-cn", "--url", "http://localhost:8080")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	fmt.Println(string(data))
}

func TestCheckUpdateCenter(t *testing.T) {
	cmd := exec.Command("jcli", "plugin", "check", "--url", "http://localhost:8080")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	fmt.Println(string(data))
}

func TestInstallPlugin(t *testing.T) {
	cmd := exec.Command("jcli", "plugin", "install", "localization-zh-cn", "--url", "http://localhost:8080")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	fmt.Println(string(data))
}

func TestDownloadPlugin(t *testing.T) {
	tempDir := os.TempDir()
	defer os.Remove(tempDir)

	cmd := exec.Command("jcli", "plugin", "download", "localization-zh-cn",
		"--download-dir", tempDir, "--url", "http://localhost:8080")
	data, err := cmd.CombinedOutput()
	assert.Nil(t, err)

	fmt.Println(string(data))
}

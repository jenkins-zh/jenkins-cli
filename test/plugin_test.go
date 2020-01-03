package test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	version := os.Getenv("JENKINS_VERSION")
	if version == "" {
		return
	}
	cmd := exec.Command("jcli", "center", "start", "--random-web-dir", "--setup-wizard=false", "--version", version)
	cmdStderrPipe, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	go func(reader io.ReadCloser, cmd *exec.Cmd) {
		buf := make([]byte, 1024, 1024)
		for {
			if strNum, err := reader.Read(buf); err != nil || strings.Contains(string(buf[:strNum]), "Jenkins is fully up and running") {
				break
			} else {
				fmt.Print(string(buf[:strNum]))
			}
		}

		m.Run()

		cmd.Process.Kill()
	}(cmdStderrPipe, cmd)

	err = cmd.Wait()
}

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

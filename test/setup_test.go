package test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
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

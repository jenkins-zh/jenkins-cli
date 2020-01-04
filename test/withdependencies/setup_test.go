package withdependencies

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var jenkinsURL string

func GetJenkinsURL() string {
	return jenkinsURL
}
func TestMain(m *testing.M) {
	version := os.Getenv("JENKINS_VERSION")
	if version == "" {
		return
	}

	jenkinsURL = "http://localhost:9090"

	cmd := exec.Command("jcli", "center", "start", "--random-web-dir", "--setup-wizard=false",
		"--version", version, "--port", "9090")
	cmdStderrPipe, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	go func(reader io.ReadCloser, cmd *exec.Cmd) {
		buf := make([]byte, 1024, 1024)
		for {
			if strNum, err := reader.Read(buf); err != nil || strings.Contains(string(buf[:strNum]), "Jenkins is fully up and running") {
				fmt.Print(string(buf[:strNum]))
				break
			} else {
				fmt.Print(string(buf[:strNum]))
			}
		}

		err = exec.Command("jcli", "plugin", "install", "localization-zh-cn", "--url", GetJenkinsURL()).Run()
		if err != nil {
			panic(err)
		}
		fmt.Println("install localization-zh-cn done")

		err = exec.Command("jcli", "center", "watch", "--util-install-complete", "--url", GetJenkinsURL()).Run()
		if err != nil {
			panic(err)
		}

		err = exec.Command("jcli", "restart", "-b", "--url", GetJenkinsURL()).Run()
		if err != nil {
			panic(err)
		}
		for {
			if strNum, err := reader.Read(buf); err != nil || strings.Contains(string(buf[:strNum]), "Jenkins is fully up and running") {
				break
			} else {
				fmt.Print(string(buf[:strNum]))
			}
		}

		exec.Command("jcli", "center", "mirror", "--url", GetJenkinsURL()).Run()
		exec.Command("jcli", "plugin", "check", "--url", GetJenkinsURL()).Run()
		err = exec.Command("jcli", "plugin", "install", "pipeline-restful-api", "--url", GetJenkinsURL()).Run()
		if err != nil {
			panic(err)
		}
		fmt.Println("install pipeline-restful-api done")

		err = exec.Command("jcli", "center", "watch", "--util-install-complete", "--url", GetJenkinsURL()).Run()
		if err != nil {
			panic(err)
		}

		err = exec.Command("jcli", "restart", "-b", "--url", GetJenkinsURL()).Run()
		if err != nil {
			panic(err)
		}
		for {
			if strNum, err := reader.Read(buf); err != nil || strings.Contains(string(buf[:strNum]), "Jenkins is fully up and running") {
				fmt.Print(string(buf[:strNum]))
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

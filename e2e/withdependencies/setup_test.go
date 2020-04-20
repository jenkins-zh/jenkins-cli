package withdependencies

import (
	"github.com/jenkins-zh/jenkins-cli/e2e"
	"io"
	"os"
	"os/exec"
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
		e2e.WaitRunningUp(reader)

		e2e.InstallPlugin("localization-zh-cn", GetJenkinsURL(), true)

		e2e.RestartAndWait(GetJenkinsURL(), reader)

		e2e.ExecuteCmd("center", "mirror", "--url", GetJenkinsURL())
		e2e.ExecuteCmd("plugin", "check", "--url", GetJenkinsURL())
		e2e.InstallPlugin("configuration-as-code", GetJenkinsURL(), true)
		e2e.InstallPlugin("pipeline-restful-api", GetJenkinsURL(), true)

		e2e.RestartAndWait(GetJenkinsURL(), reader)

		m.Run()

		if err = cmd.Process.Kill(); err != nil {
			panic(err)
		}
	}(cmdStderrPipe, cmd)

	err = cmd.Wait()
}

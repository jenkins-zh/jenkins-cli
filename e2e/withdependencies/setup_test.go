package withdependencies

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/e2e"
	"github.com/phayes/freeport"
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
	var err error

	version := os.Getenv("JENKINS_VERSION")
	os.Setenv("PATH", ".:"+os.Getenv("PATH"))

	javaHome := os.Getenv("JCLI_JAVA_HOME")
	if javaHome != "" {
		os.Setenv("PATH", javaHome+"/bin:"+os.Getenv("PATH"))
	}
	if err = os.Setenv("JCLI_CONFIG_LOAD", "false"); err != nil {
		panic(err)
	}
	if version == "" {
		return
	}

	var port int
	if port, err = freeport.GetFreePort(); err != nil {
		fmt.Println("get free port error", err)
		panic(err)
	}
	jenkinsURL = fmt.Sprintf("http://%s:%d", e2e.GetLocalIP(), port)

	cmd := exec.Command("jcli", "center", "start", "--random-web-dir", "--setup-wizard=false", "--port", fmt.Sprintf("%d", port), "--version", version)
	fmt.Println(cmd.String())
	e2e.RunAndWait(cmd, func(reader io.ReadCloser) {
		e2e.WaitJenkinsRunningUp(reader)

		e2e.ExecuteCmd("plugin", "check", "--url", GetJenkinsURL())
		e2e.InstallPlugin("localization-zh-cn", GetJenkinsURL(), true)

		e2e.RestartAndWait(GetJenkinsURL(), reader)

		e2e.ExecuteCmd("center", "mirror", "--url", GetJenkinsURL())
		e2e.ExecuteCmd("plugin", "check", "--url", GetJenkinsURL())
		e2e.InstallPlugin("configuration-as-code", GetJenkinsURL(), true)
		e2e.InstallPlugin("pipeline-restful-api", GetJenkinsURL(), true)

		e2e.RestartAndWait(GetJenkinsURL(), reader)

		m.Run()
	})
}

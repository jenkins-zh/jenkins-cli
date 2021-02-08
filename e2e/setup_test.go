package e2e

import (
	"fmt"
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
	jenkinsURL = fmt.Sprintf("http://%s:%d", GetLocalIP(), port)

	cmd := exec.Command("jcli", "center", "start", "--random-web-dir", "--setup-wizard=false",
		"--port", fmt.Sprintf("%d", port), "--version", version, "--thread", "10", "--clean-home")
	fmt.Println(cmd.String())
	RunAndWait(cmd, func(reader io.ReadCloser) {
		WaitJenkinsRunningUp(reader)

		m.Run()
	})
}

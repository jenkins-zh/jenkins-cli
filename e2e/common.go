package e2e

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// ExecuteCmd execute a jcli command
func ExecuteCmd(args ...string) {
	if err := exec.Command("jcli", args...).Run(); err != nil {
		panic(err)
	}
}

// InstallPlugin install a plugin by jcli
func InstallPlugin(name, jenkins string, wait bool) {
	ExecuteCmd("plugin", "install", name, "--url", jenkins)
	fmt.Printf("install %s done\n", name)
	if wait {
		ExecuteCmd("center", "watch", "--util-install-complete", "--url", jenkins)
	}
}

// RestartAndWait restart Jenkins then wait it
func RestartAndWait(jenkins string, outputReader io.ReadCloser) {
	buf := make([]byte, 1024, 1024)
	// should assert the error of restart
	_ = exec.Command("jcli", "restart", "-b", "--url", jenkins).Run()
	for {
		if strNum, err := outputReader.Read(buf); err != nil || strings.Contains(string(buf[:strNum]), "Jenkins is fully up and running") {
			break
		} else {
			fmt.Print(string(buf[:strNum]))
		}
	}
}

// WaitRunningUp wait until Jenkins running up
func WaitRunningUp(outputReader io.ReadCloser) {
	buf := make([]byte, 1024, 1024)
	for {
		if strNum, err := outputReader.Read(buf); err != nil || strings.Contains(string(buf[:strNum]), "Jenkins is fully up and running") {
			fmt.Print(string(buf[:strNum]))
			break
		} else {
			fmt.Print(string(buf[:strNum]))
		}
	}
}

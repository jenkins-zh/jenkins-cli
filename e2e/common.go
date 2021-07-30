package e2e

import (
	"fmt"
	"io"
	"net"
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

// RunAndWait run command and wait the callback function
func RunAndWait(cmd *exec.Cmd, callback func(reader io.ReadCloser)) {
	var err error
	cmdStderrPipe, _ := cmd.StderrPipe()
	if err = cmd.Start(); err != nil {
		panic(err)
	}

	go func(reader io.ReadCloser, cmd *exec.Cmd) {
		if callback != nil {
			callback(reader)
		}

		if err = cmd.Process.Kill(); err != nil {
			panic(err)
		}
	}(cmdStderrPipe, cmd)

	err = cmd.Wait()
}

// WaitJenkinsRunningUp wait until Jenkins running up
func WaitJenkinsRunningUp(outputReader io.ReadCloser) {
	WaitUntilExpect(outputReader, "Jenkins is fully up and running")
}

// WaitAgentRunningUp wait until Jenkins agent running up
func WaitAgentRunningUp(outputReader io.ReadCloser) {
	WaitUntilExpect(outputReader, "INFO: Connected")
}

// WaitUntilExpect wait until find the expect string
func WaitUntilExpect(outputReader io.ReadCloser, expect string) {
	buf := make([]byte, 1024, 1024)
	for {
		if strNum, err := outputReader.Read(buf); err != nil || strings.Contains(string(buf[:strNum]), expect) {
			fmt.Print(string(buf[:strNum]))
			break
		} else {
			fmt.Print(string(buf[:strNum]))
		}
	}
}

// GetLocalIP returns the local ip address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, value := range addrs {
			if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}
	return "127.0.0.1"
}

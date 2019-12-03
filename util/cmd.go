package util

import (
	"os"
	"os/exec"
	"runtime"
)

// Open a URL in a browser
func Open(url string, cmdContext ExecContext) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "Open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-Open"
	}
	args = append(args, url)

	if cmdContext == nil {
		cmdContext = exec.Command
	}
	return cmdContext(cmd, args...).Start()
}

// ExecContext is the context of system command caller
type ExecContext = func(name string, arg ...string) *exec.Cmd

// FakeExecCommandSuccess is a function that initialises a new exec.Cmd, one which will
// simply call TestShellProcessSuccess rather than the command it is provided. It will
// also pass through the command and its arguments as an argument to TestShellProcessSuccess
func FakeExecCommandSuccess(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestShellProcessSuccess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_TEST_PROCESS=1"}
	return cmd
}

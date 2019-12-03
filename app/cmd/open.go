package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// OpenOption is the open cmd option
type OpenOption struct {
	InteractiveOption

	Name   string
	Config bool
}

var openOption OpenOption

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().StringVarP(&openOption.Name, "name", "n", "",
		i18n.T("Open a specific Jenkins by name"))
	openCmd.Flags().BoolVarP(&openOption.Config, "config", "c", false,
		i18n.T("Open the configuration page of Jenkins"))
	openOption.SetFlag(openCmd)
}

var openCmd = &cobra.Command{
	Use:     "open [config name]",
	Short:   i18n.T("Open your Jenkins with a browse"),
	Long:    i18n.T(`Open your Jenkins with a browse`),
	Example: `jcli open -n <config name>`,
	Run: func(_ *cobra.Command, args []string) {
		var jenkins *JenkinsServer

		var configName string
		if len(args) > 0 {
			configName = args[0]
		} else if openOption.Name != "" {
			configName = openOption.Name
		}

		if configName == "" && openOption.Interactive {
			jenkinsNames := getJenkinsNames()
			prompt := &survey.Select{
				Message: "Choose a Jenkins that you want to Open:",
				Options: jenkinsNames,
			}
			survey.AskOne(prompt, &(configName))
		}

		if configName != "" {
			jenkins = findJenkinsByName(configName)
		} else {
			jenkins = getCurrentJenkins()
		}

		if jenkins != nil && jenkins.URL != "" {
			url := jenkins.URL
			if openOption.Config {
				url = fmt.Sprintf("%s/configure", url)
			}
			Open(url, exec.Command)
		} else {
			log.Fatalf("No URL found with Jenkins %s", configName)
		}
	},
}

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
	return cmdContext(cmd, args...).Start()
}

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

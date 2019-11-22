package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"log"
	"os/exec"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
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
				Message: "Choose a Jenkins that you want to open:",
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
			open(url)
		} else {
			log.Fatalf("No URL found with Jenkins %s", configName)
		}
	},
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

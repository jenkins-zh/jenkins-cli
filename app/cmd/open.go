package cmd

import (
	"fmt"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// OpenOption is the open cmd option
type OpenOption struct {
	InteractiveOption

	Name   string
	Config bool

	ExecContext util.ExecContext
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
	Short:   i18n.T("Open your Jenkins with a browser"),
	Long:    i18n.T(`Open your Jenkins with a browser`),
	Example: `jcli open -n [config name]`,
	RunE: func(_ *cobra.Command, args []string) (err error) {
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
				Message: i18n.T("Choose a Jenkins which you want to open:"),
				Options: jenkinsNames,
			}
			if err = survey.AskOne(prompt, &(configName)); err != nil {
				return
			}
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
			err = util.Open(url, openOption.ExecContext)
		} else {
			err = fmt.Errorf("no URL found with Jenkins %s", configName)
		}
		return
	},
}

// Open a URL in a browser
func Open(url string, cmdContext util.ExecContext) error {
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

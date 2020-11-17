package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

// ConfigSelectOptions is the option for select a config
type ConfigSelectOptions struct {
	common.Option
}

func init() {
	configCmd.AddCommand(configSelectCmd)
	configSelectOptions.Stdio = common.GetSystemStdio()
}

var configSelectOptions ConfigSelectOptions

var configSelectCmd = &cobra.Command{
	Use:   "select",
	Short: i18n.T("Select one config as current Jenkins"),
	Long:  i18n.T("Select one config as current Jenkins"),
	RunE: func(_ *cobra.Command, args []string) (err error) {
		var jenkinsName string
		if len(args) > 0 {
			jenkinsName = args[0]
		} else {
			target := ""
			if currentJenkins := getCurrentJenkins(); currentJenkins != nil {
				target = currentJenkins.Name
			}

			jenkinsName, err = configSelectOptions.Select(getJenkinsNames(),
				"Choose a Jenkins as the current one:", target)
		}

		if err == nil {
			setCurrentJenkins(jenkinsName)
		}
		return
	},
}

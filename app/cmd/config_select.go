package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configSelectCmd)
}

var configSelectCmd = &cobra.Command{
	Use:   "select [<name>]",
	Short: i18n.T("Select one config as current Jenkins"),
	Long:  i18n.T("Select one config as current Jenkins"),
	Run: func(_ *cobra.Command, args []string) {
		if len(args) > 0 {
			jenkinsName := args[0]

			setCurrentJenkins(jenkinsName)
		} else {
			selectByManual()
		}
	},
}

func selectByManual() {
	target := ""
	if currentJenkins := getCurrentJenkins(); currentJenkins != nil {
		target = currentJenkins.Name
	}

	prompt := &survey.Select{
		Message: "Choose a Jenkins as the current one:",
		Options: getJenkinsNames(),
		Default: target,
	}
	if err := survey.AskOne(prompt, &target); err == nil && target != "" {
		setCurrentJenkins(target)
	}
}

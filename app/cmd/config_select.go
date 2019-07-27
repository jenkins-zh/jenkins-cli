package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configSelectCmd)
}

var configSelectCmd = &cobra.Command{
	Use:   "select [<name>]",
	Short: "Select one config as current Jenkins",
	Long:  `Select one config as current Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
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

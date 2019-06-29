package cmd

import (
	"github.com/AlecAivazis/survey"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configSelectCmd)
}

var configSelectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select one config as current Jenkins",
	Long:  `Select one config as current Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		target := ""
		if currentJenkins := getCurrentJenkins(); currentJenkins != nil {
			target = currentJenkins.Name
		}

		prompt := &survey.Select{
			Message: "Choose a Jenkins as the current one:",
			Options: getJenkinsNames(),
			Default: target,
		}
		survey.AskOne(prompt, &target)

		if target != "" {
			setCurrentJenkins(target)
		}
	},
}

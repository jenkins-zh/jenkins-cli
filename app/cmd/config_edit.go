package cmd

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"io/ioutil"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configEditCmd)
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: i18n.T("Edit a Jenkins config"),
	Long:  i18n.T(`Edit a Jenkins config`),
	Run: func(_ *cobra.Command, _ []string) {
		current := getCurrentJenkinsFromOptionsOrDie()
		configPath := configOptions.ConfigFileLocation

		var data []byte
		var err error
		if data, err = ioutil.ReadFile(configPath); err != nil {
			log.Fatal(err)
		}

		content := string(data)
		prompt := &survey.Editor{
			Message:       fmt.Sprintf("Edit config item %s", current.Name),
			FileName:      "*.yaml",
			Help:          fmt.Sprintf("Config file path: %s", configPath),
			Default:       content,
			HideDefault:   true,
			AppendDefault: true,
		}

		if err := survey.AskOne(prompt, &content); err == nil {
			if err = ioutil.WriteFile(configPath, []byte(content), 0644); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

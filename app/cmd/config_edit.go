package cmd

import (
	"fmt"
	"log"

	"io/ioutil"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configEditCmd)
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a Jenkins config",
	Long:  `Edit a Jenkins config`,
	Run: func(_ *cobra.Command, _ []string) {
		current := getCurrentJenkins()
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

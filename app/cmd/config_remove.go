package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configRemoveCmd)
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a Jenkins config",
	Long:  `Remove a Jenkins config`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("You need to give a name")
		}

		target := args[0]
		if err := removeJenkins(target); err != nil {
			log.Fatal(err)
		}
	},
}

func removeJenkins(name string) (err error) {
	current := getCurrentJenkins()
	if name == current.Name {
		err = fmt.Errorf("You cannot remove current Jenkins, if you want to remove it, can select other items before")
		return
	}

	index := -1
	config := getConfig()
	for i, jenkins := range config.JenkinsServers {
		if name == jenkins.Name {
			index = i
			break
		}
	}

	if index == -1 {
		err = fmt.Errorf("Cannot found by name %s", name)
	} else {
		config.JenkinsServers = append(config.JenkinsServers[:index], config.JenkinsServers[index+1:]...)
		err = saveConfig()
	}
	return
}

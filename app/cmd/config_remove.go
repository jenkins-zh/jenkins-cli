package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/keyring"
	"go.uber.org/zap"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configRemoveCmd)
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: i18n.T("Remove a Jenkins config"),
	Long:  i18n.T("Remove a Jenkins config"),
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		target := args[0]
		return removeJenkins(target)
	},
}

func removeJenkins(name string) (err error) {
	current := getCurrentJenkins()
	if name == current.Name {
		err = fmt.Errorf("you cannot remove current Jenkins, if you want to remove it, can select other items before")
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
		err = fmt.Errorf("cannot found by name %s", name)
	} else {
		targetJenkins := config.JenkinsServers[index]
		if err = keyring.DelToken(targetJenkins); err == nil {
			logger.Info("delete keyring item successfully", zap.String("name", targetJenkins.Name))
		} else {
			logger.Error("cannot delete keyring item", zap.String("name", targetJenkins.Name),
				zap.Error(err))
		}

		config.JenkinsServers = append(config.JenkinsServers[:index], config.JenkinsServers[index+1:]...)
		err = saveConfig()
	}
	return
}

package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ConfigDataOptions is the config data option
type ConfigDataOptions struct {
	Key   string
	Value string
}

var configDataOptions ConfigDataOptions

func init() {
	configCmd.AddCommand(configDataCmd)
	configDataCmd.Flags().StringVarP(&configDataOptions.Key, "key", "k", "",
		i18n.T("The key of config data"))
	configDataCmd.Flags().StringVarP(&configDataOptions.Value, "value", "v", "",
		i18n.T("The value of config data"))

	configDataCmd.MarkFlagRequired("key")
	configDataCmd.MarkFlagRequired("value")
}

var configDataCmd = &cobra.Command{
	Use:               "data",
	Short:             i18n.T("Add a key/value to a config item"),
	Long:              i18n.T("Add a key/value to a config item"),
	ValidArgsFunction: ValidJenkinsNames,
	RunE: func(_ *cobra.Command, args []string) (err error) {
		var jenkinsName string
		if len(args) <= 0 {
			target := ""
			if currentJenkins := getCurrentJenkins(); currentJenkins != nil {
				target = currentJenkins.Name
			}
			jenkinsName, err = configSelectOptions.Select(getJenkinsNames(),
				"Choose a Jenkins to add key/value:", target)
		} else {
			jenkinsName = args[0]
		}

		found := false
		for i, cfg := range config.JenkinsServers {
			if cfg.Name == jenkinsName {
				if config.JenkinsServers[i].Data == nil {
					config.JenkinsServers[i].Data = make(map[string]string, 1)
				}
				config.JenkinsServers[i].Data[configDataOptions.Key] = configDataOptions.Value
				err = saveConfig()
				found = true
				break
			}
		}

		if !found {
			err = fmt.Errorf("jenkins '%s' does not exist", jenkinsName)
		}
		return
	},
	Annotations: map[string]string{
		common.Since: "v0.0.31",
	},
	Example: "jcli config data local -k jcli -v https://github.com/jenkins-zh/jenkins-cli",
}

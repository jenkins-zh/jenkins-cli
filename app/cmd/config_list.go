package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ConfigListOption option for config list command
type ConfigListOption struct {
	common.OutputOption

	Config string
}

var configListOption ConfigListOption

func init() {
	configCmd.AddCommand(configListCmd)
	configListCmd.Flags().StringVarP(&configListOption.Config, "config", "", "JenkinsServers",
		i18n.T("The type of config items, contains PreHooks, PostHooks, Mirrors, PluginSuites"))
	configListOption.SetFlagWithHeaders(configListCmd, "Name,URL,Description")
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("List all Jenkins config items"),
	Long:  i18n.T("List all Jenkins config items"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		configListOption.Writer = cmd.OutOrStdout()

		config := getConfig()
		if config == nil {
			return fmt.Errorf("no config file found")
		}

		switch configListOption.Config {
		case "JenkinsServers":
			err = configListOption.OutputV2(config.JenkinsServers)
		case "PreHooks":
			configListOption.Columns = "Path,Command"
			err = configListOption.OutputV2(config.PreHooks)
		case "PostHooks":
			configListOption.Columns = "Path,Command"
			err = configListOption.OutputV2(config.PostHooks)
		case "Mirrors":
			configListOption.Columns = "Name,URL"
			err = configListOption.OutputV2(config.Mirrors)
		case "PluginSuites":
			configListOption.Columns = "Name,Description"
			err = configListOption.OutputV2(config.PluginSuites)
		default:
			err = fmt.Errorf("unknow config %s", configListOption.Config)
		}
		return
	},
}

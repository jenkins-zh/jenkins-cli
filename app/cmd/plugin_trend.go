package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	pluginCmd.AddCommand(pluginTrendCmd)
}

var pluginTrendCmd = &cobra.Command{
	Use:   "trend <pluginName>",
	Short: "Show the trend of the plugin",
	Long:  `Show the trend of the plugin`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		pluginName := args[0]

		jclient := &client.PluginAPI{}
		jclient.ShowTrend(pluginName)
	},
}

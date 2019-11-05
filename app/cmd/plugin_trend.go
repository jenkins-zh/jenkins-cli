package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// PluginTreadOption is the option of plugin trend command
type PluginTreadOption struct {
	RoundTripper http.RoundTripper
}

var pluginTreadOption PluginTreadOption

func init() {
	pluginCmd.AddCommand(pluginTrendCmd)
}

var pluginTrendCmd = &cobra.Command{
	Use:   "trend <pluginName>",
	Short: "Show the trend of the plugin",
	Long:  `Show the trend of the plugin`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]

		jclient := &client.PluginAPI{
			RoundTripper: pluginTreadOption.RoundTripper,
		}
		tread, err := jclient.ShowTrend(pluginName)
		if err == nil {
			cmd.Print(tread)
		}
		helper.CheckErr(cmd, err)
	},
}

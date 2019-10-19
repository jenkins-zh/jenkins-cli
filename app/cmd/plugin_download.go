package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	pluginCmd.AddCommand(pluginDownloadCmd)
}

var pluginDownloadCmd = &cobra.Command{
	Use:   "download <keyword>",
	Short: "Download the plugins",
	Long:  `Download the plugins`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jclient := &client.PluginAPI{}
		jclient.DownloadPlugins(args)
	},
}

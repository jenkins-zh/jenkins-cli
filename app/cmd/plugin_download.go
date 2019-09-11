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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		jclient := &client.PluginAPI{}
		jclient.DownloadPlugins(args)
	},
}

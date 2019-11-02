package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type PluginDownloadOption struct {
	SkipDependency bool
	SkipOptional   bool
	UseMirror      bool
}

var pluginDownloadOption PluginDownloadOption

func init() {
	pluginCmd.AddCommand(pluginDownloadCmd)
	pluginDownloadCmd.Flags().BoolVarP(&pluginDownloadOption.SkipDependency, "skip-dependency", "", false,
		"If you want to skip download dependency of plugin")
	pluginDownloadCmd.Flags().BoolVarP(&pluginDownloadOption.SkipOptional, "skip-optional", "", true,
		"If you want to skip download optional dependency of plugin")
	pluginDownloadCmd.Flags().BoolVarP(&pluginDownloadOption.UseMirror, "use-mirror", "", true,
		"If you want to download plugin from a mirror site")
}

var pluginDownloadCmd = &cobra.Command{
	Use:   "download <keyword>",
	Short: "Download the plugins",
	Long:  `Download the plugins`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jclient := &client.PluginAPI{
			SkipDependency: pluginDownloadOption.SkipDependency,
			SkipOptional:   pluginDownloadOption.SkipOptional,
			UseMirror:      pluginDownloadOption.UseMirror,
			MirrorURL:      getDefaultMirror(),
		}
		jclient.DownloadPlugins(args)
	},
}

package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"net/http"
)

// PluginDownloadOption is the option for plugin download command
type PluginDownloadOption struct {
	SkipDependency bool
	SkipOptional   bool
	UseMirror      bool
	ShowProgress   bool

	RoundTripper http.RoundTripper
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
	pluginDownloadCmd.Flags().BoolVarP(&pluginDownloadOption.ShowProgress, "show-progress", "", true,
		"If you want to show the progress of download a plugin")
}

var pluginDownloadCmd = &cobra.Command{
	Use:   "download <keyword>",
	Short: "Download the plugins",
	Long:  `Download the plugins which contain the target plugin and its dependencies`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jClient := &client.PluginAPI{
			SkipDependency: pluginDownloadOption.SkipDependency,
			SkipOptional:   pluginDownloadOption.SkipOptional,
			UseMirror:      pluginDownloadOption.UseMirror,
			ShowProgress:   pluginDownloadOption.ShowProgress,
			MirrorURL:      getDefaultMirror(),
			RoundTripper:   pluginDownloadOption.RoundTripper,
		}
		jClient.DownloadPlugins(args)
	},
}

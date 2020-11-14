package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// PluginDownloadOption is the option for plugin download command
type PluginDownloadOption struct {
	SkipDependency bool
	SkipOptional   bool
	UseMirror      bool
	Mirror         string
	ShowProgress   bool
	DownloadDir    string

	RoundTripper http.RoundTripper
}

var pluginDownloadOption PluginDownloadOption

func init() {
	pluginCmd.AddCommand(pluginDownloadCmd)
	pluginDownloadCmd.Flags().BoolVarP(&pluginDownloadOption.SkipDependency, "skip-dependency", "", false,
		i18n.T("If you want to skip download dependency of plugin"))
	pluginDownloadCmd.Flags().BoolVarP(&pluginDownloadOption.SkipOptional, "skip-optional", "", true,
		i18n.T("If you want to skip download optional dependency of plugin"))
	pluginDownloadCmd.Flags().BoolVarP(&pluginDownloadOption.UseMirror, "use-mirror", "", true,
		i18n.T("If you want to download plugin from a mirror site"))
	pluginDownloadCmd.Flags().StringVarP(&pluginDownloadOption.Mirror, "mirror", "", "default",
		i18n.T("The mirror name"))
	pluginDownloadCmd.Flags().BoolVarP(&pluginDownloadOption.ShowProgress, "show-progress", "", true,
		i18n.T("If you want to show the progress of download a plugin"))
	pluginDownloadCmd.Flags().StringVarP(&pluginDownloadOption.DownloadDir, "download-dir", "", "",
		i18n.T("The directory which you want to download to"))
}

var pluginDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: i18n.T("Download the plugins"),
	Long: i18n.T(`Download the plugins which contain the target plugin and its dependencies.
It cannot collect its dependencies if you give a specific version.`),
	Args: cobra.MinimumNArgs(1),
	Example: `jcli plugin download localization-zh-cn
jcli plugin download localization-zh-cn@1.0.23`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jClient := &client.PluginAPI{
			SkipDependency: pluginDownloadOption.SkipDependency,
			SkipOptional:   pluginDownloadOption.SkipOptional,
			UseMirror:      pluginDownloadOption.UseMirror,
			ShowProgress:   pluginDownloadOption.ShowProgress,
			DownloadDir:    pluginDownloadOption.DownloadDir,
			RoundTripper:   pluginDownloadOption.RoundTripper,
		}
		if pluginDownloadOption.UseMirror {
			jClient.MirrorURL = getMirror(pluginDownloadOption.Mirror)
		}
		return jClient.DownloadPlugins(args)
	},
}

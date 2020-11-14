package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"net/http"
)

// PluginListOption option for plugin list command
type PluginListOption struct {
	common.OutputOption

	RoundTripper http.RoundTripper
}

var pluginListOption PluginListOption

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginListOption.SetFlagWithHeaders(pluginListCmd, "ShortName,Version,HasUpdate")
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("Print all the plugins which are installed"),
	Long:  i18n.T("Print all the plugins which are installed"),
	Example: `  jcli plugin list --filter ShortName=github
  jcli plugin list --no-headers`,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginListOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jClient.JenkinsCore))

		var plugins *client.InstalledPluginList
		if plugins, err = jClient.GetPlugins(1); err == nil {
			pluginListOption.Writer = cmd.OutOrStdout()
			err = pluginListOption.OutputV2(plugins.Plugins)
		}
		return
	},
}

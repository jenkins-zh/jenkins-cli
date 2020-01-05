package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// PluginUninstallOption the option of uninstall a plugin
type PluginUninstallOption struct {
	RoundTripper http.RoundTripper
}

var pluginUninstallOption PluginUninstallOption

func init() {
	pluginCmd.AddCommand(pluginUninstallCmd)
}

var pluginUninstallCmd = &cobra.Command{
	Use:   "uninstall [pluginName]",
	Short: i18n.T("Uninstall the plugins"),
	Long:  i18n.T("Uninstall the plugins"),
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginUninstallOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		err := jclient.UninstallPlugin(pluginName)
		helper.CheckErr(cmd, err)
	},
}

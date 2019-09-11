package cmd

import (
	"net/http"

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
	Short: "Uninstall the plugins",
	Long:  `Uninstall the plugins`,
	Run: func(cmd *cobra.Command, args []string) {
		var pluginName string
		if len(args) == 0 {
			cmd.Help()
			return
		}

		pluginName = args[0]

		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginUninstallOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if err := jclient.UninstallPlugin(pluginName); err != nil {
			cmd.PrintErr(err)
		}
	},
}

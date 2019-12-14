package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// PluginOptions contains the command line options
type PluginOptions struct {
	CommonOption

	Suite string
}

var pluginOpt PluginOptions

func init() {
	rootCmd.AddCommand(pluginCmd)
}

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: i18n.T("Manage the plugins of Jenkins"),
	Long:  i18n.T("Manage the plugins of Jenkins"),
	Example: `  jcli plugin list
  jcli plugin search github
  jcli plugin check`,
}

// FindPlugin find a plugin by name
func (o *PluginOptions) FindPlugin(name string) (plugin *client.InstalledPlugin, err error) {
	jClient := &client.PluginManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: o.RoundTripper,
		},
	}
	getCurrentJenkinsAndClient(&(jClient.JenkinsCore))
	if plugin, err = jClient.FindInstalledPlugin(name); err == nil && plugin == nil {
		err = fmt.Errorf(fmt.Sprintf("lack of plugin %s", name))
	}
	return
}

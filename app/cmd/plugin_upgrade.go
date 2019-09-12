package cmd

import (
	"net/http"
	"strings"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// PluginUpgradeOption option for plugin list command
type PluginUpgradeOption struct {
	Filter []string

	RoundTripper http.RoundTripper
}

var pluginUpgradeOption PluginUpgradeOption

func init() {
	pluginCmd.AddCommand(pluginUpgradeCmd)
	pluginUpgradeCmd.Flags().StringArrayVarP(&pluginUpgradeOption.Filter, "filter", "", []string{}, "Filter for the list, like: name=foo")
}

var pluginUpgradeCmd = &cobra.Command{
	Use:   "upgrade [plugin name]",
	Short: "Upgrade the specific plugin",
	Long:  `Upgrade the specific plugin`,
	Run: func(cmd *cobra.Command, args []string) {
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginUpgradeOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		var err error
		targetPlugins := make([]string, 0)
		if len(args) == 0 {
			var upgradeablePlugins []client.InstalledPlugin
			if upgradeablePlugins, err = pluginUpgradeOption.findUpgradeablePlugins(jclient); err == nil {
				prompt := &survey.MultiSelect{
					Message: fmt.Sprintf("Please select the plugins(%d) which you want to upgrade:", len(upgradeablePlugins)),
					Options: pluginUpgradeOption.convertToArray(upgradeablePlugins),
				}
				err = survey.AskOne(prompt, &targetPlugins)
			}
		} else {
			targetPlugins = args
		}

		if err != nil {
			cmd.PrintErrln(err)
		} else {
			if err = jclient.InstallPlugin(targetPlugins); err != nil {
				cmd.PrintErrln(err)
			}
		}
	},
}

func (p *PluginUpgradeOption) convertToArray(installedPlugins []client.InstalledPlugin) (plugins []string) {
	plugins = make([]string, 0)

	for _, plugin := range installedPlugins {
		plugins = append(plugins, plugin.ShortName)
	}
	return
}

func (p *PluginUpgradeOption) findUpgradeablePlugins(jclient *client.PluginManager) (
	filteredPlugins []client.InstalledPlugin, err error) {
	var (
		pluginName string
	)
	if p.Filter != nil {
		for _, f := range p.Filter {
			if strings.HasPrefix(f, "name=") {
				pluginName = strings.TrimPrefix(f, "name=")
			}
		}
	}

	var plugins *client.InstalledPluginList
	if plugins, err = jclient.GetPlugins(); err == nil {
		for _, plugin := range plugins.Plugins {
			if !plugin.HasUpdate {
				continue
			}

			if pluginName != "" && !strings.Contains(plugin.ShortName, pluginName) {
				continue
			}

			filteredPlugins = append(filteredPlugins, plugin)
		}
	}
	return
}

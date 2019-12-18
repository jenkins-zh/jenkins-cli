package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// PluginUpgradeOption option for plugin list command
type PluginUpgradeOption struct {
	Filter []string
	All    bool

	RoundTripper http.RoundTripper
}

var pluginUpgradeOption PluginUpgradeOption

func init() {
	pluginCmd.AddCommand(pluginUpgradeCmd)
	pluginUpgradeCmd.Flags().StringArrayVarP(&pluginUpgradeOption.Filter, "filter", "", []string{}, i18n.T("Filter for the list, like: name=foo"))
	pluginUpgradeCmd.Flags().BoolVarP(&pluginUpgradeOption.All, "all", "", false, i18n.T("Upgrade all plugins for updated"))

}

var pluginUpgradeCmd = &cobra.Command{
	Use:     "upgrade [plugin name]",
	Short:   i18n.T("Upgrade the specific plugin"),
	Long:    i18n.T("Upgrade the specific plugin"),
	Example: `jcli plugin upgrade [tab][tab]`,
	Run: func(cmd *cobra.Command, args []string) {
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginUpgradeOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		var err error
		targetPlugins := make([]string, 0)
		if pluginUpgradeOption.All {
			if upgradeablePlugins, err := pluginUpgradeOption.findUpgradeablePlugins(jclient); err == nil {
				targetPlugins = pluginUpgradeOption.convertToArray(upgradeablePlugins)
			}
		} else if len(args) == 0 {
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

		if err == nil && len(targetPlugins) != 0 {
			err = jclient.InstallPlugin(targetPlugins)
		}
		helper.CheckErr(cmd, err)
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
		plugins    *client.InstalledPluginList
	)
	if p.Filter != nil {
		for _, f := range p.Filter {
			if strings.HasPrefix(f, "name=") {
				pluginName = strings.TrimPrefix(f, "name=")
				break
			}
		}
	}

	if plugins, err = jclient.GetPlugins(1); err != nil {
		return
	}
	for _, plugin := range plugins.Plugins {
		if !plugin.HasUpdate {
			continue
		}

		if pluginName != "" && !strings.Contains(plugin.ShortName, pluginName) {
			continue
		}

		filteredPlugins = append(filteredPlugins, plugin)
	}
	return
}

/*func (p *PluginUpgradeOption) findCompatiblePlugins(installedPlugins []client.InstalledPlugin) (plugins []string) {
	plugins = make([]string, 0)
	var pluginNames string
	for i, plugin := range installedPlugins {
		if !strings.Contains(pluginNames, plugin.ShortName) {
			if len(installedPlugins) > i+1 {
				pluginNames += plugin.ShortName + "|"
			} else {
				pluginNames += plugin.ShortName
			}
		}
	}
	plugins = pluginUpgradeOption.assembleData(installedPlugins, pluginNames)
	return
}

func (p *PluginUpgradeOption) assembleData(installedPlugins []client.InstalledPlugin, pluginNames string) (plugins []string) {
	pluginAPI := client.PluginAPI{}
	if pluginsList, err := pluginAPI.BatchSearchPlugins(pluginNames); err == nil {
		for _, pluginInfo := range pluginsList {
			for _, plugin := range installedPlugins {
				if plugin.ShortName == pluginInfo.Name {
					var hasSecurity bool
					securityWarnings := pluginInfo.SecurityWarnings
					hasSecurity = p.checkSecurity(securityWarnings)
					if !hasSecurity {
						plugins = append(plugins, pluginInfo.Name)
					}
				}
			}
		}
	}
	return
}

func (p *PluginUpgradeOption) checkSecurity(securityWarnings []client.SecurityWarning) (hasSecurity bool) {
	for _, securityWarning := range securityWarnings {
		if securityWarning.Active {
			hasSecurity = true
			break
		}
	}
	return
}*/

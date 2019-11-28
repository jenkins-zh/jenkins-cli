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
	Filter     []string
	All        bool
	Compatible bool

	RoundTripper http.RoundTripper
}

var pluginUpgradeOption PluginUpgradeOption

func init() {
	pluginCmd.AddCommand(pluginUpgradeCmd)
	pluginUpgradeCmd.Flags().StringArrayVarP(&pluginUpgradeOption.Filter, "filter", "", []string{}, "Filter for the list, like: name=foo")
	pluginUpgradeCmd.Flags().BoolVarP(&pluginUpgradeOption.All, "all", "", false, "upgrade all plugins")
	pluginUpgradeCmd.Flags().BoolVarP(&pluginUpgradeOption.Compatible, "compatible", "", false, "upgrade all plugins")

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
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		var err error
		targetPlugins := make([]string, 0)
		if cmd.Flags() != nil && (cmd.Flag("all").Value.String() == "true" || cmd.Flag("compatible").Value.String() == "true") {
			if upgradeablePlugins, err := pluginUpgradeOption.findUpgradeablePlugins(jclient); err == nil {
				if cmd.Flag("all").Value.String() == "true" {
					targetPlugins = pluginUpgradeOption.convertToArray(upgradeablePlugins)
				} else {
					targetPlugins = pluginUpgradeOption.findCompatiblePlugins(upgradeablePlugins)
				}
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

func (p *PluginUpgradeOption) findCompatiblePlugins(installedPlugins []client.InstalledPlugin) (plugins []string) {
	plugins = make([]string, 0)

	for _, plugin := range installedPlugins {
		var hasSecurity bool
		pluginAPI := client.PluginAPI{}
		if pluginInfo, err := pluginAPI.GetPlugin(plugin.ShortName); err == nil {
			securityWarnings := pluginInfo.SecurityWarnings
			for _, securityWarning := range securityWarnings {
				if securityWarning.Active {
					hasSecurity = true
				}
			}
		}
		if !hasSecurity {
			plugins = append(plugins, plugin.ShortName)
		}
	}
	return
}

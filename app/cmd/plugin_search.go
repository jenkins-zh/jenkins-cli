package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// PluginSearchOption is the plugin search option
type PluginSearchOption struct {
	OutputOption

	RoundTripper http.RoundTripper
}

var pluginSearchOption PluginSearchOption

func init() {
	pluginCmd.AddCommand(pluginSearchCmd)
	pluginSearchCmd.PersistentFlags().StringVarP(&pluginSearchOption.Format, "output", "o", TableOutputFormat, "Format the output")
}

var pluginSearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Print the plugins of your Jenkins",
	Long:  `Print the plugins of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		keyword := args[0]

		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginSearchOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if plugins, err := jclient.GetAvailablePlugins(); err == nil {
			result := searchPlugins(plugins, keyword)
			resultData := matchPluginsData(result, jclient)
			if data, err := pluginSearchOption.Output(resultData); err == nil {
				if len(data) > 0 {
					cmd.Print(string(data))
				}
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

func searchPlugins(plugins *client.AvailablePluginList, keyword string) []client.AvailablePlugin {
	result := make([]client.AvailablePlugin, 0)

	for _, plugin := range plugins.Data {
		if strings.Contains(plugin.Name, strings.ToLower(keyword)) {
			result = append(result, plugin)
		}
	}
	return result
}

func matchPluginsData(plugins []client.AvailablePlugin, pluginJclient *client.PluginManager) (result []client.CenterPlugin) {
	if len(plugins) == 0 {
		return
	}
	result = make([]client.CenterPlugin, 0)
	jclient := &client.UpdateCenterManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: pluginSearchOption.RoundTripper,
		},
	}
	getCurrentJenkinsAndClient(&(jclient.JenkinsCore))
	site, err := jclient.GetSite()
	noSite := false
	if err != nil && site == nil {
		noSite = true
	}
	installedPlugins, err := pluginJclient.GetPlugins()
	noInstalledPlugin := false
	if err != nil && installedPlugins == nil {
		noInstalledPlugin = true
	}
	for _, plugin := range plugins {
		result = buildData(noSite, site, plugin, result, noInstalledPlugin, installedPlugins)
	}
	return
}

func buildData(noSite bool, site *client.CenterSite, plugin client.AvailablePlugin, result []client.CenterPlugin, noInstalledPlugin bool, installedPlugins *client.InstalledPluginList) (resultData []client.CenterPlugin) {
	isTrue := false
	if !noSite {
		if len(site.UpdatePlugins) > 0 {
			resultData, isTrue = buildUpdatePlugins(site, plugin, result)
		}

		if len(site.AvailablesPlugins) > 0 && !isTrue {
			resultData, isTrue = buildAvailablePlugins(site, plugin, result)
		}
	}
	if !noInstalledPlugin && len(installedPlugins.Plugins) > 0 && !isTrue {
		resultData, isTrue = buildInstalledPlugins(installedPlugins, plugin, result)
	}
	if !isTrue {
		resultData = buildNoMatchPlugins(plugin, result)
	}
	return
}

func buildUpdatePlugins(site *client.CenterSite, plugin client.AvailablePlugin, result []client.CenterPlugin) (resultData []client.CenterPlugin, isTrue bool) {
	isTrue = false
	resultData = result
	for _, updatePlugin := range site.UpdatePlugins {
		if plugin.Name == updatePlugin.Name {
			resultData = append(result, updatePlugin)
			isTrue = true
			break
		}
	}
	return
}

func buildAvailablePlugins(site *client.CenterSite, plugin client.AvailablePlugin, result []client.CenterPlugin) (resultData []client.CenterPlugin, isTrue bool) {
	resultData = result
	for _, availablePlugin := range site.AvailablesPlugins {
		if plugin.Name == availablePlugin.Name {
			resultData = append(result, availablePlugin)
			isTrue = true
			break
		}
	}
	return
}

func buildInstalledPlugins(installedPlugins *client.InstalledPluginList, plugin client.AvailablePlugin, result []client.CenterPlugin) (resultData []client.CenterPlugin, isTrue bool) {
	resultData = result
	for _, installPlugin := range installedPlugins.Plugins {
		if plugin.Name == installPlugin.ShortName {
			resultPlugin := client.CenterPlugin{}
			resultPlugin.CompatibleWithInstalledVersion = false
			resultPlugin.Name = installPlugin.ShortName
			resultPlugin.Installed.Active = true
			resultPlugin.Installed.Version = installPlugin.Version
			resultPlugin.Title = plugin.Title
			resultData = append(result, resultPlugin)
			isTrue = true
			break
		}
	}
	return
}

func buildNoMatchPlugins(plugin client.AvailablePlugin, result []client.CenterPlugin) (resultData []client.CenterPlugin) {
	resultData = result
	resultPlugin := client.CenterPlugin{}
	resultPlugin.CompatibleWithInstalledVersion = false
	resultPlugin.Name = plugin.Name
	resultPlugin.Installed.Active = plugin.Installed
	resultPlugin.Title = plugin.Title
	resultData = append(result, resultPlugin)
	return
}

// Output output the data into buffer
func (o *PluginSearchOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		pluginList := obj.([]client.CenterPlugin)
		buf := new(bytes.Buffer)

		if len(pluginList) != 0 {
			table := util.CreateTable(buf)
			table.AddRow("number", "name", "installed", "version", "installedVersion", "title")

			for i, plugin := range pluginList {
				formatTab(&table, i, plugin)
			}
			table.Render()
		}
		err = nil
		data = buf.Bytes()
	}
	return
}

func formatTab(table *util.Table, i int, plugin client.CenterPlugin) {
	installed := plugin.Installed
	if installed != (client.InstalledPlugin{}) {
		table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
			fmt.Sprintf("%t", true), plugin.Version, installed.Version, plugin.Title)
	} else {
		table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
			fmt.Sprintf("%t", false), plugin.Version, " ", plugin.Title)
	}
}

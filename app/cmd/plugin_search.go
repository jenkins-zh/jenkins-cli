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
			resultData := matchPluginsData(result)
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

func matchPluginsData(plugins []client.AvailablePlugin) (result []client.CenterPlugin) {
	result = make([]client.CenterPlugin, 0)
	jclient := &client.UpdateCenterManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: pluginSearchOption.RoundTripper,
		},
	}
	getCurrentJenkinsAndClient(&(jclient.JenkinsCore))
	site, err := jclient.GetSite()
	if err != nil {
		return
	}
	pluginJclient := &client.PluginManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: pluginSearchOption.RoundTripper,
		},
	}
	getCurrentJenkinsAndClient(&(pluginJclient.JenkinsCore))
	plu, err := pluginJclient.GetPlugins()
	if err != nil {
		return
	}

	for _, plugin := range plugins {
		for _, updatePlugin := range site.UpdatePlugins {
			if plugin.Name == updatePlugin.Name {
				result = append(result, updatePlugin)
				break
			}
		}
		for _, availablePlugin := range site.AvailablesPlugins {
			if plugin.Name == availablePlugin.Name {
				result = append(result, availablePlugin)
				break
			}
		}
		for _, pl := range plu.Plugins {
			if plugin.Name == pl.ShortName {
				s := client.CenterPlugin{}
				s.CompatibleWithInstalledVersion = false
				s.Name = pl.ShortName
				s.Installed.Active = true
				s.Installed.Version = pl.Version
				s.Title = plugin.Title
				result = append(result, s)
				break
			}
		}
	}
	return
}

func (o *PluginSearchOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		pluginList := obj.([]client.CenterPlugin)
		buf := new(bytes.Buffer)

		if len(pluginList) != 0 {
			table := util.CreateTable(buf)
			table.AddRow("number", "name", "installed", "version", "usedVersion", "title")

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
		if len(plugin.Version) > 6 && len(installed.Version) > 6 {
			table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
				fmt.Sprintf("%t", true), fmt.Sprintf("%v...", plugin.Version[0:6]), fmt.Sprintf("%v...", installed.Version[0:6]), plugin.Title)
		} else if len(plugin.Version) > 6 {
			table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
				fmt.Sprintf("%t", true), fmt.Sprintf("%v...", plugin.Version[0:6]), installed.Version, plugin.Title)
		} else if len(installed.Version) > 6 {
			table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
				fmt.Sprintf("%t", true), plugin.Version, fmt.Sprintf("%v...", installed.Version[0:6]), plugin.Title)
		} else {
			table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
				fmt.Sprintf("%t", true), plugin.Version, installed.Version, plugin.Title)
		}
	} else {
		if len(plugin.Version) > 6 {
			table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
				fmt.Sprintf("%t", false), fmt.Sprintf("%v...", plugin.Version[0:6]), " ", plugin.Title)
		} else {
			table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
				fmt.Sprintf("%t", false), plugin.Version, " ", plugin.Title)
		}
	}
}

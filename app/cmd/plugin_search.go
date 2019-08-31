package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

type PluginSearchOption struct {
	OutputOption
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

		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if plugins, err := jclient.GetAvailablePlugins(); err == nil {
			result := searchPlugins(plugins, keyword)

			if data, err := pluginSearchOption.Output(result, keyword); err == nil {
				if len(data) > 0 {
					fmt.Println(string(data))
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

func (o *PluginSearchOption) Output(obj interface{}, keyword string) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		pluginList := obj.([]client.AvailablePlugin)
		if len(pluginList) != 0 {
			table := util.CreateTable(os.Stdout)
			table.AddRow("number", "name", "installed", "version", "title")

			jclient := &client.PluginAPI{}
			plugins := jclient.SearchPlugins(keyword)
			for i, plugin := range pluginList {
				for _, plu := range plugins.Plugins {
					if plu.Name == plugin.Name && len(plu.Version) > 6 {
						table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
							fmt.Sprintf("%v", plugin.Installed), fmt.Sprintf("%v...", plu.Version[0:5]), plugin.Title)
						break
					} else if plu.Name == plugin.Name {
						table.AddRow(fmt.Sprintf("%d", i), plugin.Name,
							fmt.Sprintf("%v", plugin.Installed), plu.Version, plugin.Title)
						break
					}
				}
			}
			table.Render()
		}
		err = nil
		data = []byte{}
	}
	return
}

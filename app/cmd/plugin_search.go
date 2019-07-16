package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type PluginSearchOption struct {
	OutputOption
}

var pluginSearchOption PluginSearchOption

func init() {
	pluginCmd.AddCommand(pluginSearchCmd)
	pluginSearchCmd.PersistentFlags().StringVarP(&queueOption.Format, "output", "o", "json", "Format the output")
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

		jenkins := getCurrentJenkins()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if plugins, err := jclient.GetAvailablePlugins(); err == nil {
			result := searchPlugins(plugins, keyword)

			var data []byte
			if data, err = Format(result, queueOption.Format); err == nil {
				fmt.Printf("%s\n", string(data))
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
		if strings.Contains(plugin.Name, keyword) {
			result = append(result, plugin)
		}
	}
	return result
}

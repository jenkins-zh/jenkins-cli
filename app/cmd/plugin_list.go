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

type PluginListOption struct {
	OutputOption

	Filter []string
}

var pluginListOption PluginListOption

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginListCmd.Flags().StringArrayVarP(&pluginListOption.Filter, "filter", "", []string{}, "Filter for the list, like: active, hasUpdate, downgradable, enable, name=foo")
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all the plugins which are installed",
	Long:  `Print all the plugins which are installed`,
	Example: `  jcli plugin list --filter name=github
  jcli plugin list --filter hasUpdate`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		var (
			filter       bool
			hasUpdate    bool
			downgradable bool
			enable       bool
			active       bool
			pluginName   string
		)
		if pluginListOption.Filter != nil {
			filter = true
			for _, f := range pluginListOption.Filter {
				switch f {
				case "hasUpdate":
					hasUpdate = true
				case "downgradable":
					downgradable = true
				case "enable":
					enable = true
				case "active":
					active = true
				case "name":
					downgradable = true
				}

				if strings.HasPrefix(f, "name=") {
					pluginName = strings.TrimPrefix(f, "name=")
				}
			}
		}

		if plugins, err := jclient.GetPlugins(); err == nil {
			filteredPlugins := make([]client.InstalledPlugin, 0)
			for _, plugin := range plugins.Plugins {
				if filter {
					if hasUpdate && !plugin.HasUpdate {
						continue
					}

					if downgradable && !plugin.Downgradable {
						continue
					}

					if enable && !plugin.Enable {
						continue
					}

					if active && !plugin.Active {
						continue
					}

					if pluginName != "" && !strings.Contains(plugin.ShortName, pluginName) {
						continue
					}

					filteredPlugins = append(filteredPlugins, plugin)
				}
			}

			if data, err := pluginListOption.Output(filteredPlugins); err == nil {
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

func (o *PluginListOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		pluginList := obj.([]client.InstalledPlugin)
		table := util.CreateTable(os.Stdout)
		table.AddRow("number", "name", "version", "update")
		for i, plugin := range pluginList {
			table.AddRow(fmt.Sprintf("%d", i), plugin.ShortName, plugin.Version, fmt.Sprintf("%v", plugin.HasUpdate))
		}
		table.Render()
		err = nil
		data = []byte{}
	}
	return
}

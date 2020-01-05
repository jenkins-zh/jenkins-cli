package cmd

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"net/http"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// PluginListOption option for plugin list command
type PluginListOption struct {
	OutputOption

	Filter []string

	RoundTripper http.RoundTripper
}

var pluginListOption PluginListOption

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginListCmd.Flags().StringArrayVarP(&pluginListOption.Filter, "filter", "", []string{}, "Filter for the list, like: active, hasUpdate, downgradable, enable, name=foo")
	pluginListOption.SetFlag(pluginListCmd)
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("Print all the plugins which are installed"),
	Long:  i18n.T("Print all the plugins which are installed"),
	Example: `  jcli plugin list --filter name=github
  jcli plugin list --filter hasUpdate
  jcli plugin list --no-headers`,
	Run: func(cmd *cobra.Command, _ []string) {
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginListOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

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
				}

				if strings.HasPrefix(f, "name=") {
					pluginName = strings.TrimPrefix(f, "name=")
				}
			}
		}

		var err error
		var plugins *client.InstalledPluginList
		if plugins, err = jclient.GetPlugins(1); err == nil {
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

			var data []byte
			if data, err = pluginListOption.Output(filteredPlugins); err == nil && len(data) > 0 {
				cmd.Print(string(data))
			}
		}
		helper.CheckErr(cmd, err)
	},
}

// Output render data into byte array as a table format
func (o *PluginListOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil && o.Format == TableOutputFormat {
		buf := new(bytes.Buffer)

		pluginList := obj.([]client.InstalledPlugin)
		table := util.CreateTableWithHeader(buf, o.WithoutHeaders)
		table.AddHeader("number", "name", "version", "update")
		for i, plugin := range pluginList {
			table.AddRow(fmt.Sprintf("%d", i), plugin.ShortName, plugin.Version, fmt.Sprintf("%v", plugin.HasUpdate))
		}
		table.Render()
		err = nil
		data = buf.Bytes()
	}
	return
}

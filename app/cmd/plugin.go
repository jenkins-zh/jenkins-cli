package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/linuxsuren/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// PluginOptions contains the command line options
type PluginOptions struct {
	OutputOption

	Upload      bool
	CheckUpdate bool
	Open        bool
	List        bool

	Install   []string
	Uninstall string

	Filter []string
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.Flags().BoolVarP(&pluginOpt.Upload, "upload", "u", false, "Upload plugin to your Jenkins server")
	pluginCmd.Flags().BoolVarP(&pluginOpt.CheckUpdate, "check", "c", false, "Checkout update center server")
	pluginCmd.Flags().BoolVarP(&pluginOpt.Open, "open", "o", false, "Open the browse with the address of plugin manager")
	pluginCmd.Flags().BoolVarP(&pluginOpt.List, "list", "l", false, "Print all the plugins which are installed")
	pluginCmd.Flags().StringVarP(&pluginOpt.Format, "format", "", TableOutputFormat, "Format the output")
	pluginCmd.Flags().StringVarP(&pluginOpt.Uninstall, "uninstall", "", "", "Uninstall a plugin by shortName")
	pluginCmd.Flags().StringArrayVarP(&pluginOpt.Filter, "filter", "", []string{}, "Filter for the list, like: active, hasUpdate, downgradable, enable, name=foo")
}

var pluginOpt PluginOptions

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage the plugins of Jenkins",
	Long:  `Manage the plugins of Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if pluginOpt.Upload {
			jclient.Upload()
		}

		if pluginOpt.CheckUpdate {
			jclient.CheckUpdate(func(response *http.Response) {
				code := response.StatusCode
				if code == 200 {
					fmt.Println("update site updated.")
				} else {
					contentData, _ := ioutil.ReadAll(response.Body)
					log.Fatal(fmt.Sprintf("response code is %d, content: %s",
						code, string(contentData)))
				}
			})
		}

		if pluginOpt.Open {
			if jenkins.URL != "" {
				open(fmt.Sprintf("%s/pluginManager", jenkins.URL))
			} else {
				log.Fatal(fmt.Sprintf("No URL fond from %s", jenkins.Name))
			}
		}

		if pluginOpt.List {
			var (
				filter       bool
				hasUpdate    bool
				downgradable bool
				enable       bool
				active       bool
				pluginName   string
			)
			if pluginOpt.Filter != nil {
				filter = true
				for _, f := range pluginOpt.Filter {
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

				if data, err := pluginOpt.Output(filteredPlugins); err == nil {
					if len(data) > 0 {
						fmt.Println(string(data))
					}
				} else {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}

		if pluginOpt.Uninstall != "" {
			if err := jclient.UninstallPlugin(pluginOpt.Uninstall); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func (o *PluginOptions) Output(obj interface{}) (data []byte, err error) {
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

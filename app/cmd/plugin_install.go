package cmd

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginInstallCmd.Flags().StringVarP(&pluginOpt.Suite, "suite", "", "", "Suite of plugins")
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [pluginName]",
	Short: "Install the plugins",
	Long:  `Install the plugins`,
	Run: func(_ *cobra.Command, args []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth
		plugins := make([]string, len(args))
		plugins = append(plugins, args...)

		if pluginOpt.Suite != "" {
			if suite := findSuiteByName(pluginOpt.Suite); suite != nil {
				plugins = append(plugins, suite.Plugins...)
			} else {
				log.Fatal("Cannot found suite", pluginOpt.Suite)
			}
		}

		if len(plugins) == 0 {
			for {
				var keyword string
				prompt := &survey.Input{Message: "Please input the keyword to search your plugin!"}
				if err := survey.AskOne(prompt, &keyword); err != nil {
					log.Fatal(err)
				}

				if availablePlugins, err := jclient.GetAvailablePlugins(); err == nil {
					matchedPlugins := searchPlugins(availablePlugins, keyword)
					optinalPlugins := convertToArray(matchedPlugins)

					if len(optinalPlugins) == 0 {
						fmt.Println("Cannot find any plugins by your keyword, or they already installed.")
						continue
					}

					prompt := &survey.MultiSelect{
						Message: "Please select the plugins whose you want to install:",
						Options: convertToArray(matchedPlugins),
					}
					selectedPlugins := []string{}
					survey.AskOne(prompt, &selectedPlugins)
					plugins = append(plugins, selectedPlugins...)

					break
				} else {
					log.Fatal(err)
				}
			}

			fmt.Println("Going to install", plugins)
		}

		jclient.InstallPlugin(plugins)
	},
}

func convertToArray(aviablePlugins []client.AvailablePlugin) (plugins []string) {
	plugins = make([]string, 0)

	for _, plugin := range aviablePlugins {
		if plugin.Installed {
			continue
		}

		plugins = append(plugins, plugin.Name)
	}
	return
}

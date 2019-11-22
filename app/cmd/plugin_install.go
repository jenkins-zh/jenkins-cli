package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// PluginInstallOption is the option for plugin install
type PluginInstallOption struct {
	UseMirror    bool
	ShowProgress bool

	RoundTripper http.RoundTripper
}

var pluginInstallOption PluginInstallOption

func init() {
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginInstallCmd.Flags().BoolVarP(&pluginInstallOption.UseMirror, "use-mirror", "", true,
		i18n.T("If you want to download plugin from a mirror site"))
	pluginInstallCmd.Flags().BoolVarP(&pluginInstallOption.ShowProgress, "show-progress", "", true,
		i18n.T("If you want to show the progress of download a plugin"))
	pluginInstallCmd.Flags().StringVarP(&pluginOpt.Suite, "suite", "", "", "Suite of plugins")
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [pluginName]",
	Short: i18n.T("Install the plugins"),
	Long: i18n.T(`Install the plugins
Allow you to install a plugin with or without the version`),
	Example: `jcli plugin install localization-zh-cn
jcli plugin install localization-zh-cn@1.0.9
`,
	Run: func(cmd *cobra.Command, args []string) {
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginInstallOption.RoundTripper,
			},
			ShowProgress: pluginInstallOption.ShowProgress,
			UseMirror:    pluginInstallOption.UseMirror,
			MirrorURL:    getDefaultMirror(),
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		plugins := make([]string, len(args))
		plugins = append(plugins, args...)

		var err error
		if pluginOpt.Suite != "" {
			if suite := findSuiteByName(pluginOpt.Suite); suite != nil {
				plugins = append(plugins, suite.Plugins...)
			} else {
				err = fmt.Errorf("cannot found suite %s", pluginOpt.Suite)
			}
		}

		if err == nil && len(plugins) == 0 {
			for {
				var keyword string
				prompt := &survey.Input{Message: "Please input the keyword to search your plugin!"}
				if err = survey.AskOne(prompt, &keyword); err != nil {
					break
				}

				var availablePlugins *client.AvailablePluginList
				if availablePlugins, err = jclient.GetAvailablePlugins(); err == nil {
					matchedPlugins := searchPlugins(availablePlugins, keyword)
					optionalPlugins := convertToArray(matchedPlugins)

					if len(optionalPlugins) == 0 {
						cmd.Println("Cannot find any plugins by your keyword, or they already installed.")
						continue
					}

					prompt := &survey.MultiSelect{
						Message: "Please select the plugins whose you want to install:",
						Options: convertToArray(matchedPlugins),
					}
					selectedPlugins := []string{}
					if err = survey.AskOne(prompt, &selectedPlugins); err != nil {
						break
					}
					plugins = append(plugins, selectedPlugins...)
				}
				break
			}

			cmd.Println("Going to install", plugins)
		}

		if err == nil {
			err = jclient.InstallPlugin(plugins)
		}
		helper.CheckErr(cmd, err)
	},
}

func convertToArray(availablePlugins []client.AvailablePlugin) (plugins []string) {
	plugins = make([]string, 0)

	for _, plugin := range availablePlugins {
		if plugin.Installed {
			continue
		}

		plugins = append(plugins, plugin.Name)
	}
	return
}

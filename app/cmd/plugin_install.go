package cmd

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	jenkinsFormula "github.com/jenkins-zh/jenkins-formulas/pkg/common"
	"github.com/spf13/cobra"
)

// PluginInstallOption is the option for plugin install
type PluginInstallOption struct {
	UseMirror    bool
	ShowProgress bool
	Formula      string

	RoundTripper http.RoundTripper
}

var pluginInstallOption PluginInstallOption

func init() {
	pluginCmd.AddCommand(pluginInstallCmd)

	flags := pluginInstallCmd.Flags()
	flags.BoolVarP(&pluginInstallOption.UseMirror, "use-mirror", "", true,
		i18n.T("If you want to download plugin from a mirror site"))
	flags.BoolVarP(&pluginInstallOption.ShowProgress, "show-progress", "", true,
		i18n.T("If you want to show the progress of download a plugin"))
	flags.StringVarP(&pluginOpt.Suite, "suite", "", "", "Suite of plugins")
	flags.StringVarP(&pluginInstallOption.Formula, "formula", "", "",
		"Install plugins via a Jenkins formula. If you want to know more about Jenkins formula, please checkout https://github.com/jenkinsci/custom-war-packager")
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install",
	Short: i18n.T("Install the plugins"),
	Long: i18n.T(`Install the plugins
Allow you to install a plugin with or without the version`),
	Example: `jcli plugin install localization-zh-cn
jcli plugin install localization-zh-cn@1.0.9
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginInstallOption.RoundTripper,
			},
			ShowProgress: pluginInstallOption.ShowProgress,
			UseMirror:    pluginInstallOption.UseMirror,
			MirrorURL:    getDefaultMirror(),
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		plugins := make([]string, len(args))
		plugins = append(plugins, args...)

		if pluginOpt.Suite != "" {
			if suite := findSuiteByName(pluginOpt.Suite); suite != nil {
				plugins = append(plugins, suite.Plugins...)
			} else {
				err = fmt.Errorf("cannot found suite %s", pluginOpt.Suite)
			}
		}

		if pluginInstallOption.Formula != "" {
			var data []byte
			if data, err = ioutil.ReadFile(pluginInstallOption.Formula); err == nil {
				formula := jenkinsFormula.CustomWarPackage{}
				if err = yaml.Unmarshal(data, &formula); err == nil {
					for _, val := range formula.Plugins {
						plugins = append(plugins, val.ArtifactId)
					}
					logger.Info("prepare to install plugins", zap.Int("count", len(plugins)))
				}
			}

			if err != nil {
				err = fmt.Errorf("cannot read the Jenkins formular. %v", err)
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
		return
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

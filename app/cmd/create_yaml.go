package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	jenkinsFormula "github.com/jenkins-zh/jenkins-formulas/pkg/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var formulaYaml jenkinsFormula.CustomWarPackage
var pluginFormulaWithUpgradeOption PluginFormulaOption
var all bool

func init() {
	rootCmd.AddCommand(createYamlCmd)
	createYamlCmd.Flags().BoolVarP(&all, "all", "", false, i18n.T("Upgrade jenkins core and all plugins to update"))
	createYamlCmd.Flags().StringVarP(&formulaYaml.Bundle.GroupId, "bundle-groupId", "", "io.jenkins.tools.jcli.yaml.demo", i18n.T("GourpId of Bundle in yaml"))
	createYamlCmd.Flags().StringVarP(&formulaYaml.Bundle.ArtifactId, "bundle-artifactId", "", "jcli-yaml-demo", i18n.T("ArtifactId of Bundle in yaml"))
	createYamlCmd.Flags().StringVarP(&formulaYaml.Bundle.Vendor, "bundle-vendor", "", "jenkins-cli", i18n.T("Vendor of Bundle in yaml"))
	createYamlCmd.Flags().StringVarP(&formulaYaml.Bundle.Description, "bundle-description", "", "Upgraded jenkins core and plugins in a YAML specification", i18n.T("Description of Bundle in yaml"))
	healthCheckRegister.Register(getCmdPath(createYamlCmd), &pluginFormulaOption)
	createYamlCmd.Flags().BoolVarP(&pluginFormulaWithUpgradeOption.OnlyRelease, "only-release", "", true,
		`Indicated that we only output the release version of plugins`)
	createYamlCmd.Flags().BoolVarP(&pluginFormulaWithUpgradeOption.DockerBuild, "docker-build", "", false,
		`Indicated if build docker image`)
	createYamlCmd.Flags().BoolVarP(&pluginFormulaWithUpgradeOption.SortPlugins, "sort-plugins", "", true,
		`Indicated if sort the plugins by name`)
}

var createYamlCmd = &cobra.Command{
	Use:     "create yaml",
	Short:   i18n.T("Print a formula which contains all plugins come from current Jenkins server and upgraded plugins which were chosen by user"),
	Long:    i18n.T("Print a formula which contains all plugins come from current Jenkins server and upgraded plugins which were chosen by user"),
	Example: `create yaml --all
	create yaml`,
	RunE:    multipleChoice,
}

func multipleChoice(cmd *cobra.Command, args []string) (err error) {
	targetPlugins := make([]string, 0)
	formulaYaml.War.GroupId = "org.jenkins-ci.main"
	formulaYaml.War.ArtifactId = "jenkins-war"
	if all {
		if _, err := getLocalJenkinsAndPlugins(); err == nil {
			if pluginFormulaOption.OnlyRelease {
				formulaYaml.Plugins = removeSnapshotPluginsAndUpgradeOrNot(formulaYaml.Plugins, true)
			}
			if pluginFormulaOption.SortPlugins {
				formulaYaml.Plugins = SortPlugins(formulaYaml.Plugins)
			}
			if items, _, err := GetVersionData(LtsURL); err == nil {
				formulaYaml.War.Source.Version = "\"" + items[0].Title[8:] + "\""
			}
			formulaYaml.BuildSettings.Docker = jenkinsFormula.BuildDockerSetting{
				Base:  fmt.Sprintf("jenkins/jenkins:%s", formulaYaml.War.Source.Version),
				Tag:   "jenkins/jenkins-formula:v0.0.1",
				Build: pluginFormulaWithUpgradeOption.DockerBuild,
			}
		}
	} else if !all {
		var coreTemp bool
		if jenkinsVersion, err := getLocalJenkinsAndPlugins(); err == nil {
			promptCore := &survey.Confirm{
				Message: fmt.Sprintf("Please indicate whether do you want to upgrade or not"),
			}
			err = survey.AskOne(promptCore, &coreTemp)
			if err != nil {
				return err
			}
			if coreTemp {
				if items, _, err := GetVersionData(LtsURL); err == nil {
					formulaYaml.War.Source.Version = "\"" + items[0].Title[8:] + "\""
				}
			} else if !coreTemp {
				formulaYaml.War.Source.Version = jenkinsVersion
			}
			prompt := &survey.MultiSelect{
				Message: fmt.Sprintf("Please select the plugins(%d) which you want to upgrade to the latest: ", len(formulaYaml.Plugins)),
				Options: ConvertPluginsToArray(formulaYaml.Plugins),
			}
			err = survey.AskOne(prompt, &targetPlugins)

			if err != nil {
				return err
			}
			tempMap := make(map[string]bool)
			for _, plugin := range targetPlugins {
				tempMap[plugin] = true
			}
			for index, plugin := range formulaYaml.Plugins {
				if _, exist := tempMap[plugin.ArtifactId]; exist {
					formulaYaml.Plugins[index].Source.Version, err = getNewVersionOfPlugin(plugin.ArtifactId)
				}
				if err != nil {
					return err
				}
			}
			formulaYaml.BuildSettings.Docker = jenkinsFormula.BuildDockerSetting{
				Base:  fmt.Sprintf("jenkins/jenkins:%s", formulaYaml.War.Source.Version),
				Tag:   "jenkins/jenkins-formula:v0.0.1",
				Build: pluginFormulaWithUpgradeOption.DockerBuild,
			}
		}
	}
	renderYaml(formulaYaml)
	return nil
}

func getNewVersionOfPlugin(pluginName string) (version string, err error) {
	api := "https://plugins.jenkins.io/api/plugin/" + pluginName
	resp, err := http.Get(api)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	var newPluginOption NewPluginOption
	err = json.Unmarshal(bytes, &newPluginOption)
	if err != nil {
		return "", err
	}
	version = trimToID(newPluginOption.Version)

	return version, nil
}

func trimToID(content string) (version string) {
	startOfVersionNumber := strings.LastIndex(content, ":")
	version = content[startOfVersionNumber+1:]
	return version
}

func renderYaml(yamlTemp jenkinsFormula.CustomWarPackage) (err error) {
	data, err := yaml.Marshal(&yamlTemp)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("test.yaml", data, 0)
	if err != nil {
		return err
	}
	return nil
}

func removeSnapshotPluginsAndUpgradeOrNot(plugins []jenkinsFormula.Plugin, allTmp bool) (result []jenkinsFormula.Plugin) {
	result = make([]jenkinsFormula.Plugin, 0)
	if allTmp {
		for i := range plugins {
			if strings.Contains(plugins[i].Source.Version, "SNAPSHOT") {
				continue
			}
			plugins[i].Source.Version, _ = getNewVersionOfPlugin(plugins[i].ArtifactId)
			result = append(result, plugins[i])
		}
	} else if !allTmp {
		for i := range plugins {
			if strings.Contains(plugins[i].Source.Version, "SNAPSHOT") {
				continue
			}
			result = append(result, plugins[i])
		}
	}
	return result
}
func getLocalJenkinsAndPlugins() (jenkinsVersion string, err error) {
	jClientPlugin := &client.PluginManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: pluginListOption.RoundTripper,
		},
	}
	GetCurrentJenkinsAndClient(&(jClientPlugin.JenkinsCore))
	if err = jClientPlugin.GetPluginsFormula(&(formulaYaml.Plugins)); err != nil {
		return "", err
	}

	jClientCore := &client.JenkinsStatusClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: centerOption.RoundTripper,
		},
	}
	GetCurrentJenkinsAndClient(&(jClientCore.JenkinsCore))
	status, error := jClientCore.Get()
	if error != nil {
		return "", error
	}
	jenkinsVersion = status.Version
	return jenkinsVersion, nil
}
//ConvertPluginsToArray convert jenkinsFormula.Plugin to slice for the sake of multiple select
func ConvertPluginsToArray(plugins []jenkinsFormula.Plugin) (pluginArray []string) {
	pluginArray = make([]string, 0)
	for _, plugin := range plugins {
		pluginArray = append(pluginArray, plugin.ArtifactId)
	}
	return pluginArray
}

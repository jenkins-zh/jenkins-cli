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
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type YamlOption struct {
	Bundle  Bundle   `yaml:"bundle"`
	War     War      `yaml:"war"`
	Plugins []Plugin `yaml:"Plugins`
}
type Bundle struct {
	GroupId     string `yaml:"groupId"`
	ArtifactId  string `yaml:"artifactId"`
	Vendor      string `yaml:"vendor"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}
type War struct {
	GroupId    string `yaml:"groupId"`
	ArtifactId string `yaml:"artifactId"`
	Source     Source `yaml:"source"`
}
type Source struct {
	Version string `yaml:"version"`
}
type Plugin struct {
	GroupId    string `yaml:"groupId"`
	ArtifactId string `yaml:"artifactId"`
	Source     Source `yaml:"source"`
}
type coreAndPluginOption struct {
	plugin PluginUpgradeOption
	core   bool
	all    bool
}

var coreAndPlugin coreAndPluginOption
var yamlOption YamlOption

func init() {
	rootCmd.AddCommand(createYamlCmd)
	// createYamlCmd.Flags().StringArrayVarP(&coreAndPlugin.plugin.Filter, "filter", "", []string{}, i18n.T("Filter for the list"))
	createYamlCmd.Flags().BoolVarP(&coreAndPlugin.all, "all", "", false, i18n.T("Upgrade jenkins core and all plugins to update"))
	// createYamlCmd.Flags().BoolVarP(&coreAndPlugin.core, "core", "", false, i18n.T("Only upgrade jenkins core"))
	// createYamlCmd.Flags().BoolVarP(&coreAndPlugin.plugin.All, "plugin-all", "", false, i18n.T("Upgrade all plugins to update"))
	createYamlCmd.Flags().StringVarP(&yamlOption.Bundle.GroupId, "bundle-groupId", "", "io.jenkins.tools.jcli.yaml.demo", i18n.T("GourpId of Bundle in yaml"))
	createYamlCmd.Flags().StringVarP(&yamlOption.Bundle.ArtifactId, "bundle-artifactId", "", "jcli-yaml-demo", i18n.T("ArtifactId of Bundle in yaml"))
	createYamlCmd.Flags().StringVarP(&yamlOption.Bundle.Vendor, "bundle-vendor", "", "jenkins-cli", i18n.T("Vendor of Bundle in yaml"))
	createYamlCmd.Flags().StringVarP(&yamlOption.Bundle.Title, "bundle-title", "", "jcli create yaml demo", i18n.T("Title of Bundle in yaml"))
	createYamlCmd.Flags().StringVarP(&yamlOption.Bundle.Description, "bundle-description", "", "Upgraded jenkins core and plugins in a YAML specification", i18n.T("Description of Bundle in yaml"))
}

var createYamlCmd = &cobra.Command{
	Use:     "create yaml",
	Short:   "",
	Long:    "",
	Example: ``,
	RunE:    coreAndPlugin.multipleChoice,
}

func getLocalJenkinsAndPlugins() (jenkinsVersion string, pluginList []client.InstalledPlugin, err error) {
	jClientPlugin := &client.PluginManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: pluginListOption.RoundTripper,
		},
	}
	GetCurrentJenkinsAndClient(&(jClientPlugin.JenkinsCore))
	if plugins, err := jClientPlugin.GetPlugins(1); err == nil {
		pluginList = plugins.Plugins
	}
	jClientCore := &client.JenkinsStatusClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: centerOption.RoundTripper,
		},
	}
	GetCurrentJenkinsAndClient(&(jClientCore.JenkinsCore))
	status, error := jClientCore.Get()
	if error != nil {
		return "", nil, error
	}
	jenkinsVersion = status.Version
	return jenkinsVersion, pluginList, nil
}

func (c *coreAndPluginOption) multipleChoice(cmd *cobra.Command, args []string) (err error) {
	// var yamlOption YamlOption
	targetPlugins := make([]string, 0)
	yamlOption.War.GroupId = "org.jenkins-ci.main"
	yamlOption.War.ArtifactId = "jenkins-war"
	if c.all {
		if _, pluginList, err := getLocalJenkinsAndPlugins(); err == nil {
			yamlOption.Plugins = make([]Plugin, len(pluginList))
			for index, plugin := range pluginList {
				yamlOption.Plugins[index].GroupId, yamlOption.Plugins[index].ArtifactId, yamlOption.Plugins[index].Source.Version, err = getGroupIdAndArtifactId(plugin.ShortName)
				if err != nil {
					return err
				}
			}
			if items, _, err := GetVersionData(LtsURL); err == nil {
				yamlOption.War.Source.Version = "\"" + items[0].Title[8:] + "\""
			}
		}
		renderYaml(yamlOption)
		return nil
	} else if !c.all {
		var coreTemp bool
		if version, pluginList, err := getLocalJenkinsAndPlugins(); err == nil {
			yamlOption.Plugins = make([]Plugin, len(pluginList))
			promptCore := &survey.Confirm{
				Message: fmt.Sprintf("Please indicate whether do you want to upgrade or not"),
			}
			err = survey.AskOne(promptCore, &coreTemp)
			if err != nil {
				return err
			}
			if coreTemp {
				if items, _, err := GetVersionData(LtsURL); err == nil {
					yamlOption.War.Source.Version = "\"" + items[0].Title[8:] + "\""
				}
			} else if !coreTemp {
				yamlOption.War.Source.Version = "\"" + version + "\""
			}
			prompt := &survey.MultiSelect{
				Message: fmt.Sprintf("Please select the plugins(%d) which you want to upgrade to the latest: ", len(pluginList)),
				Options: coreAndPlugin.plugin.convertToArray(pluginList),
			}
			err = survey.AskOne(prompt, &targetPlugins)

			if err != nil {
				return err
			}
			tempMap := make(map[string]bool)
			for _, plugin := range targetPlugins {
				tempMap[plugin] = true
			}
			for index, plugin := range pluginList {
				if _, exist := tempMap[plugin.ShortName]; exist {
					yamlOption.Plugins[index].GroupId, yamlOption.Plugins[index].ArtifactId, yamlOption.Plugins[index].Source.Version, err = getGroupIdAndArtifactId(plugin.ShortName)
				} else {
					yamlOption.Plugins[index].GroupId, yamlOption.Plugins[index].ArtifactId, _, err = getGroupIdAndArtifactId(plugin.ShortName)
					yamlOption.Plugins[index].Source.Version = plugin.Version
				}
				if err != nil {
					return err
				}
			}
			renderYaml(yamlOption)
			return nil
		}
	}
	if c.all && len(targetPlugins) != 0 {
		cmd.Println("If you want to upgrade jenkins and all plugins, please use the flag --all. Otherwise, please do not append anything after `create yaml` and use the prompt instead.")
	}
	return nil
}

func getGroupIdAndArtifactId(pluginName string) (groupId string, artifactId string, version string, err error) {
	api := "https://plugins.jenkins.io/api/plugin/" + pluginName
	resp, err := http.Get(api)
	if err != nil {
		return "", "", "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", "", "", err
	}
	var newPluginOption NewPluginOption
	err = json.Unmarshal(bytes, &newPluginOption)
	if err != nil {
		return "", "", "", err
	}
	version, groupId, artifactId = trimToId(newPluginOption.Version)

	return groupId, artifactId, version, nil
}

func trimToId(content string) (version string, groupId string, artifactId string) {
	startOfVersionNumber := strings.LastIndex(content, ":")
	endOfGroupId := strings.Index(content, ":")
	version = content[startOfVersionNumber+1:]
	groupId = content[:endOfGroupId]
	artifactId = content[endOfGroupId+1 : startOfVersionNumber]
	return version, groupId, artifactId
}

func renderYaml(yamlTemp YamlOption) (err error) {
	bundle := Bundle{
		yamlOption.Bundle.GroupId,
		yamlOption.Bundle.ArtifactId,
		yamlOption.Bundle.Vendor,
		yamlOption.Bundle.Title,
		yamlOption.Bundle.Description,
	}
	yamlTemp.Bundle = bundle
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

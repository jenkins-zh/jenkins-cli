package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
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
	filter []string
	all    bool
	core   bool
}

var coreAndPlugin coreAndPluginOption

func init() {
	rootCmd.AddCommand(createYamlCmd)
	createYamlCmd.Flags().StringArrayVarP(&coreAndPlugin.filter, "filter", "", []string{}, i18n.T("Filter for the list"))
	createYamlCmd.Flags().BoolVarP(&coreAndPlugin.all, "all", "", false, i18n.T("Upgrade jenkins core and all plugins to update"))
	createYamlCmd.Flags().BoolVarP(&coreAndPlugin.core, "core", "", false, i18n.T("Only upgrade jenkins core"))
}

var createYamlCmd = &cobra.Command{
	Use:     "",
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
	var yamlOption YamlOption
	if c.all {
		if _, pluginList, err := getLocalJenkinsAndPlugins(); err == nil {
			for index, plugin := range pluginList {
				yamlOption.Plugins[index].GroupId, yamlOption.Plugins[index].ArtifactId, yamlOption.Plugins[index].Source.Version, err = getGroupIdAndArtifactId(plugin.ShortName)
				if err != nil {
					return err
				}
			}
			yamlOption.War.GroupId = "org.jenkins-ci.main"
			yamlOption.War.ArtifactId = "jenkins-war"
			if items, _, err := GetVersionData(LtsURL); err == nil {
				yamlOption.War.Source.Version = items[0].Title[8:]
			}
		}
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
	version, groupId, artifactId = trimToId(pluginName)

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

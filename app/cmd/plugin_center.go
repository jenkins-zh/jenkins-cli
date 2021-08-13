package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

//NewPluginOption consists of four options
type NewPluginOption struct {
	Name         string `json:"name"`
	Version      string `json:"gav"`
	Date         string `json:"releaseTimestamp"`
	RequiredCore string `json:"requiredCore"`
}

func init() {
	pluginCmd.AddCommand(pluginCenterCmd)
}

var pluginCenterCmd = &cobra.Command{
	Use:     "center",
	Short:   i18n.T("Print information about new version of the plugins which are installed"),
	Long:    i18n.T("Print information about new version of the plugins which are installed"),
	Example: `jcli plugin center`,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginListOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		var plugins *client.InstalledPluginList
		t := table.NewWriter()
		t.AppendHeader(table.Row{"ShortName", "Version", "Released Date", "Requires Jenkins"})
		if plugins, err = jClient.GetPlugins(1); err == nil {
			for _, plugin := range plugins.Plugins {
				version, date, requireCore, err := searchNewPlugin(plugin.ShortName)
				if err != nil {
					return err
				}
				if version != plugin.Version {
					t.AppendRow([]interface{}{plugin.ShortName, version, date, requireCore})
					t.AppendSeparator()
				}
			}
		}
		cmd.Print(t.Render())
		return
	},
}

func searchNewPlugin(pluginName string) (version string, date string, requireCore string, err error) {
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
	version = trimToVersionNumber(newPluginOption.Version)
	date = newPluginOption.Date[:10]
	requireCore = "jenkins " + newPluginOption.RequiredCore
	return version, date, requireCore, nil
}

func trimToVersionNumber(content string) string {
	startOfVersionNumber := strings.LastIndex(content, ":")
	return content[startOfVersionNumber+1:]
}

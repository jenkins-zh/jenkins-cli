package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

type ConfigDoctorOption struct {
	OutputOption

	RoundTripper http.RoundTripper
}

var configDoctorOption ConfigDoctorOption

func init() {
	configCmd.AddCommand(configDoctorCmd)
	configDoctorCmd.PersistentFlags().StringVarP(&configDoctorOption.Format, "output", "o", TableOutputFormat, "Format the output")
}

var configDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check unuseful Jenkins config items",
	Long:  `Check unuseful Jenkins config items`,
	Run: func(_ *cobra.Command, _ []string) {
		current := getCurrentJenkinsFromOptionsOrDie()

		table := util.CreateTable(os.Stdout)
		table.AddRow("number", "name", "url", "description", "status")
		for i, jenkins := range getConfig().JenkinsServers {
			jclient := &client.PluginManager{}
			jclient.URL = jenkins.URL
			jclient.UserName = jenkins.UserName
			jclient.Token = jenkins.Token
			jclient.Proxy = jenkins.Proxy
			jclient.ProxyAuth = jenkins.ProxyAuth
			plugin, err := jclient.GetAvailablePlugins()
			name := jenkins.Name
			if name == current.Name {
				name = fmt.Sprintf("*%s", name)
			}
			if err == nil && plugin.Status == "ok" {
				if len(jenkins.Description) > 15 {
					table.AddRow(fmt.Sprintf("%d", i), name, jenkins.URL, jenkins.Description[0:15], "available")
				} else {
					table.AddRow(fmt.Sprintf("%d", i), name, jenkins.URL, jenkins.Description, "available")
				}
			} else {
				if len(jenkins.Description) > 15 {
					table.AddRow(fmt.Sprintf("%d", i), name, jenkins.URL, jenkins.Description[0:15], "unavailable")
				} else {
					table.AddRow(fmt.Sprintf("%d", i), name, jenkins.URL, jenkins.Description, "unavailable")
				}
			}
		}
		table.Render()
	},
}

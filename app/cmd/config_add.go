package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ConfigAddOptions is the config ad option
type ConfigAddOptions struct {
	JenkinsServer
}

var configAddOptions ConfigAddOptions

func init() {
	configCmd.AddCommand(configAddCmd)
	configAddCmd.Flags().StringVarP(&configAddOptions.Name, "name", "n", "",
		i18n.T("Name of the Jenkins"))
	configAddCmd.Flags().StringVarP(&configAddOptions.URL, "url", "", "",
		i18n.T("URL of the Jenkins"))
	configAddCmd.Flags().StringVarP(&configAddOptions.UserName, "username", "u", "",
		i18n.T("UserName of the Jenkins"))
	configAddCmd.Flags().StringVarP(&configAddOptions.Token, "token", "t", "",
		i18n.T("Token of the Jenkins"))
	configAddCmd.Flags().StringVarP(&configAddOptions.Proxy, "proxy", "p", "",
		i18n.T("Proxy of the Jenkins"))
	configAddCmd.Flags().StringVarP(&configAddOptions.ProxyAuth, "proxyAuth", "a", "",
		i18n.T("ProxyAuth of the Jenkins"))
	configAddCmd.Flags().StringVarP(&configAddOptions.Description, "description", "d", "",
		i18n.T("Description of the Jenkins"))
}

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: i18n.T("Add a Jenkins config item"),
	Long:  i18n.T("Add a Jenkins config item"),
	RunE: func(_ *cobra.Command, _ []string) error {
		return addJenkins(configAddOptions.JenkinsServer)
	},
	Example: "jcli config add -n demo",
}

func addJenkins(jenkinsServer JenkinsServer) (err error) {
	jenkinsName := jenkinsServer.Name
	if jenkinsName == "" {
		err = fmt.Errorf("name cannot be empty")
		return
	}

	if findJenkinsByName(jenkinsName) != nil {
		err = fmt.Errorf("jenkins %s is existed", jenkinsName)
		return
	}

	config.JenkinsServers = append(config.JenkinsServers, jenkinsServer)
	err = saveConfig()
	return
}

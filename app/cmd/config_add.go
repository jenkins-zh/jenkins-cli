package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// ConfigAddOptions is the config ad option
type ConfigAddOptions struct {
	JenkinsServer
}

var configAddOptions ConfigAddOptions

func init() {
	configCmd.AddCommand(configAddCmd)
	configAddCmd.Flags().StringVarP(&configAddOptions.Name, "name", "n", "", "Name of the Jenkins")
	configAddCmd.Flags().StringVarP(&configAddOptions.URL, "url", "", "", "URL of the Jenkins")
	configAddCmd.Flags().StringVarP(&configAddOptions.UserName, "username", "u", "", "UserName of the Jenkins")
	configAddCmd.Flags().StringVarP(&configAddOptions.Token, "token", "t", "", "Token of the Jenkins")
	configAddCmd.Flags().StringVarP(&configAddOptions.Proxy, "proxy", "p", "", "Proxy of the Jenkins")
	configAddCmd.Flags().StringVarP(&configAddOptions.ProxyAuth, "proxyAuth", "a", "", "ProxyAuth of the Jenkins")
	configAddCmd.Flags().StringVarP(&configAddOptions.Description, "description", "d", "", "Description of the Jenkins")
}

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Jenkins config item",
	Long:  `Add a Jenkins config item`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := addJenkins(configAddOptions.JenkinsServer); err != nil {
			log.Fatal(err)
		}
	},
	Example: "jcli config add -n demo",
}

func addJenkins(jenkinsServer JenkinsServer) (err error) {
	jenkinsName := jenkinsServer.Name
	if jenkinsName == "" {
		err = fmt.Errorf("Name cannot be empty")
		return
	}

	if findJenkinsByName(jenkinsName) != nil {
		err = fmt.Errorf("Jenkins %s is existed", jenkinsName)
		return
	}

	config.JenkinsServers = append(config.JenkinsServers, jenkinsServer)
	err = saveConfig()
	return
}

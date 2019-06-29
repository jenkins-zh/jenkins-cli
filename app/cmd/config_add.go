package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

type ConfigAddOptions struct {
	JenkinsServer
}

var configAddOptions ConfigAddOptions

func init() {
	configCmd.AddCommand(configAddCmd)
	configAddCmd.PersistentFlags().StringVarP(&configAddOptions.Name, "name", "n", "", "Name of the Jenkins")
	configAddCmd.PersistentFlags().StringVarP(&configAddOptions.URL, "url", "", "", "URL of the Jenkins")
	configAddCmd.PersistentFlags().StringVarP(&configAddOptions.UserName, "username", "u", "", "UserName of the Jenkins")
	configAddCmd.PersistentFlags().StringVarP(&configAddOptions.Token, "token", "t", "", "Token of the Jenkins")
	configAddCmd.PersistentFlags().StringVarP(&configAddOptions.Proxy, "proxy", "p", "", "Proxy of the Jenkins")
	configAddCmd.PersistentFlags().StringVarP(&configAddOptions.ProxyAuth, "proxyAuth", "a", "", "ProxyAuth of the Jenkins")
}

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Jenkins config",
	Long:  `Add a Jenkins config`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := addJenkins(configAddOptions.JenkinsServer); err != nil {
			log.Fatal(err)
		}
	},
}

package cmd

import (
	"log"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobBuildCmd)
}

var jobBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the job of your Jenkins",
	Long:  `Build the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if jobOption.Name == "" {
			log.Fatal("need a name")
		}

		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		jclient.Build(jobOption.Name)
	},
}

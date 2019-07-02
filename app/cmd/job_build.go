package cmd

import (
	"fmt"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobBuildOption struct {
	BatchOption
}

var jobBuildOption JobBuildOption

func init() {
	jobCmd.AddCommand(jobBuildCmd)
	jobBuildCmd.Flags().BoolVarP(&jobBuildOption.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

var jobBuildCmd = &cobra.Command{
	Use:   "build -n",
	Short: "Build the job of your Jenkins",
	Long:  `Build the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if jobOption.Name == "" {
			cmd.Help()
			return
		}

		if !jobBuildOption.Confirm(fmt.Sprintf("Are you sure to build job %s", jobOption.Name)) {
			return
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

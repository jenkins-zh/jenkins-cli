package cmd

import (
	"fmt"
	"log"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobDeleteOption struct {
	BatchOption
}

var jobDeleteOption JobDeleteOption

func init() {
	jobCmd.AddCommand(jobDeleteCmd)
	jobDeleteOption.SetFlag(jobDeleteCmd)
}

var jobDeleteCmd = &cobra.Command{
	Use:   "delete <jobName>",
	Short: "Delete a job in your Jenkins",
	Long:  `Delete a job in your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		jobName := args[0]
		if !jobDeleteOption.Confirm(fmt.Sprintf("Are you sure to delete job %s ?", jobName)) {
			return
		}

		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if err := jclient.Delete(jobName); err != nil {
			log.Fatal(err)
		}
	},
}

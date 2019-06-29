package cmd

import (
	"fmt"
	"log"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobLogOption struct {
	History int
}

var jobLogOption JobLogOption

func init() {
	jobCmd.AddCommand(jobLogCmd)
	jobLogCmd.PersistentFlags().IntVarP(&jobLogOption.History, "history", "s", -1, "Specific build history of log")
}

var jobLogCmd = &cobra.Command{
	Use:   "log -n",
	Short: "Print the job of your Jenkins",
	Long:  `Print the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if jobOption.Name == "" {
			cmd.Help()
			return
		}

		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		printLog(jclient, jobOption.Name, jobLogOption.History, 0)
	},
}

func printLog(jclient *client.JobClient, jobName string, history int, start int64) {
	if status, err := jclient.Log(jobName, history, start); err == nil {
		fmt.Print(status.Text)
		if status.HasMore {
			printLog(jclient, jobName, history, status.NextStart)
		}
	} else {
		log.Fatal(err)
	}
}

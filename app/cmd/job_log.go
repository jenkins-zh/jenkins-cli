package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobLogOption struct {
	WatchOption
	History int

	LogText      string
	LastBuildID  int
	LastBuildURL string
}

var jobLogOption JobLogOption

func init() {
	jobCmd.AddCommand(jobLogCmd)
	jobLogCmd.Flags().IntVarP(&jobLogOption.History, "history", "s", -1, "Specific build history of log")
	jobLogCmd.Flags().BoolVarP(&jobLogOption.Watch, "watch", "w", false, "Watch the job logs")
	jobLogCmd.Flags().IntVarP(&jobLogOption.Interval, "interval", "i", 1, "Interval of watch")
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

		lastBuildID := -1
		for {
			if build, err := jclient.GetBuild(jobOption.Name, -1); err == nil {
				jobLogOption.LastBuildID = build.Number
				jobLogOption.LastBuildURL = build.URL
			}
			if lastBuildID != jobLogOption.LastBuildID {
				lastBuildID = jobLogOption.LastBuildID
				fmt.Println("Current build number:", jobLogOption.LastBuildID)
				fmt.Println("Current build url:", jobLogOption.LastBuildURL)

				printLog(jclient, jobOption.Name, jobLogOption.History, 0)
			}

			if !jobLogOption.Watch {
				break
			}

			time.Sleep(time.Duration(jobLogOption.Interval) * time.Second)
		}
	},
}

func printLog(jclient *client.JobClient, jobName string, history int, start int64) {
	if status, err := jclient.Log(jobName, history, start); err == nil {
		isNew := false

		if jobLogOption.LogText != status.Text {
			jobLogOption.LogText = status.Text
			isNew = true
		} else if history == -1 {
			if build, err := jclient.GetBuild(jobName, -1); err == nil && jobLogOption.LastBuildID != build.Number {
				jobLogOption.LastBuildID = build.Number
				jobLogOption.LastBuildURL = build.URL
				isNew = true
			}
		}

		if isNew {
			fmt.Print(status.Text)
		}

		if status.HasMore {
			printLog(jclient, jobName, history, status.NextStart)
		}
	} else {
		log.Fatal(err)
	}
}

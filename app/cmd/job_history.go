package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/linuxsuren/jenkins-cli/util"
	"github.com/spf13/cobra"
)

type JobHistoryOption struct {
	OutputOption
}

var jobHistoryOption JobHistoryOption

func init() {
	jobCmd.AddCommand(jobHistoryCmd)
	jobHistoryOption.SetFlag(jobHistoryCmd)
}

var jobHistoryCmd = &cobra.Command{
	Use:   "history <jobName>",
	Short: "Print the history of job in your Jenkins",
	Long:  `Print the history of job in your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		jobName := args[0]

		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if builds, err := jclient.GetHistory(jobName); err == nil {
			if data, err := jobHistoryOption.Output(builds); err == nil {
				if len(data) > 0 {
					fmt.Println(string(data))
				}
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

func (o *JobHistoryOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		buildList := obj.([]client.JobBuild)
		table := util.CreateTable(os.Stdout)
		table.AddRow("number", "displayname", "building", "result")
		for i, build := range buildList {
			table.AddRow(fmt.Sprintf("%d", i), build.DisplayName,
				fmt.Sprintf("%v", build.Building), build.Result)
		}
		table.Render()
		err = nil
		data = []byte{}
	}
	return
}

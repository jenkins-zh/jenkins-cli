package cmd

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"net/http"
)

// JobHistoryOption is the job history option
type JobHistoryOption struct {
	OutputOption

	RoundTripper http.RoundTripper
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
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobName := args[0]

		jClient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobHistoryOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		builds, err := jClient.GetHistory(jobName)
		if err == nil {
			var data []byte
			data, err = jobHistoryOption.Output(builds)
			if err == nil && len(data) > 0 {
				cmd.Print(string(data))
			}
		}
		helper.CheckErr(cmd, err)
	},
}

// Output print the output
func (o *JobHistoryOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		buildList := obj.([]*client.JobBuild)
		buf := new(bytes.Buffer)
		table := util.CreateTable(buf)
		table.AddRow("number", "displayname", "building", "result")
		for i, build := range buildList {
			table.AddRow(fmt.Sprintf("%d", i), build.DisplayName,
				fmt.Sprintf("%v", build.Building), build.Result)
		}
		table.Render()
		data = buf.Bytes()
		err = nil
	}
	return
}

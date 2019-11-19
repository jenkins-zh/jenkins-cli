package cmd

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/util"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// QueueListOption represents the option of queue list command
type QueueListOption struct {
	OutputOption

	RoundTripper http.RoundTripper
}

var queueListOption QueueListOption

func init() {
	queueCmd.AddCommand(queueListCmd)
	queueListOption.SetFlag(queueListCmd)
}

var queueListCmd = &cobra.Command{
	Use:   "list",
	Short: "Print the queue of your Jenkins",
	Long:  `Print the queue of your Jenkins`,
	Run: func(cmd *cobra.Command, _ []string) {
		jclient := &client.QueueClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: queueListOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		var err error
		var jobQueue *client.JobQueue
		if jobQueue, err = jclient.Get(); err == nil {
			var data []byte
			if data, err = queueListOption.Output(jobQueue); err == nil && len(data) > 0 {
				cmd.Print(string(data))
			}
		}
		helper.CheckErr(cmd, err)
	},
}

// Output render data into byte array as a table format
func (o *QueueListOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil && o.Format == TableOutputFormat {
		buf := new(bytes.Buffer)

		jobQueue := obj.(*client.JobQueue)
		table := util.CreateTableWithHeader(buf, o.WithoutHeaders)
		table.AddHeader("number", "id", "why", "url")
		for i, item := range jobQueue.Items {
			table.AddRow(fmt.Sprintf("%d", i), fmt.Sprintf("%d", item.ID), item.Why, item.URL)
		}
		table.Render()
		err = nil
		data = buf.Bytes()
	}
	return
}

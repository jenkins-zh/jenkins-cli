package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	cobra_ext "github.com/linuxsuren/cobra-extension/pkg"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// QueueListOption represents the option of queue list command
type QueueListOption struct {
	cobra_ext.OutputOption

	RoundTripper http.RoundTripper
}

var queueListOption QueueListOption

func init() {
	queueCmd.AddCommand(queueListCmd)
	queueListOption.SetFlagWithHeaders(queueListCmd, "ID,Why,URL")
}

var queueListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("Print the queue of your Jenkins"),
	Long:  i18n.T("Print the queue of your Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.QueueClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: queueListOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		var jobQueue *client.JobQueue
		if jobQueue, err = jClient.Get(); err == nil {
			queueListOption.Writer = cmd.OutOrStdout()
			err = queueListOption.OutputV2(jobQueue.Items)
		}
		return
	},
}

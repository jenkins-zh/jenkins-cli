package cmd

import (
	"net/http"
	"strconv"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// QueueCancelOption represents the option of queue cancel command
type QueueCancelOption struct {
	RoundTripper http.RoundTripper
}

var queueCancelOption QueueCancelOption

func init() {
	queueCmd.AddCommand(queueCancelCmd)
}

var queueCancelCmd = &cobra.Command{
	Use:   "cancel <id>",
	Short: "Cancel the queue of your Jenkins",
	Long:  `Cancel the queue of your Jenkins`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			queueID int
			err     error
		)
		if queueID, err = strconv.Atoi(args[0]); err != nil {
			cmd.PrintErrln(err)
			return
		}

		jclient := &client.QueueClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: queueCancelOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		err = jclient.Cancel(queueID)
		helper.CheckErr(cmd, err)
	},
}

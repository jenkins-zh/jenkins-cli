package cmd

import (
	"strconv"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// QueueCancelOption represents the option of queue cancel command
type QueueCancelOption struct {
}

var queueCancelOption QueueCancelOption

func init() {
	queueCmd.AddCommand(queueCancelCmd)
}

var queueCancelCmd = &cobra.Command{
	Use:   "cancel <id>",
	Short: "Cancel the queue of your Jenkins",
	Long:  `Cancel the queue of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

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
				RoundTripper: pluginListOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if err = jclient.Cancel(queueID); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

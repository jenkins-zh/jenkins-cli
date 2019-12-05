package cmd

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"

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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		for _, arg := range args {
			if err = queueCancelOption.cancel(arg); err != nil {
				break
			}
		}
		return
	},
}

func (c *QueueCancelOption) cancel(id string) (err error) {
	var (
		queueID int
	)

	if queueID, err = strconv.Atoi(id); err != nil {
		return
	}

	jclient := &client.QueueClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: queueCancelOption.RoundTripper,
			Debug:        rootOptions.Debug,
		},
	}
	getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

	logger.Debug("cancel queue by id,", zap.Int("id", queueID))

	err = jclient.Cancel(queueID)
	return
}

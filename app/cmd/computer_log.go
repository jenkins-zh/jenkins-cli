package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerLogOption option for config list command
type ComputerLogOption struct {
	CommonOption
}

var computerLogOption ComputerLogOption

func init() {
	computerCmd.AddCommand(computerLogCmd)
}

var computerLogCmd = &cobra.Command{
	Use:   "log <name>",
	Short: i18n.T("Output the log of the agent"),
	Long:  i18n.T("Output the log of the agent"),
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) (err error) {
		jClient := &client.ComputerClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: computerLogOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		var log string
		if log, err = jClient.GetLog(args[0]); err == nil {
			fmt.Print(log)
		}
		return
	},
}

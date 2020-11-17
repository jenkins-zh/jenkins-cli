package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerLogOption option for config list command
type ComputerLogOption struct {
	common.Option
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jClient, _ := GetComputerClient(computerLogOption.Option)

		var log string
		if log, err = jClient.GetLog(args[0]); err == nil {
			cmd.Print(log)
		}
		return
	},
}

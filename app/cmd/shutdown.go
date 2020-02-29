package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// ShutdownOption holds the options for shutdown cmd
type ShutdownOption struct {
	BatchOption
	CommonOption

	Safe          bool
	Prepare       bool
	CancelPrepare bool
}

var shutdownOption ShutdownOption

func init() {
	rootCmd.AddCommand(shutdownCmd)
	shutdownOption.SetFlag(shutdownCmd)
	shutdownCmd.Flags().BoolVarP(&shutdownOption.Safe, "safe", "s", true,
		i18n.T("Puts Jenkins into the quiet mode, wait for existing builds to be completed, and then shut down Jenkins"))
	shutdownCmd.Flags().BoolVarP(&shutdownOption.Prepare, "prepare", "", false,
		i18n.T("Put Jenkins in a Quiet mode, in preparation for a restart. In that mode Jenkins don’t start any build"))
	shutdownCmd.Flags().BoolVarP(&shutdownOption.CancelPrepare, "prepare-cancel", "", false,
		i18n.T(" Cancel the effect of the “quiet-down” command"))
	shutdownOption.BatchOption.Stdio = GetSystemStdio()
	shutdownOption.CommonOption.Stdio = GetSystemStdio()
}

var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: i18n.T("Puts Jenkins into the quiet mode, wait for existing builds to be completed, and then shut down Jenkins"),
	Long:  i18n.T("Puts Jenkins into the quiet mode, wait for existing builds to be completed, and then shut down Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jenkins := getCurrentJenkinsFromOptions()
		if !shutdownOption.Confirm(fmt.Sprintf("Are you sure to shutdown Jenkins %s?", jenkins.URL)) {
			return
		}

		jClient := &client.CoreClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: shutdownOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		if shutdownOption.CancelPrepare {
			err = jClient.PrepareShutdown(true)
		} else if shutdownOption.Prepare {
			err = jClient.PrepareShutdown(false)
		} else {
			err = jClient.Shutdown(shutdownOption.Safe)
		}
		return
	},
}

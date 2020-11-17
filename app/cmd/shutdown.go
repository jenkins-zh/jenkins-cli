package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// ShutdownOption holds the options for shutdown cmd
type ShutdownOption struct {
	common.BatchOption
	common.Option
	RootOptions *RootOptions

	Safe          bool
	Prepare       bool
	CancelPrepare bool
}

const (
	// SafeShutdown the text about shutdown safely
	SafeShutdown = "Puts Jenkins into the quiet mode, wait for existing builds to be completed, and then shut down Jenkins"
)

// NewShutdownCmd create the shutdown command
func NewShutdownCmd(rootOpt *RootOptions) (cmd *cobra.Command) {
	shutdownOption := &ShutdownOption{
		RootOptions: rootOpt,
	}
	cmd = &cobra.Command{
		Use:   "shutdown",
		Short: i18n.T(SafeShutdown),
		Long:  i18n.T(SafeShutdown),
		RunE:  shutdownOption.runE,
	}
	shutdownOption.init(cmd)
	return
}

func (o *ShutdownOption) runE(cmd *cobra.Command, _ []string) (err error) {
	jenkins := getCurrentJenkinsFromOptions()
	if !o.Confirm(fmt.Sprintf("Are you sure to shutdown Jenkins %s?", jenkins.URL)) {
		return
	}

	jClient := &client.CoreClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: o.RootOptions.CommonOption.RoundTripper,
			Debug:        o.RootOptions.Debug,
		},
	}
	getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

	if o.CancelPrepare {
		err = jClient.PrepareShutdown(true)
	} else if o.Prepare {
		err = jClient.PrepareShutdown(false)
	} else {
		err = jClient.Shutdown(o.Safe)
	}
	return
}

func (o *ShutdownOption) init(shutdownCmd *cobra.Command) {
	rootCmd.AddCommand(shutdownCmd)
	o.SetFlag(shutdownCmd)
	shutdownCmd.Flags().BoolVarP(&o.Safe, "safe", "s", true,
		i18n.T(SafeShutdown))
	shutdownCmd.Flags().BoolVarP(&o.Prepare, "prepare", "", false,
		i18n.T("Put Jenkins in a Quiet mode, in preparation for a restart. In that mode Jenkins don’t start any build"))
	shutdownCmd.Flags().BoolVarP(&o.CancelPrepare, "prepare-cancel", "", false,
		i18n.T(" Cancel the effect of the “quiet-down” command"))
	o.BatchOption.Stdio = common.GetSystemStdio()
	o.Option.Stdio = common.GetSystemStdio()
}

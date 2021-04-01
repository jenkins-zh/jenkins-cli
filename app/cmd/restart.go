package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// RestartOption holds the options for restart cmd
type RestartOption struct {
	common.BatchOption
	common.Option

	Safe   bool
	Reload bool
}

var restartOption RestartOption

func init() {
	rootCmd.AddCommand(restartCmd)
	restartOption.SetFlag(restartCmd)
	restartCmd.Flags().BoolVarP(&restartOption.Safe, "safe", "s", true,
		i18n.T("Puts Jenkins into the quiet mode, wait for existing builds to be completed, and then restart Jenkins"))
	restartCmd.Flags().BoolVarP(&restartOption.Reload, "reload", "r", false,
		i18n.T("Reload configuration from disk, this action would not restart your Jenkins"))
	restartOption.BatchOption.Stdio = common.GetSystemStdio()
	restartOption.Option.Stdio = common.GetSystemStdio()
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: i18n.T("Restart your Jenkins"),
	Long:  i18n.T("Restart your Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jenkins := GetCurrentJenkinsFromOptions()
		if !restartOption.Confirm(fmt.Sprintf("Are you sure to restart/reload Jenkins %s?", jenkins.URL)) {
			return
		}

		jClient := &client.CoreClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: restartOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		if restartOption.Reload {
			err = jClient.Reload()
		} else if restartOption.Safe {
			err = jClient.Restart()
		} else {
			err = jClient.RestartDirectly()
		}

		if err == nil {
			if restartOption.Reload {
				cmd.Println("Please wait while Jenkins is reloading")
			} else {
				cmd.Println("Please wait while Jenkins is restarting")
			}
		}
		return
	},
}

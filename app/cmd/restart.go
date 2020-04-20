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
	common.CommonOption

	Safe bool
}

var restartOption RestartOption

func init() {
	rootCmd.AddCommand(restartCmd)
	restartOption.SetFlag(restartCmd)
	restartCmd.Flags().BoolVarP(&restartOption.Safe, "safe", "s", true,
		i18n.T("Puts Jenkins into the quiet mode, wait for existing builds to be completed, and then restart Jenkins"))
	restartOption.BatchOption.Stdio = common.GetSystemStdio()
	restartOption.CommonOption.Stdio = common.GetSystemStdio()
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: i18n.T("Restart your Jenkins"),
	Long:  i18n.T("Restart your Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jenkins := getCurrentJenkinsFromOptions()
		if !restartOption.Confirm(fmt.Sprintf("Are you sure to restart Jenkins %s?", jenkins.URL)) {
			return
		}

		jClient := &client.CoreClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: restartOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		if restartOption.Safe {
			err = jClient.Restart()
		} else {
			err = jClient.RestartDirectly()
		}

		if err == nil {
			cmd.Println("Please wait while Jenkins is restarting")
		}
		return
	},
}

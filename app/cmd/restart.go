package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey"
	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// RestartOption holds the options for restart cmd
type RestartOption struct {
	BatchOption
}

var restartOption RestartOption

func init() {
	rootCmd.AddCommand(restartCmd)
	restartCmd.Flags().BoolVarP(&restartOption.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart your Jenkins",
	Long:  `Restart your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		if !restartOption.Batch {
			confirm := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Are you sure to restart Jenkins %s?", jenkins.URL),
			}
			survey.AskOne(prompt, &confirm)
			if !confirm {
				return
			}
		}

		jclient := &client.CoreClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		jclient.Restart()
	},
}

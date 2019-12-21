package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerCreateOption option for config list command
type ComputerCreateOption struct {
	OutputOption

	RoundTripper http.RoundTripper
}

var computerCreateOption ComputerCreateOption

func init() {
	computerCmd.AddCommand(computerCreateCmd)
}

var computerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: i18n.T("Create an Jenkins agent"),
	Long: i18n.T(`Create an Jenkins agent
It can only create a JNLP agent.`),
	Example: `jcli agent create agent-name`,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jClient := &client.ComputerClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: computerCreateOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		err = jClient.Create(args[0])
		return
	},
	Annotations: map[string]string{
		since: "v0.0.24",
	},
}

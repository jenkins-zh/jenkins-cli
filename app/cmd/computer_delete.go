package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerDeleteOption option for agent delete command
type ComputerDeleteOption struct {
	CommonOption
}

var computerDeleteOption ComputerDeleteOption

func init() {
	computerCmd.AddCommand(computerDeleteCmd)
}

var computerDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: GetAliasesDel(),
	Short:   i18n.T("Delete an agent from Jenkins"),
	Long:    i18n.T("Delete an agent from Jenkins"),
	Args:    cobra.MinimumNArgs(1),
	Example: `jcli agent delete agent-name`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jClient := &client.ComputerClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: computerDeleteOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		err = jClient.Delete(args[0])
		return
	},
	Annotations: map[string]string{
		since: "v0.0.24",
	},
}

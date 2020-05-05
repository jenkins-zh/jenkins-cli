package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerCreateOption option for config list command
type ComputerCreateOption struct {
	common.CommonOption
	common.OutputOption
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
		jClient, _ := GetComputerClient(computerCreateOption.CommonOption)
		return jClient.Create(args[0])
	},
	Annotations: map[string]string{
		common.Since: common.VersionSince0024,
	},
}

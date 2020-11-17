package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerDeleteOption option for agent delete command
type ComputerDeleteOption struct {
	common.Option
}

var computerDeleteOption ComputerDeleteOption

func init() {
	computerCmd.AddCommand(computerDeleteCmd)
}

var computerDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: common.GetAliasesDel(),
	Short:   i18n.T("Delete an agent from Jenkins"),
	Long:    i18n.T("Delete an agent from Jenkins"),
	Args:    cobra.MinimumNArgs(1),
	Example: `jcli agent delete agent-name`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jClient, _ := GetComputerClient(computerDeleteOption.Option)
		return jClient.Delete(args[0])
	},
	Annotations: map[string]string{
		common.Since: common.VersionSince0024,
	},
}

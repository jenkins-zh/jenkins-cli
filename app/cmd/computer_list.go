package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerListOption option for config list command
type ComputerListOption struct {
	common.CommonOption
	common.OutputOption
}

var computerListOption ComputerListOption

func init() {
	computerCmd.AddCommand(computerListCmd)
	computerListOption.SetFlagWithHeaders(computerListCmd, "DisplayName,NumExecutors,Description,Offline")
}

var computerListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("List all Jenkins agents"),
	Long:  i18n.T("List all Jenkins agents"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient, config := GetComputerClient(computerListOption.CommonOption)
		if config == nil {
			err = fmt.Errorf("cannot found the configuration")
			return
		}

		var computers client.ComputerList
		if computers, err = jClient.List(); err == nil {
			computerListOption.Writer = cmd.OutOrStdout()
			computerListOption.CellRenderMap = map[string]common.RenderCell{
				"Offline": func(offline string) string {
					switch offline {
					case "true":
						return "yes"
					}
					return "no"
				},
			}
			err = computerListOption.OutputV2(computers.Computer)
		}
		return
	},
}

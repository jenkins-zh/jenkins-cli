package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-client/pkg/computer"
	cobra_ext "github.com/linuxsuren/cobra-extension/pkg"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerListOption option for config list command
type ComputerListOption struct {
	common.Option
	cobra_ext.OutputOption
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
		jClient, config := GetComputerClient(computerListOption.Option)
		if config == nil {
			err = fmt.Errorf("cannot found the configuration")
			return
		}

		var computers computer.List
		if computers, err = jClient.List(); err == nil {
			computerListOption.Writer = cmd.OutOrStdout()
			computerListOption.CellRenderMap = map[string]cobra_ext.RenderCell{
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

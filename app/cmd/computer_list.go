package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// ComputerListOption option for config list command
type ComputerListOption struct {
	OutputOption

	RoundTripper http.RoundTripper
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
		jClient := &client.ComputerClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: computerListOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		var computers client.ComputerList
		if computers, err = jClient.List(); err == nil {
			computerListOption.Writer = cmd.OutOrStdout()
			computerListOption.CellRenderMap = map[string]RenderCell{
				"Offline": func(offline string) string {
					switch offline {
					case "true":
						return util.ColorWarning("yes")
					}
					return "no"
				},
			}
			err = computerListOption.OutputV2(computers.Computer)
		}
		return
	},
}

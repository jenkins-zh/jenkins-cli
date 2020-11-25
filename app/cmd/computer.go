package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	appCfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(computerCmd)
}

var computerCmd = &cobra.Command{
	Use:     "computer",
	Aliases: []string{"cpu", "agent"},
	Short:   i18n.T("Manage the computers of your Jenkins"),
	Long:    i18n.T(`Manage the computers of your Jenkins`),
}

// GetComputerClient returns the client of computer
func GetComputerClient(option common.Option) (*client.ComputerClient, *appCfg.JenkinsServer) {
	jClient := &client.ComputerClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: option.RoundTripper,
		},
	}
	return jClient, getCurrentJenkinsAndClient(&(jClient.JenkinsCore))
}

// ValidAgentNames autocomplete with agent names
func ValidAgentNames(cmd *cobra.Command, args []string, prefix string) (agentNames []string, directive cobra.ShellCompDirective) {
	directive = cobra.ShellCompDirectiveNoFileComp
	agentNames = make([]string, 0)

	jClient, _ := GetComputerClient(computerListOption.Option)
	if jClient != nil {
		if computers, err := jClient.List(); err == nil {
			for i := range computers.Computer {
				agent := computers.Computer[i]

				// handle it according different cmd
				if (cmd.Use == "start" || cmd.Use == "launch") && !agent.Offline {
					continue
				}

				duplicated := false
				for j := range args {
					if agent.DisplayName == args[j] {
						duplicated = true
						break
					}
				}

				if !duplicated && strings.HasPrefix(agent.DisplayName, prefix) {
					agentNames = append(agentNames, agent.DisplayName)
				}
			}
		}
	}
	return
}

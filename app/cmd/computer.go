package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
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
func GetComputerClient(option CommonOption) (*client.ComputerClient, *JenkinsServer) {
	jClient := &client.ComputerClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: option.RoundTripper,
		},
	}
	return jClient, getCurrentJenkinsAndClient(&(jClient.JenkinsCore))
}

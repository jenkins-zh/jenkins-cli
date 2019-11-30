package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(credentialCmd)
}

var credentialCmd = &cobra.Command{
	Use:     "credential",
	Aliases: []string{"secret", "cred"},
	Short:   i18n.T("Manage the credentials of your Jenkins"),
	Long:    i18n.T(`Manage the credentials of your Jenkins`),
}

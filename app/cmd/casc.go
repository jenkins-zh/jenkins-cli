package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cascCmd)
}

var cascCmd = &cobra.Command{
	Use:   "casc",
	Short: i18n.T("Configuration as Code"),
	Long:  i18n.T("Configuration as Code"),
}

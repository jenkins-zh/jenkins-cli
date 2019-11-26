package cmd

import (
	"os"
	"os/exec"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

func init() {
	pluginCmd.AddCommand(pluginCreateCmd)
}

var pluginCreateCmd = &cobra.Command{
	Use:   "create",
	Short: i18n.T("Create a plugin project from the archetypes"),
	Long: i18n.T(`Create a plugin project from the archetypes
Plugin tutorial is here https://jenkins.io/doc/developer/tutorial/`),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		e := exec.Command("mvn", "-U", "archetype:generate", `-Dfilter="io.jenkins.archetypes:"`)
		e.Stdout = cmd.OutOrStdout()
		e.Stderr = cmd.OutOrStderr()
		e.Stdin = os.Stdin
		err = e.Run()
		return
	},
}

package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

// PluginOptions contains the command line options
type PluginOptions struct {
	Suite string
}

var pluginOpt PluginOptions

func init() {
	rootCmd.AddCommand(pluginCmd)
}

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: i18n.T("Manage the plugins of Jenkins"),
	Long:  i18n.T("Manage the plugins of Jenkins"),
	Example: `  jcli plugin list
  jcli plugin search github
  jcli plugin check`,
}

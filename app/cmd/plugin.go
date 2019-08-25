package cmd

import (
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
	Short: "Manage the plugins of Jenkins",
	Long:  `Manage the plugins of Jenkins`,
	Example: `  jcli plugin list
  jcli plugin search github
  jcli plugin check`,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

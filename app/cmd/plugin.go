package cmd

import (
	"github.com/spf13/cobra"
)

// PluginOptions contains the command line options
type PluginOptions struct {
}

var pluginOpt PluginOptions

func init() {
	rootCmd.AddCommand(pluginCmd)
}

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage the plugins of Jenkins",
	Long:  `Manage the plugins of Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(queueCmd)
}

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Manage the queue of your Jenkins",
	Long:  `Manage the queue of your Jenkins`,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

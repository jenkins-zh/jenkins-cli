package cmd

import (
	"github.com/spf13/cobra"
)

type JobOption struct {
	OutputOption
}

var jobOption JobOption

func init() {
	rootCmd.AddCommand(jobCmd)
	jobCmd.PersistentFlags().StringVarP(&jobOption.Format, "output", "o", "json", "Format the output")
}

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Print the job of your Jenkins",
	Long:  `Print the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

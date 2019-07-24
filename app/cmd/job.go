package cmd

import (
	"github.com/spf13/cobra"
)

type JobOption struct {
	OutputOption

	Name    string
	History bool
}

var jobOption JobOption

func init() {
	rootCmd.AddCommand(jobCmd)
	jobCmd.PersistentFlags().StringVarP(&jobOption.Format, "output", "o", "json", "Format the output")
	jobCmd.PersistentFlags().StringVarP(&jobOption.Name, "name", "n", "", "Name of the job")
}

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Print the job of your Jenkins",
	Long:  `Print the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

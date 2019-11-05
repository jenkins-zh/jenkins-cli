package cmd

import (
	"github.com/spf13/cobra"
)

// JobOption is the job cmd option
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
	Short: "Manage the job of your Jenkins",
	Long: `Manage the job of your Jenkins
Editing the pipeline job needs to install a plugin which is pipeline-restful-api
https://plugins.jenkins.io/pipeline-restful-api`,
}

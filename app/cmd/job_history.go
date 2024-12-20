package cmd

import (
	cobra_ext "github.com/linuxsuren/cobra-extension/pkg"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobHistoryOption is the job history option
type JobHistoryOption struct {
	cobra_ext.OutputOption

	Delete       int
	RoundTripper http.RoundTripper
}

var jobHistoryOption JobHistoryOption

func init() {
	jobCmd.AddCommand(jobHistoryCmd)
	jobHistoryOption.SetFlagWithHeaders(jobHistoryCmd, "DisplayName,Building,Result")
	jobHistoryCmd.Flags().IntVarP(&jobHistoryOption.Delete, "delete", "d", -1, "Delete a history item")
}

var jobHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: i18n.T("Print the history of job in your Jenkins"),
	Long:  i18n.T(`Print the history of job in your Jenkins`),
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jobName := args[0]

		jClient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobHistoryOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jClient.JenkinsCore))

		if jobHistoryOption.Delete != -1 {
			err = jClient.DeleteHistory(jobName, jobHistoryOption.Delete)
			return
		}

		var builds []*client.JobBuild
		builds, err = jClient.GetHistory(jobName)
		if err == nil {
			jobHistoryOption.Writer = cmd.OutOrStdout()
			err = jobHistoryOption.OutputV2(builds)
		}
		return
	},
}

package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobDisableOption is the job delete option
type JobDisableOption struct {
	BatchOption
	CommonOption
}

var jobDisableOption JobDisableOption

func init() {
	jobCmd.AddCommand(jobDisableCmd)
	jobDisableOption.SetFlag(jobDisableCmd)
}

var jobDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: i18n.T("Disable a job in your Jenkins"),
	Long:  i18n.T("Disable a job in your Jenkins"),
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jobName := args[0]
		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobDisableOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		err = jclient.DisableJob(jobName)
		return
	},
}

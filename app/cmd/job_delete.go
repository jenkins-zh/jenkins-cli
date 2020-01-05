package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobDeleteOption is the job delete option
type JobDeleteOption struct {
	BatchOption

	RoundTripper http.RoundTripper
}

var jobDeleteOption JobDeleteOption

func init() {
	jobCmd.AddCommand(jobDeleteCmd)
	jobDeleteOption.SetFlag(jobDeleteCmd)
}

var jobDeleteCmd = &cobra.Command{
	Use:   "delete <jobName>",
	Short: i18n.T("Delete a job in your Jenkins"),
	Long:  i18n.T("Delete a job in your Jenkins"),
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobName := args[0]
		if !jobDeleteOption.Confirm(fmt.Sprintf("Are you sure to delete job %s ?", jobName)) {
			return
		}

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobDeleteOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		err := jclient.Delete(jobName)
		helper.CheckErr(cmd, err)
	},
}

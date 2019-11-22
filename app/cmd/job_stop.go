package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"net/http"
	"strconv"
)

// JobStopOption is the job stop option
type JobStopOption struct {
	BatchOption

	RoundTripper http.RoundTripper
}

var jobStopOption JobStopOption

func init() {
	jobCmd.AddCommand(jobStopCmd)
	jobStopOption.SetFlag(jobStopCmd)
}

var jobStopCmd = &cobra.Command{
	Use:   "stop <jobName> [buildNumber]",
	Short: i18n.T("Stop a job build in your Jenkins"),
	Long:  i18n.T("Stop a job build in your Jenkins"),
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		buildNum := -1
		if len(args) > 1 {
			if buildNum, err = strconv.Atoi(args[1]); err != nil {
				return
			}
		}

		jobName := args[0]
		if !jobStopOption.Confirm(fmt.Sprintf("Are you sure to stop job %s ?", jobName)) {
			return
		}

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobStopOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		return jclient.StopJob(jobName, buildNum)
	},
}

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
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
	Use:   "stop <jobName> <buildNumber>",
	Short: "Stop a job build in your Jenkins",
	Long:  `Stop a job build in your Jenkins`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			buildNum int
			err      error
		)
		if buildNum, err = strconv.Atoi(args[1]); err != nil {
			cmd.PrintErrln(err)
			return
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

		if err := jclient.StopJob(jobName, buildNum); err != nil {
			log.Fatal(err)
		}
	},
}

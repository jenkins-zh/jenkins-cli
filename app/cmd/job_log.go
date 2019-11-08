package cmd

import (
	"net/http"
	"time"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobLogOption is the job log option
type JobLogOption struct {
	WatchOption
	History int

	LogText      string
	LastBuildID  int
	LastBuildURL string

	RoundTripper http.RoundTripper
}

var jobLogOption JobLogOption

func init() {
	jobCmd.AddCommand(jobLogCmd)
	jobLogCmd.Flags().IntVarP(&jobLogOption.History, "history", "s", -1, "Specific build history of log")
	jobLogCmd.Flags().BoolVarP(&jobLogOption.Watch, "watch", "w", false, "Watch the job logs")
	jobLogCmd.Flags().IntVarP(&jobLogOption.Interval, "interval", "i", 1, "Interval of watch")
}

var jobLogCmd = &cobra.Command{
	Use:   "log <jobName> [buildID]",
	Short: "Print the job's log of your Jenkins",
	Long: `Print the job's log of your Jenkins
It'll print the log text of the last build if you don't give the build id.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobLogOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		lastBuildID := -1
		var jobBuild *client.JobBuild
		var err error
		for {
			if jobBuild, err = jclient.GetBuild(name, -1); err == nil {
				jobLogOption.LastBuildID = jobBuild.Number
				jobLogOption.LastBuildURL = jobBuild.URL
			}
			if lastBuildID != jobLogOption.LastBuildID {
				lastBuildID = jobLogOption.LastBuildID
				cmd.Println("Current build number:", jobLogOption.LastBuildID)
				cmd.Println("Current build url:", jobLogOption.LastBuildURL)

				err = printLog(jclient, cmd, name, jobLogOption.History, 0)
			}

			if err != nil || !jobLogOption.Watch {
				break
			}

			time.Sleep(time.Duration(jobLogOption.Interval) * time.Second)
		}
	},
}

func printLog(jclient *client.JobClient, cmd *cobra.Command, jobName string, history int, start int64) (err error) {
	var jobLog client.JobLog
	if jobLog, err = jclient.Log(jobName, history, start); err == nil {
		isNew := false

		if jobLogOption.LogText != jobLog.Text {
			jobLogOption.LogText = jobLog.Text
			isNew = true
		} else if history == -1 {
			if build, err := jclient.GetBuild(jobName, -1); err == nil && jobLogOption.LastBuildID != build.Number {
				jobLogOption.LastBuildID = build.Number
				jobLogOption.LastBuildURL = build.URL
				isNew = true
			}
		}

		if isNew {
			cmd.Print(jobLog.Text)
		}

		if jobLog.HasMore {
			err = printLog(jclient, cmd, jobName, history, jobLog.NextStart)
		}
	}
	return
}

package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobLogOption is the job log option
type JobLogOption struct {
	common.WatchOption
	History       int
	LogText       string
	LastBuildID   int
	LastBuildURL  string
	NumberOfLines int

	RoundTripper http.RoundTripper
}

var jobLogOption JobLogOption

func init() {
	jobCmd.AddCommand(jobLogCmd)
	jobLogCmd.Flags().IntVarP(&jobLogOption.History, "history", "s", -1,
		i18n.T("Specific build history of log"))
	jobLogCmd.Flags().BoolVarP(&jobLogOption.Watch, "watch", "w", false,
		i18n.T("Watch the job logs"))
	jobLogCmd.Flags().IntVarP(&jobLogOption.Interval, "interval", "i", 1,
		i18n.T("Interval of watch"))
	jobLogCmd.Flags().IntVarP(&jobLogOption.NumberOfLines, "tail", "t", -1,
		i18n.T("the last number of lines of the log"))
}

var jobLogCmd = &cobra.Command{
	Use:   "log",
	Short: i18n.T("Print the job's log of your Jenkins"),
	Long: i18n.T(`Print the job's log of your Jenkins
It'll print the log text of the last build if you don't give the build id.`),
	Args: cobra.MinimumNArgs(1),
	Example: `jcli job log <jobName> [buildID]
jcli job log <jobName> --history 1
jcli job log <jobName> --watch
jcli job log <jobName> --tail <numberOfLines>`,
	PreRunE: func(_ *cobra.Command, args []string) (err error) {
		if len(args) >= 3 && jobLogOption.History == -1 {
			var history int
			historyStr := args[1]
			if history, err = strconv.Atoi(historyStr); err == nil {
				jobLogOption.History = history
			} else {
				err = fmt.Errorf("job history must be a number instead of '%s'", historyStr)
			}
			numberOfLinesStr := args[2]
			numberOfLines, err := strconv.Atoi(numberOfLinesStr)
			if err != nil || numberOfLines <= 0 {
				err = fmt.Errorf(err.Error(), "lines of job must be a positive integer instead of '%s'", numberOfLinesStr)
			} else {
				jobLogOption.NumberOfLines = numberOfLines
			}
		}
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobLogOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		lastBuildID := -1
		var jobBuild *client.JobBuild
		var err error
		for {
			if jobBuild, err = jclient.GetBuild(name, jobLogOption.History); err == nil {
				jobLogOption.LastBuildID = jobBuild.Number
				jobLogOption.LastBuildURL = jobBuild.URL
			} else {
				break
			}

			if lastBuildID != jobLogOption.LastBuildID {
				lastBuildID = jobLogOption.LastBuildID
				cmd.Println("Current build number:", jobLogOption.LastBuildID)
				cmd.Println("Current build url:", jobLogOption.LastBuildURL)

				err = printLog(jclient, cmd, name, jobLogOption.History, 0, jobLogOption.NumberOfLines)
			}

			if err != nil || !jobLogOption.Watch {
				break
			}

			time.Sleep(time.Duration(jobLogOption.Interval) * time.Second)
		}
	},
}

func printLog(jclient *client.JobClient, cmd *cobra.Command, jobName string, history int, start int64, numberOfLines int) (err error) {
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

		numberOfLinesOfJobLogText := strings.Count(jobLog.Text, "\n")
		if isNew && (numberOfLines > 0 || numberOfLines == -1) {
			if numberOfLines >= numberOfLinesOfJobLogText || numberOfLines == -1 {
				cmd.Print(jobLog.Text)
				numberOfLines -= numberOfLinesOfJobLogText

			} else if numberOfLines < numberOfLinesOfJobLogText {
				text := jobLog.Text
				for i := 0; i <= numberOfLinesOfJobLogText-numberOfLines; i++ {
					temp := strings.Index(text, "\n")
					text = text[temp+1:]
				}
				cmd.Print(text)
				numberOfLines = 0
			}
		}

		if jobLog.HasMore {
			err = printLog(jclient, cmd, jobName, history, jobLog.NextStart, numberOfLines)
		}
	}
	return
}

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobInputOption is the job delete option
type JobInputOption struct {
	BatchOption
	RoundTripper http.RoundTripper
}

var jobInputOption JobInputOption

func init() {
	jobCmd.AddCommand(jobInputCmd)
}

var jobInputCmd = &cobra.Command{
	Use:   "input <jobName> [buildID]",
	Short: "Input a job in your Jenkins",
	Long:  `Input a job in your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		jobName := args[0]
		buildID := -1

		if len(args) >= 2 {
			var err error
			if buildID, err = strconv.Atoi(args[1]); err != nil {
				cmd.PrintErrln(err)
			}
		}

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobInputOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if actions, err := jclient.GetJobInputActions(jobName, buildID); err != nil {
			log.Fatal(err)
		} else if len(actions) >= 1 {
			action := ""
			prompt := &survey.Input{
				Message: fmt.Sprintf("Are you going to process or abort this input: %s?", actions[0].Message),
			}
			survey.AskOne(prompt, &action)

			fmt.Println(actions[0])
			if action == "process" {
				err = jclient.JobInputSubmitTest(jobName, actions[0].ID, buildID, false, nil)
			} else if action == "abort" {
				err = jclient.JobInputSubmitTest(jobName, actions[0].ID, buildID, true, nil)
			} else {
				cmd.PrintErrln("Only process or abort is accepted!")
			}

			if err != nil {
				cmd.PrintErrln(err)
			}
		}
	},
}

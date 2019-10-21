package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobInputOption is the job delete option
type JobInputOption struct {
	BatchOption

	Action string

	RoundTripper http.RoundTripper
	Stdio        terminal.Stdio
}

var jobInputOption JobInputOption

func init() {
	jobCmd.AddCommand(jobInputCmd)
	jobInputCmd.Flags().StringVarP(&jobInputOption.Action, "action", "", "", "The action wether you want to process or abort.")
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

		if inputActions, err := jclient.GetJobInputActions(jobName, buildID); err != nil {
			log.Fatal(err)
		} else if len(inputActions) >= 1 {
			inputAction := inputActions[0]
			params := make(map[string]string)

			if len(inputAction.Inputs) > 0 {
				inputsJSON, _ := json.MarshalIndent(inputAction.Inputs, "", " ")
				content := string(inputsJSON)

				prompt := &survey.Editor{
					Message:       "Edit your pipeline input parameters",
					FileName:      "*.json",
					Default:       content,
					HideDefault:   true,
					AppendDefault: true,
				}

				if err = survey.AskOne(prompt, &content); err != nil {
					log.Fatal(err)
				}

				if err = json.Unmarshal([]byte(content), &(inputAction.Inputs)); err != nil {
					log.Fatal(err)
				}

				for _, input := range inputAction.Inputs {
					params[input.Name] = input.Value
				}
			}

			render := &survey.Renderer{}
			render.WithStdio(jobInputOption.Stdio)

			// allow users make their choice through cli arguments
			action := jobInputOption.Action
			if action == "" {
				prompt := &survey.Input{
					Renderer: *render,
					Message:  fmt.Sprintf("Are you going to process or abort this input: %s?", inputAction.Message),
				}
				survey.AskOne(prompt, &action)
			}

			if action == "process" {
				err = jclient.JobInputSubmit(jobName, inputAction.ID, buildID, false, params)
			} else if action == "abort" {
				err = jclient.JobInputSubmit(jobName, inputAction.ID, buildID, true, params)
			} else {
				cmd.PrintErrln("Only process or abort is accepted!")
			}

			if err != nil {
				cmd.PrintErrln(err)
			}
		}
	},
}

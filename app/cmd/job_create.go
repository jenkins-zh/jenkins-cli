package cmd

import (
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobCreateOption is the job create option
type JobCreateOption struct {
	Copy string
	Type string

	RoundTripper http.RoundTripper
}

var jobCreateOption JobCreateOption

func init() {
	jobCmd.AddCommand(jobCreateCmd)
	jobCreateCmd.Flags().StringVarP(&jobCreateOption.Copy, "copy", "", "", "Copy an exists job")
	jobCreateCmd.Flags().StringVarP(&jobCreateOption.Type, "type", "", "", "Which type do you want to create")
}

var jobCreateCmd = &cobra.Command{
	Use:   "create <jobName>",
	Short: "Create a job in your Jenkins",
	Long:  `Create a job in your Jenkins`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobName := args[0]
		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobCreateOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		var createMode string
		var err error
		if createMode, err = jobCreateOption.getCreateMode(jclient); err == nil {
			payload := client.CreateJobPayload{
				Name: jobName,
				Mode: createMode,
				From: jobCreateOption.Copy,
			}

			if jobCreateOption.Copy != "" {
				payload.Mode = "copy"
			}
			err = jclient.Create(payload)
		}
		helper.CheckErr(cmd, err)
	},
}

func (j *JobCreateOption) getCreateMode(jclient *client.JobClient) (mode string, err error) {
	mode = j.Type
	if j.Copy != "" || mode != "" {
		return
	}

	var types []string
	var typeMap map[string]string
	typeMap, types, err = GetCategories(jclient)
	if err != nil {
		return
	}

	var jobType string
	prompt := &survey.Select{
		Message: "Choose a job type:",
		Options: types,
		Default: jobType,
	}
	if err = survey.AskOne(prompt, &jobType); err != nil {
		return
	}

	mode = typeMap[jobType]
	return
}

package cmd

import (
	"net/http"

	"github.com/AlecAivazis/survey/v2"
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
		if createMode, err = jobCreateOption.getCreateMode(jclient); err != nil {
			cmd.PrintErrln(err)
			return
		}

		payload := client.CreateJobPayload{
			Name: jobName,
			Mode: createMode,
			From: jobCreateOption.Copy,
		}

		if jobCreateOption.Copy != "" {
			payload.Mode = "copy"
		}
		if err := jclient.Create(payload); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func (j *JobCreateOption) getCreateMode(jclient *client.JobClient) (mode string, err error) {
	mode = j.Type
	if j.Copy != "" || mode != "" {
		return
	}

	types := make(map[string]string)
	var categories []client.JobCategory
	if categories, err = jclient.GetJobTypeCategories(); err == nil {
		for _, category := range categories {
			for _, item := range category.Items {
				types[item.DisplayName] = item.Class
			}
		}
	} else {
		return
	}

	typesArray := make([]string, 0)
	for tp := range types {
		typesArray = append(typesArray, tp)
	}

	var jobType string
	prompt := &survey.Select{
		Message: "Choose a job type:",
		Options: typesArray,
		Default: jobType,
	}
	if err = survey.AskOne(prompt, &jobType); err != nil {
		return
	}

	mode = types[jobType]
	return
}

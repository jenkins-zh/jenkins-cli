package cmd

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobCreateOption struct {
}

var jobCreateOption JobCreateOption

func init() {
	jobCmd.AddCommand(jobCreateCmd)
}

var jobCreateCmd = &cobra.Command{
	Use:   "create <jobName>",
	Short: "Create a job in your Jenkins",
	Long:  `Create a job in your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		jobName := args[0]

		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		types := make(map[string]string)
		if categories, err := jclient.GetJobTypeCategories(); err == nil {
			for _, category := range categories {
				for _, item := range category.Items {
					types[item.DisplayName] = item.Class
				}
			}
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
		survey.AskOne(prompt, &jobType)

		if err := jclient.Create(jobName, types[jobType]); err != nil {
			log.Fatal(err)
		}
	},
}

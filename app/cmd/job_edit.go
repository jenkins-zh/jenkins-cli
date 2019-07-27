package cmd

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobEditCmd)
}

var jobEditCmd = &cobra.Command{
	Use:   "edit <jobName>",
	Short: "Edit the job of your Jenkins",
	Long:  `Edit the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		name := args[0]
		var content string
		var err error
		if content, err = getPipeline(name); err != nil {
			log.Fatal(err)
		}

		prompt := &survey.Editor{
			Message:       "Edit your pipeline script",
			FileName:      "*.sh",
			Default:       content,
			HideDefault:   true,
			AppendDefault: true,
		}

		if err = survey.AskOne(prompt, &content); err != nil {
			log.Fatal(err)
		}

		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth
		if err = jclient.UpdatePipeline(name, content); err != nil {
			fmt.Println("update failed")
			log.Fatal(err)
		}
	},
}

func getPipeline(name string) (script string, err error) {
	jenkins := getCurrentJenkins()
	jclient := &client.JobClient{}
	jclient.URL = jenkins.URL
	jclient.UserName = jenkins.UserName
	jclient.Token = jenkins.Token
	jclient.Proxy = jenkins.Proxy
	jclient.ProxyAuth = jenkins.ProxyAuth

	var job *client.Pipeline
	if job, err = jclient.GetPipeline(name); err == nil {
		script = job.Script
	}
	return
}

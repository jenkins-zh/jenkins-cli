package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/AlecAivazis/survey"
	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobEditCmd)
}

var jobEditCmd = &cobra.Command{
	Use:   "edit -n",
	Short: "Edit the job of your Jenkins",
	Long:  `Edit the job of your Jenkins`,
	Args: func(cmd *cobra.Command, args []string) error {
		if jobOption.Name == "" {
			return errors.New("requires job name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var content string
		var err error
		if content, err = getPipeline(jobOption.Name); err != nil {
			log.Fatal(err)
		}

		prompt := &survey.Editor{
			Message:       "Edit your pipeline script",
			FileName:      "*.sh",
			Default:       content,
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
		if err = jclient.UpdatePipeline(jobOption.Name, content); err != nil {
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

package cmd

import (
	"fmt"
	"log"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(centerCmd)
}

var centerCmd = &cobra.Command{
	Use:   "center",
	Short: "Manage your update center",
	Long:  `Manage your update center`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		printJenkinsStatus(jenkins)
		printUpdateCenter(jenkins)
	},
}

func printUpdateCenter(jenkins *JenkinsServer) {
	jclient := &client.UpdateCenterManager{}
	jclient.URL = jenkins.URL
	jclient.UserName = jenkins.UserName
	jclient.Token = jenkins.Token

	if status, err := jclient.Status(); err == nil {
		fmt.Println("RestartRequiredForCompletion:", status.RestartRequiredForCompletion)
		if status.Jobs != nil {
			for i, job := range status.Jobs {
				fmt.Printf("%d, %s, %s\n", i, job.Type, job.ErrorMessage)
			}
		}
	} else {
		log.Fatal(err)
	}
}

func printJenkinsStatus(jenkins *JenkinsServer) {
	jclient := &client.JenkinsStatusClient{}
	jclient.URL = jenkins.URL
	jclient.UserName = jenkins.UserName
	jclient.Token = jenkins.Token

	if status, err := jclient.Get(); err == nil {
		fmt.Println("Jenkins Version:", status.Version)
	} else {
		log.Fatal(err)
	}
}

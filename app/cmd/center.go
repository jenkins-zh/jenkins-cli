package cmd

import (
	"fmt"
	"log"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type CenterOption struct {
	WatchOption

	CeneterStatus string
}

var centerOption CenterOption

func init() {
	rootCmd.AddCommand(centerCmd)
}

var centerCmd = &cobra.Command{
	Use:   "center",
	Short: "Manage your update center",
	Long:  `Manage your update center`,
	Run: func(_ *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		printJenkinsStatus(jenkins)

		printUpdateCenter(jenkins)
	},
}

func printUpdateCenter(jenkins *JenkinsServer) {
	jclient := &client.UpdateCenterManager{}
	jclient.URL = jenkins.URL
	jclient.UserName = jenkins.UserName
	jclient.Token = jenkins.Token
	jclient.Proxy = jenkins.Proxy
	jclient.ProxyAuth = jenkins.ProxyAuth

	var centerStatus string
	if status, err := jclient.Status(); err == nil {
		centerStatus += fmt.Sprintf("RestartRequiredForCompletion: %v\n", status.RestartRequiredForCompletion)
		if status.Jobs != nil {
			for i, job := range status.Jobs {
				if job.Type == "InstallationJob" {
					centerStatus += fmt.Sprintf("%d, %s, %s, %v, %s\n", i, job.Type, job.Name, job.Status, job.ErrorMessage)
				} else {
					centerStatus += fmt.Sprintf("%d, %s, %s\n", i, job.Type, job.ErrorMessage)
				}
			}
		}

		if centerOption.CeneterStatus != centerStatus {
			centerOption.CeneterStatus = centerStatus

			fmt.Printf("%s", centerStatus)
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
	jclient.Proxy = jenkins.Proxy
	jclient.ProxyAuth = jenkins.ProxyAuth

	if status, err := jclient.Get(); err == nil {
		fmt.Println("Jenkins Version:", status.Version)
	} else {
		log.Fatal(err)
	}
}

package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterOption is the center cmd option
type CenterOption struct {
	WatchOption

	RoundTripper  http.RoundTripper
	CeneterStatus string
}

var centerOption CenterOption

func init() {
	rootCmd.AddCommand(centerCmd)
}

var centerCmd = &cobra.Command{
	Use:   "center",
	Short: i18n.T("Manage your update center"),
	Long:  `Manage your update center`,
	Run: func(cmd *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		printJenkinsStatus(jenkins, cmd, centerOption.RoundTripper)

		printUpdateCenter(jenkins, cmd, centerOption.RoundTripper)
	},
}

func printUpdateCenter(jenkins *JenkinsServer, cmd *cobra.Command, roundTripper http.RoundTripper) (
	status *client.UpdateCenter, err error) {
	jclient := &client.UpdateCenterManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: roundTripper,
		},
	}
	getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

	var centerStatus string
	if status, err = jclient.Status(); err == nil {
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

			cmd.Printf("%s", centerStatus)
		}
	}
	return
}

func printJenkinsStatus(jenkins *JenkinsServer, cmd *cobra.Command, roundTripper http.RoundTripper) {
	jclient := &client.JenkinsStatusClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: roundTripper,
		},
	}
	jclient.URL = jenkins.URL
	jclient.UserName = jenkins.UserName
	jclient.Token = jenkins.Token
	jclient.Proxy = jenkins.Proxy
	jclient.ProxyAuth = jenkins.ProxyAuth

	status, err := jclient.Get()
	if err == nil {
		cmd.Println("Jenkins Version:", status.Version)
	}
	helper.CheckErr(cmd, err)
}

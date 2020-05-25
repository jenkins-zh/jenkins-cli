package center

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"net/http"

	. "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func NewCenterCmd(client common.JenkinsClient, jenkinsConfigMgr common.JenkinsConfigMgr) (cmd *cobra.Command) {
	opt := &CenterOption{
		JenkinsClient:    client,
		JenkinsConfigMgr: jenkinsConfigMgr,
	}

	cmd = &cobra.Command{
		Use:   "center",
		Short: i18n.T("Manage your update center"),
		Long:  i18n.T("Manage your update center"),
		RunE:  opt.RunE,
	}

	cmd.AddCommand(NewCenterStartcmd(client),
		NewCenterUpgradecmd(client),
		NewCenterWatchCmd(client, opt),
		NewCenterMirrorCmd(client),
		NewCenterIdentityCmd(client))
	return
}

func (o *CenterOption) RunE(cmd *cobra.Command, _ []string) (err error) {
	jenkins := o.JenkinsClient.GetCurrentJenkinsFromOptions()
	o.printJenkinsStatus(jenkins, cmd, o.RoundTripper)

	_, err = o.printUpdateCenter(jenkins, cmd, o.RoundTripper)
	return
}

func (o *CenterOption) printUpdateCenter(jenkins *JenkinsServer, cmd *cobra.Command, roundTripper http.RoundTripper) (
	status *client.UpdateCenter, err error) {
	jclient := &client.UpdateCenterManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: roundTripper,
		},
	}
	o.JenkinsClient.GetCurrentJenkinsAndClient(&(jclient.JenkinsCore))

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

		if o.CenterStatus != centerStatus {
			o.CenterStatus = centerStatus

			cmd.Printf("%s", centerStatus)
		}
	}
	return
}

func (o *CenterOption) printJenkinsStatus(jenkins *JenkinsServer, cmd *cobra.Command, roundTripper http.RoundTripper) {
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

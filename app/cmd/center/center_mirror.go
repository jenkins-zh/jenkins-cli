package center

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func NewCenterMirrorCmd(client common.JenkinsClient) (cmd *cobra.Command) {
	opt := &CenterMirrorOption{
		JenkinsClient: client,
	}
	cmd = &cobra.Command{
		Use:   "mirror",
		Short: i18n.T("Set the update center to a mirror address"),
		Long:  i18n.T("Set the update center to a mirror address"),
		RunE:  opt.RunE,
	}

	cmd.Flags().BoolVarP(&opt.Enable, "enable", "", true,
		i18n.T("If you want to enable update center server"))
	cmd.Flags().StringVarP(&opt.MirrorURL, "mirror-url", "", "https://updates.jenkins-zh.cn/update-center.json",
		i18n.T("The address of update center site mirror"))
	return
}

func (o *CenterMirrorOption) RunE(cmd *cobra.Command, _ []string) (err error) {
	jclient := &client.UpdateCenterManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: o.RoundTripper,
		},
	}
	o.JenkinsClient.GetCurrentJenkinsAndClient(&(jclient.JenkinsCore))

	var siteURL string
	if o.Enable {
		siteURL = o.MirrorURL
	} else {
		siteURL = "https://updates.jenkins.io/update-center.json"
	}

	if err = jclient.ChangeUpdateCenterSite("default", siteURL); err == nil {
		err = jclient.SetMirrorCertificate(o.Enable)
	}
	return
}

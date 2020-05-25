package center

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// NewCenterUpgradecmd creates the center upgrade command
func NewCenterUpgradecmd(jenkinsClient common.JenkinsClient) (cmd *cobra.Command) {
	opt := &CenterUpgradeOption{
		JenkinsClient: jenkinsClient,
	}
	cmd = &cobra.Command{
		Use:   "upgrade",
		Short: i18n.T("Upgrade your Jenkins"),
		Long:  i18n.T("Upgrade your Jenkins"),
		RunE:  opt.RunE,
	}
	return
}

// RunE is the main entry point
func (o *CenterUpgradeOption) RunE(cmd *cobra.Command, _ []string) (err error) {
	jclient := &client.UpdateCenterManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: o.RoundTripper,
		},
	}
	o.JenkinsClient.GetCurrentJenkinsAndClient(&(jclient.JenkinsCore))
	return jclient.Upgrade()
}

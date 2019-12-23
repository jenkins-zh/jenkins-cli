package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterUpgradeOption option for upgrade Jenkins
type CenterUpgradeOption struct {
	RoundTripper http.RoundTripper
}

var centerUpgradeOption CenterUpgradeOption

func init() {
	centerCmd.AddCommand(centerUpgradeCmd)
}

var centerUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: i18n.T("Upgrade your Jenkins"),
	Long:  i18n.T("Upgrade your Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jclient := &client.UpdateCenterManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerUpgradeOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))
		return jclient.Upgrade()
	},
}

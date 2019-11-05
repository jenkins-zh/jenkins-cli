package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

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
	Short: "Upgrade your Jenkins",
	Long:  `Upgrade your Jenkins`,
	Run: func(cmd *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()

		jclient := &client.UpdateCenterManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerUpgradeOption.RoundTripper,
			},
		}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		err := jclient.Upgrade()
		helper.CheckErr(cmd, err)
	},
}

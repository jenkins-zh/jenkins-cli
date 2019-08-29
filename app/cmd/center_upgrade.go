package cmd

import (
	"log"
	"net/http"

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
	Run: func(_ *cobra.Command, _ []string) {
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

		if err := jclient.Upgrade(); err != nil {
			log.Fatal(err)
		}
	},
}

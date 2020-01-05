package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterMirrorOption option for upgrade Jenkins
type CenterMirrorOption struct {
	RoundTripper http.RoundTripper

	Enable    bool
	MirrorURL string
}

var centerMirrorOption CenterMirrorOption

func init() {
	centerCmd.AddCommand(centerMirrorCmd)
	centerMirrorCmd.Flags().BoolVarP(&centerMirrorOption.Enable, "enable", "", true,
		i18n.T("If you want to enable update center server"))
	centerMirrorCmd.Flags().StringVarP(&centerMirrorOption.MirrorURL, "mirror-url", "", "https://updates.jenkins-zh.cn/update-center.json",
		i18n.T("The address of update center site mirror"))
}

var centerMirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: i18n.T("Set the update center to a mirror address"),
	Long:  i18n.T("Set the update center to a mirror address"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jclient := &client.UpdateCenterManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerUpgradeOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		var siteURL string
		if centerMirrorOption.Enable {
			siteURL = centerMirrorOption.MirrorURL
		} else {
			siteURL = "https://updates.jenkins.io/update-center.json"
		}

		if err = jclient.ChangeUpdateCenterSite("default", siteURL); err == nil {
			err = jclient.SetMirrorCertificate(centerMirrorOption.Enable)
		}
		return
	},
}

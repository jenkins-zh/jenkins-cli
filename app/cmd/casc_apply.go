package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// CASCApplyOption as the options of apply configuration as code
type CASCApplyOption struct {
	RoundTripper http.RoundTripper
}

var cascApplyOption CASCApplyOption

func init() {
	cascCmd.AddCommand(cascApplyCmd)
}

var cascApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: i18n.T("Apply config through configuration-as-code"),
	Long:  i18n.T("Apply config through configuration-as-code"),
	RunE: func(cmd *cobra.Command, _ []string) error {
		jClient := &client.CASCManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: cascApplyOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jClient.JenkinsCore))
		return jClient.Apply()
	},
}

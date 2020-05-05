package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// CASCExportOption as the options of reload configuration as code
type CASCExportOption struct {
	RoundTripper http.RoundTripper
}

var cascExportOption CASCExportOption

func init() {
	cascCmd.AddCommand(cascExportCmd)
}

var cascExportCmd = &cobra.Command{
	Use:   "export",
	Short: i18n.T("Export the config from configuration-as-code"),
	Long:  i18n.T("Export the config from configuration-as-code"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.CASCManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: cascExportOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jClient.JenkinsCore))

		var config string
		if config, err = jClient.Export(); err == nil {
			cmd.Print(config)
		}
		return
	},
	Annotations: map[string]string{
		common.Since: common.VersionSince0024,
	},
}

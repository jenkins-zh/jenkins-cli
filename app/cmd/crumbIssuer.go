package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/spf13/cobra"
)

// CrumbIssuerOptions contains the command line options
type CrumbIssuerOptions struct {
	RoundTripper http.RoundTripper
}

func init() {
	rootCmd.AddCommand(crumbIssuerCmd)
}

var crumbIssuerOptions CrumbIssuerOptions

var crumbIssuerCmd = &cobra.Command{
	Use:   "crumb",
	Short: "Print crumbIssuer of Jenkins",
	Long:  `Print crumbIssuer of Jenkins`,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jenkinsCore := &client.JenkinsCore{RoundTripper: crumbIssuerOptions.RoundTripper}
		getCurrentJenkinsAndClient(jenkinsCore)

		var crumb *client.JenkinsCrumb
		if crumb, err = jenkinsCore.GetCrumb(); err == nil && crumb != nil {
			cmd.Printf("%s=%s\n", crumb.CrumbRequestField, crumb.Crumb)
		}
		return
	},
}

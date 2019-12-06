package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

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
	Run: func(cmd *cobra.Command, _ []string) {
		jenkinsCore := &client.JenkinsCore{RoundTripper: crumbIssuerOptions.RoundTripper}
		getCurrentJenkinsAndClientOrDie(jenkinsCore)

		crumb, err := jenkinsCore.GetCrumb()
		if err == nil {
			cmd.Printf("%s=%s\n", crumb.CrumbRequestField, crumb.Crumb)
		}
		helper.CheckErr(cmd, err)
	},
}

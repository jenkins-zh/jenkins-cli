package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// CASCSchemaOption as the options of reload configuration as code
type CASCSchemaOption struct {
	RoundTripper http.RoundTripper
}

var cascSchemaOption CASCSchemaOption

func init() {
	cascCmd.AddCommand(cascSchemaCmd)
}

var cascSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: i18n.T("Get the schema of configuration-as-code"),
	Long:  i18n.T("Get the schema of configuration-as-code"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.CASCManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: cascSchemaOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jClient.JenkinsCore))

		var config string
		if config, err = jClient.Schema(); err == nil {
			cmd.Print(config)
		}
		return
	},
}

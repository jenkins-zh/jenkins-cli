package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// CredentialListOption option for credential list command
type CredentialListOption struct {
	common.OutputOption

	Store string

	RoundTripper http.RoundTripper
}

var credentialListOption CredentialListOption

func init() {
	credentialCmd.AddCommand(credentialListCmd)
	credentialListCmd.Flags().StringVarP(&credentialListOption.Store, "store", "", "system",
		i18n.T("The store name of Jenkins credentials"))
	credentialListOption.SetFlagWithHeaders(credentialListCmd, "DisplayName,ID,TypeName,Description")
}

var credentialListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("List all credentials of Jenkins"),
	Long:  i18n.T("List all credentials of Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.CredentialsManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: credentialListOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		var credentialList client.CredentialList
		if credentialList, err = jClient.GetList(credentialListOption.Store); err == nil {
			credentialListOption.Writer = cmd.OutOrStdout()
			err = credentialListOption.OutputV2(credentialList.Credentials)
		}
		return
	},
	Annotations: map[string]string{
		common.Since: "v0.0.24",
	},
}

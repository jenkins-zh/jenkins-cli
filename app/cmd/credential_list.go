package cmd

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/client"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// CredentialListOption option for credential list command
type CredentialListOption struct {
	OutputOption

	Store string

	RoundTripper http.RoundTripper
}

var credentialListOption CredentialListOption

func init() {
	credentialCmd.AddCommand(credentialListCmd)
	credentialListCmd.Flags().StringVarP(&credentialListOption.Store, "store", "", "system",
		i18n.T("The store name of Jenkins credentials"))
	credentialListOption.SetFlag(credentialListCmd)
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
		var data []byte
		if credentialList, err = jClient.GetList(credentialListOption.Store); err == nil {
			if data, err = credentialListOption.Output(credentialList); err == nil {
				cmd.Print(string(data))
			}
		}
		return
	},
}

// Output render data into byte array as a table format
func (o *CredentialListOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil && o.Format == TableOutputFormat {
		credentialList := obj.(client.CredentialList)

		buf := new(bytes.Buffer)
		table := util.CreateTableWithHeader(buf, o.WithoutHeaders)
		table.AddHeader("number", "displayName", "id", "type", "description")
		for i, cred := range credentialList.Credentials {
			table.AddRow(fmt.Sprintf("%d", i), cred.DisplayName, cred.ID, cred.TypeName, cred.Description)
		}
		table.Render()
		err = nil
		data = buf.Bytes()
	}
	return
}

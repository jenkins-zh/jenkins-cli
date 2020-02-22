package cmd

import (
	"fmt"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CredentialCreateOption option for credential delete command
type CredentialCreateOption struct {
	Description string
	ID          string
	Store       string

	Username string
	Password string

	Secret string

	Scope string
	Type  string

	RoundTripper http.RoundTripper
}

var credentialCreateOption CredentialCreateOption

func init() {
	credentialCmd.AddCommand(credentialCreateCmd)
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Store, "store", "", "system",
		i18n.T("The store name of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Type, "type", "", "basic",
		i18n.T("The type of Jenkins credentials which could be: basic, secret"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Scope, "scope", "", "GLOBAL",
		i18n.T("The scope of Jenkins credentials which might be GLOBAL or SYSTEM"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.ID, "credential-id", "", "",
		i18n.T("The ID of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Username, "credential-username", "", "",
		i18n.T("The Username of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Password, "credential-password", "", "",
		i18n.T("The Password of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Description, "desc", "", "",
		i18n.T("The Description of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Secret, "secret", "", "",
		i18n.T("The Secret of Jenkins credentials"))
}

var credentialCreateCmd = &cobra.Command{
	Use:   "create",
	Short: i18n.T("Create a credential from Jenkins"),
	Long:  i18n.T("Create a credential from Jenkins"),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) >= 1 {
			credentialCreateOption.Store = args[0]
		}

		if credentialCreateOption.Store == "" {
			err = fmt.Errorf("the store or id of target credential is empty")
		}
		return
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jClient := &client.CredentialsManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: credentialCreateOption.RoundTripper,
				Debug:        true,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		switch credentialCreateOption.Type {
		case "basic":
			err = jClient.CreateUsernamePassword(credentialCreateOption.Store, client.UsernamePasswordCredential{
				Username: credentialCreateOption.Username,
				Password: credentialCreateOption.Password,
				Credential: client.Credential{
					Scope:       credentialCreateOption.Scope,
					ID:          credentialCreateOption.ID,
					Description: credentialCreateOption.Description,
				},
			})
		case "secret":
			err = jClient.CreateSecret(credentialCreateOption.Store, client.StringCredentials{
				Secret: credentialCreateOption.Secret,
				Credential: client.Credential{
					Scope:       credentialCreateOption.Scope,
					ID:          credentialCreateOption.ID,
					Description: credentialCreateOption.Description,
				},
			})
		default:
			err = fmt.Errorf("unknow credential type: %s", credentialCreateOption.Type)
		}
		return
	},
	Annotations: map[string]string{
		since: "v0.0.24",
	},
}

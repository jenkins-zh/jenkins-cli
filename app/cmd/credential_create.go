package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"net/http"
)

// CredentialCreateOption option for credential delete command
type CredentialCreateOption struct {
	Description string
	ID          string
	Store       string

	Username string
	Password string

	Secret string

	Type string

	RoundTripper http.RoundTripper
}

var credentialCreateOption CredentialCreateOption

func init() {
	credentialCmd.AddCommand(credentialCreateCmd)
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Store, "store", "", "system",
		i18n.T("The store name of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Type, "type", "", "basic",
		i18n.T("The type of Jenkins credentials which could be: basic, secret"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.ID, "id", "", "",
		i18n.T("The ID of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Username, "username", "", "",
		i18n.T("The Username of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Password, "password", "", "",
		i18n.T("The Password of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Description, "desc", "", "",
		i18n.T("The Description of Jenkins credentials"))
	credentialCreateCmd.Flags().StringVarP(&credentialCreateOption.Secret, "secret", "", "",
		i18n.T("The Secret of Jenkins credentials"))
}

var credentialCreateCmd = &cobra.Command{
	Use:   "create [store] [id]",
	Short: i18n.T("Create a credential from Jenkins"),
	Long:  i18n.T("Create a credential from Jenkins"),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		switch credentialCreateOption.Type {
		case "basic", "secret":
		default:
			err = fmt.Errorf("unknow credential type: %s", credentialCreateOption.Type)
			return
		}

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
					ID:          credentialCreateOption.ID,
					Description: credentialCreateOption.Description,
				},
			})
		case "secret":
			err = jClient.CreateSecret(credentialCreateOption.Store, client.StringCredentials{
				Secret: credentialCreateOption.Secret,
				Credential: client.Credential{
					ID:          credentialCreateOption.ID,
					Description: credentialCreateOption.Description,
				},
			})
		default:
			err = fmt.Errorf("unknow credential type: %s", credentialCreateOption.Type)
		}
		return
	},
}

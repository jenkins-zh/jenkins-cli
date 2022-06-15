package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// UserTokenOption represents a user token cmd option
type UserTokenOption struct {
	Generate   bool
	Name       string
	TargetUser string

	RoundTripper http.RoundTripper
}

var userTokenOption UserTokenOption

func init() {
	userCmd.AddCommand(userTokenCmd)
	userTokenCmd.Flags().BoolVarP(&userTokenOption.Generate, "generate", "g", false,
		i18n.T("Generate the token"))
	userTokenCmd.Flags().StringVarP(&userTokenOption.Name, "name", "n", "",
		i18n.T("Name of the token"))
	userTokenCmd.Flags().StringVarP(&userTokenOption.TargetUser, "target-user", "", "",
		i18n.T("The target user of the new token"))
}

var userTokenCmd = &cobra.Command{
	Use:     "token",
	Short:   i18n.T("Token the user of your Jenkins"),
	Long:    i18n.T("Token the user of your Jenkins"),
	Example: `jcli user token -g`,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		if !userTokenOption.Generate {
			cmd.Help()
			return
		}

		jclient := &client.UserClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: userTokenOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		var token *client.Token
		token, err = jclient.CreateToken(userTokenOption.TargetUser, userTokenOption.Name)
		if err == nil {
			var data []byte
			data, err = userOption.Output(token)
			if err == nil {
				cmd.Println(string(data))
			}
		}
		return
	},
}

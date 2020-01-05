package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// UserTokenOption represents a user token cmd option
type UserTokenOption struct {
	Generate bool
	Name     string

	RoundTripper http.RoundTripper
}

var userTokenOption UserTokenOption

func init() {
	userCmd.AddCommand(userTokenCmd)
	userTokenCmd.Flags().BoolVarP(&userTokenOption.Generate, "generate", "g", false, "Generate the token")
	userTokenCmd.Flags().StringVarP(&userTokenOption.Name, "name", "n", "", "Name of the token")
}

var userTokenCmd = &cobra.Command{
	Use:   "token -g",
	Short: "Token the user of your Jenkins",
	Long:  `Token the user of your Jenkins`,
	Run: func(cmd *cobra.Command, _ []string) {
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

		tokenName := userTokenOption.Name
		status, err := jclient.CreateToken(tokenName)
		if err == nil {
			var data []byte
			data, err = userOption.Output(status)
			if err == nil {
				cmd.Println(string(data))
			}
		}
		helper.CheckErr(cmd, err)
	},
}

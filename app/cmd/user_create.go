package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// UserCreateOption is user create cmd option
type UserCreateOption struct {
	RoundTripper http.RoundTripper
}

var userCreateOption UserCreateOption

func init() {
	userCmd.AddCommand(userCreateCmd)
}

var userCreateCmd = &cobra.Command{
	Use:   "create <username>",
	Short: "Create a user for your Jenkins",
	Long:  `Create a user for your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		username := args[0]

		jclient := &client.UserClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: userCreateOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if user, err := jclient.Create(username); err == nil {
			cmd.Println("create user success. Password is:", user.Password1)
		} else {
			cmd.PrintErrln(err)
		}
	},
}

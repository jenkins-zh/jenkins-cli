package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// UserDeleteOption is user delete cmd option
type UserDeleteOption struct {
	BatchOption

	RoundTripper http.RoundTripper
}

var userDeleteOption UserDeleteOption

func init() {
	userCmd.AddCommand(userDeleteCmd)
	userDeleteCmd.Flags().BoolVarP(&userDeleteOption.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

var userDeleteCmd = &cobra.Command{
	Use:     "delete <username>",
	Aliases: GetAliasesDel(),
	Short:   "Delete a user for your Jenkins",
	Long:    `Delete a user for your Jenkins`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		if !userDeleteOption.Confirm(fmt.Sprintf("Are you sure to delete user %s ?", username)) {
			return
		}

		jclient := &client.UserClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: userDeleteOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		err := jclient.Delete(username)
		helper.CheckErr(cmd, err)
	},
}

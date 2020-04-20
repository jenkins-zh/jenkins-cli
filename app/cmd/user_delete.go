package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// UserDeleteOption is user delete cmd option
type UserDeleteOption struct {
	common.BatchOption

	RoundTripper http.RoundTripper
}

var userDeleteOption UserDeleteOption

func init() {
	userCmd.AddCommand(userDeleteCmd)
	userDeleteCmd.Flags().BoolVarP(&userDeleteOption.Batch, "batch", "b", false,
		i18n.T("Batch mode, no need confirm"))
	userDeleteOption.BatchOption.Stdio = common.GetSystemStdio()
}

var userDeleteCmd = &cobra.Command{
	Use:     "delete <username>",
	Aliases: common.GetAliasesDel(),
	Short:   "Delete a user for your Jenkins",
	Long:    `Delete a user for your Jenkins`,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
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
		return jclient.Delete(username)
	},
}

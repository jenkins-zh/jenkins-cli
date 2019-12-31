package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// UserEditOption is the user edit cmd option
type UserEditOption struct {
	CommonOption

	Description string
}

var userEditOption UserEditOption

func init() {
	userCmd.AddCommand(userEditCmd)
	userEditCmd.Flags().StringVarP(&userEditOption.Description, "desc", "d", "",
		i18n.T("Edit the description"))
}

var userEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the user of your Jenkins",
	Long:  `Edit the user of your Jenkins`,
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		jClient := &client.UserClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: userEditOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		var user *client.User
		if user, err = jClient.Get(); err == nil {
			var content string
			content, err = userEditOption.Editor(user.Description, "Edit user description")
			if err == nil {
				err = jClient.EditDesc(content)
			}
		}
		return
	},
}

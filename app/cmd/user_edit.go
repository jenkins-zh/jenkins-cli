package cmd

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type UserEditOption struct {
	Description bool
}

var userEditOption UserEditOption

func init() {
	userCmd.AddCommand(userEditCmd)
	userEditCmd.Flags().BoolVarP(&userEditOption.Description, "desc", "d", false, "Edit the description")
}

var userEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the user of your Jenkins",
	Long:  `Edit the user of your Jenkins`,
	Run: func(_ *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.UserClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if status, err := jclient.Get(); err == nil {
			description := status.Description

			prompt := &survey.Editor{
				Message:       "Edit user description",
				FileName:      "*.sh",
				Default:       description,
				HideDefault:   true,
				AppendDefault: true,
			}

			if err = survey.AskOne(prompt, &description); err != nil {
				log.Fatal(err)
			} else {
				if err = jclient.EditDesc(description); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			log.Fatal(err)
		}
	},
}

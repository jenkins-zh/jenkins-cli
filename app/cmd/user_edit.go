package cmd

import (
	"log"

	"github.com/AlecAivazis/survey"
	"github.com/linuxsuren/jenkins-cli/client"
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
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		jclient := &client.UserClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if status, err := jclient.Get(); err == nil {
			description := status.Description

			prompt := &survey.Editor{
				Message:       "Edit your pipeline script",
				FileName:      "*.sh",
				Default:       description,
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

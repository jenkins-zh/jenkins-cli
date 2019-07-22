package cmd

import (
	"fmt"
	"log"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type UserCreateOption struct {
}

var userCreateOption UserCreateOption

func init() {
	userCmd.AddCommand(userCreateCmd)
}

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a user for your Jenkins",
	Long:  `Create a user for your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		username := args[0]

		jenkins := getCurrentJenkins()
		jclient := &client.UserClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if user, err := jclient.Create(username); err == nil {
			fmt.Printf("create user success. Password is: %s\n", user.Password1)
		} else {
			log.Fatal(err)
		}
	},
}

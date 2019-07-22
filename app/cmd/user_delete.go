package cmd

import (
	"fmt"
	"log"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type UserDeleteOption struct {
	BatchOption
}

var userDeleteOption UserDeleteOption

func init() {
	userCmd.AddCommand(userDeleteCmd)
	userDeleteCmd.Flags().BoolVarP(&userDeleteOption.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user for your Jenkins",
	Long:  `Delete a user for your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		username := args[0]

		if !userDeleteOption.Confirm(fmt.Sprintf("Are you sure to delete user %s ?", username)) {
			return
		}

		jenkins := getCurrentJenkins()
		jclient := &client.UserClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if err := jclient.Delete(username); err == nil {
			fmt.Printf("delete user success.\n")
		} else {
			log.Fatal(err)
		}
	},
}

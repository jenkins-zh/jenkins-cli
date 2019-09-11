package cmd

import (
	"fmt"
	"log"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type UserTokenOption struct {
	Generate bool
	Name     string
}

var userTokenOption UserTokenOption

func init() {
	userCmd.AddCommand(userTokenCmd)
	userTokenCmd.Flags().BoolVarP(&userTokenOption.Generate, "generate", "g", false, "Generate the token")
	userTokenCmd.Flags().StringVarP(&userTokenOption.Name, "name", "n", "", "Name of the token")
}

var userTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Token the user of your Jenkins",
	Long:  `Token the user of your Jenkins`,
	Run: func(cmd *cobra.Command, _ []string) {
		if !userTokenOption.Generate {
			cmd.Help()
			return
		}

		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.UserClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		tokenName := userTokenOption.Name
		if status, err := jclient.CreateToken(tokenName); err == nil {
			var data []byte
			if data, err = userOption.Output(status); err == nil {
				fmt.Printf("%s\n", string(data))
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

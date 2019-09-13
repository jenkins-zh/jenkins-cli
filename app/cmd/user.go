package cmd

import (
	"fmt"
	"log"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// UserOption is the user cmd option
type UserOption struct {
	OutputOption
}

var userOption UserOption

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.Flags().StringVarP(&userOption.Format, "output", "o", "json", "Format the output")
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Print the user of your Jenkins",
	Long:  `Print the user of your Jenkins`,
	Run: func(_ *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.UserClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if status, err := jclient.Get(); err == nil {
			if data, err := userOption.Output(status); err == nil {
				fmt.Println(string(data))
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// UserOption is the user cmd option
type UserOption struct {
	OutputOption

	RoundTripper http.RoundTripper
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
	Run: func(cmd *cobra.Command, _ []string) {
		jclient := &client.UserClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: userOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		status, err := jclient.Get()
		if err == nil {
			data, err := userOption.Output(status)
			if err == nil {
				cmd.Println(string(data))
			}
		}
		helper.CheckErr(cmd, err)
	},
}

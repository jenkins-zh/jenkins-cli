package cmd

import (
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

		if status, err := jclient.Get(); err == nil {
			if data, err := userOption.Output(status); err == nil {
				cmd.Println(string(data))
			} else {
				cmd.PrintErrln(err)
			}
		} else {
			cmd.PrintErrln(err)
		}
	},
}

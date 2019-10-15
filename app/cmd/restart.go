package cmd

import (
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// RestartOption holds the options for restart cmd
type RestartOption struct {
	BatchOption

	RoundTripper http.RoundTripper
}

var restartOption RestartOption

func init() {
	rootCmd.AddCommand(restartCmd)
	restartCmd.Flags().BoolVarP(&restartOption.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart your Jenkins",
	Long:  `Restart your Jenkins`,
	Run: func(cmd *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		if !restartOption.Batch {
			confirm := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Are you sure to restart Jenkins %s?", jenkins.URL),
			}
			if err := survey.AskOne(prompt, &confirm); !confirm {
				return
			} else if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		jclient := &client.CoreClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: restartOption.RoundTripper,
				Debug:        rootOptions.Debug,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if err := jclient.Restart(); err == nil {
			cmd.Println("Please wait while Jenkins is restarting")
		} else {
			cmd.PrintErrln(err)
		}
	},
}

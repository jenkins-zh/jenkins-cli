package cmd

import (
	"log"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobSearchOption struct {
	OutputOption

	RoundTripper http.RoundTripper
}

var jobSearchOption JobSearchOption

func init() {
	jobCmd.AddCommand(jobSearchCmd)
	jobSearchCmd.Flags().StringVarP(&jobSearchOption.Format, "output", "o", "json", "Format the output")
}

var jobSearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Print the job of your Jenkins",
	Long:  `Print the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		keyword := args[0]

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobSearchOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if status, err := jclient.Search(keyword); err == nil {
			var data []byte
			if data, err = Format(status, queueOption.Format); err == nil {
				cmd.Println(string(data))
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

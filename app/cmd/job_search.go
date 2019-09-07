package cmd

import (
	"log"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobSearchOption is the options of job search command
type JobSearchOption struct {
	OutputOption
	PrintAll bool
	Max      int

	RoundTripper http.RoundTripper
}

var jobSearchOption JobSearchOption

func init() {
	jobCmd.AddCommand(jobSearchCmd)
	jobSearchCmd.Flags().IntVarP(&jobSearchOption.Max, "max", "", 10, "The number of limitation to print")
	jobSearchCmd.Flags().BoolVarP(&jobSearchOption.PrintAll, "all", "", false, "Print all items if there's no keyword")
	jobSearchCmd.Flags().StringVarP(&jobSearchOption.Format, "output", "o", "json", "Format the output")
}

var jobSearchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Print the job of your Jenkins",
	Long:  `Print the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if !jobSearchOption.PrintAll && len(args) == 0 {
			cmd.Help()
			return
		}

		if jobSearchOption.PrintAll && len(args) == 0 {
			args = []string{""}
		}

		keyword := args[0]

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobSearchOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if status, err := jclient.Search(keyword, jobSearchOption.Max); err == nil {
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

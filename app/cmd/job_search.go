package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
			if data, err = jobSearchOption.Output(status); err == nil {
				cmd.Println(string(data))
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

// Output render data into byte array
func (o *JobSearchOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil && o.OutputOption.Format == "name" {
		buf := ""
		searchResult := obj.(*client.SearchResult)

		for _, item := range searchResult.Suggestions {
			buf = fmt.Sprintf("%s%s\n", buf, item.Name)
		}
		data = []byte(strings.Trim(buf, "\n"))
		err = nil
	}
	return
}

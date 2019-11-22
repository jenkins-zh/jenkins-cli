package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"net/http"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

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
	jobSearchCmd.Flags().IntVarP(&jobSearchOption.Max, "max", "", 10,
		i18n.T("The number of limitation to print"))
	jobSearchCmd.Flags().BoolVarP(&jobSearchOption.PrintAll, "all", "", false,
		i18n.T("Print all items if there's no keyword"))
	jobSearchCmd.Flags().StringVarP(&jobSearchOption.Format, "output", "o", "json",
		i18n.T(`Formats of the output which contain name, path`))
}

var jobSearchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: i18n.T("Print the job of your Jenkins"),
	Long:  i18n.T(`Print the job of your Jenkins`),
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

		status, err := jclient.Search(keyword, jobSearchOption.Max)
		if err == nil {
			var data []byte
			data, err = jobSearchOption.Output(status)
			if err == nil {
				cmd.Println(string(data))
			}
		}
		helper.CheckErr(cmd, err)
	},
}

// Output render data into byte array
func (o *JobSearchOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		var formatFunc JobNameFormat

		switch o.OutputOption.Format {
		case "name":
			formatFunc = simpleFormat
		case "path":
			formatFunc = pathFormat
		}

		if formatFunc == nil {
			err = fmt.Errorf("unknow format %s", o.OutputOption.Format)
			return
		}

		buf := ""
		searchResult := obj.(*client.SearchResult)

		for _, item := range searchResult.Suggestions {
			buf = fmt.Sprintf("%s%s\n", buf, formatFunc(item.Name))
		}
		data = []byte(strings.Trim(buf, "\n"))
		err = nil
	}
	return
}

// JobNameFormat format the job name
type JobNameFormat func(string) string

func simpleFormat(name string) string {
	return name
}

func pathFormat(name string) string {
	return client.ParseJobPath(name)
}

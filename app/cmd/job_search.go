package cmd

import (
	"fmt"
	"log"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobSearchOption struct {
	OutputOption
}

var jobSearchOption JobSearchOption

func init() {
	jobCmd.AddCommand(jobSearchCmd)
	jobSearchCmd.PersistentFlags().StringVarP(&queueOption.Format, "output", "o", "json", "Format the output")
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

		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if status, err := jclient.Search(keyword); err == nil {
			var data []byte
			if data, err = Format(status, queueOption.Format); err == nil {
				fmt.Printf("%s\n", string(data))
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

package cmd

import (
	"fmt"
	"log"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobOption struct {
	OutputOption

	Name    string
	History bool
}

var jobOption JobOption

func init() {
	rootCmd.AddCommand(jobCmd)
	jobCmd.PersistentFlags().StringVarP(&jobOption.Format, "output", "o", "json", "Format the output")
	jobCmd.PersistentFlags().StringVarP(&jobOption.Name, "name", "n", "", "Name of the job")
	jobCmd.PersistentFlags().BoolVarP(&jobOption.History, "history", "", false, "Print the build history of job")
}

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Print the job of your Jenkins",
	Long:  `Print the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if jobOption.Name == "" {
			log.Fatal("need a name")
		}

		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth
		if jobOption.History {
			if builds, err := jclient.GetHistory(jobOption.Name); err == nil {
				var data []byte
				if data, err = Format(builds, jobOption.Format); err == nil {
					fmt.Printf("%s\n", string(data))
				} else {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
	},
}

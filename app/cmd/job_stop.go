package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobStopOption struct {
	BatchOption
}

var jobStopOption JobStopOption

func init() {
	jobCmd.AddCommand(jobStopCmd)
	jobStopOption.SetFlag(jobStopCmd)
}

var jobStopCmd = &cobra.Command{
	Use:   "stop <jobName> <buildNumbe>",
	Short: "Stop a job build in your Jenkins",
	Long:  `Stop a job build in your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}

		var (
			buildNum int
			err      error
		)
		if buildNum, err = strconv.Atoi(args[1]); err != nil {
			log.Fatal(err)
		}

		jobName := args[0]
		if !jobStopOption.Confirm(fmt.Sprintf("Are you sure to stop job %s ?", jobName)) {
			return
		}

		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if err := jclient.StopJob(jobName, buildNum); err != nil {
			log.Fatal(err)
		}
	},
}

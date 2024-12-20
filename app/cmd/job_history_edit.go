package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type jobHistoryEdit struct {
	displayName string
	description string
	id          int
}

func createJobHistoryEditCmd() (cmd *cobra.Command) {
	opt := &jobHistoryEdit{}
	cmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit job history",
		RunE:  opt.RunE,
		Args:  cobra.MinimumNArgs(1),
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.displayName, "displayName", "d", "", "Display name of target job history")
	flags.StringVarP(&opt.description, "description", "m", "", "Description of target job history")
	flags.IntVarP(&opt.id, "id", "i", -1, "ID of job history")
	return
}

func (o *jobHistoryEdit) RunE(cmd *cobra.Command, args []string) (err error) {
	jobName := args[0]

	jClient := &client.JobClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: jobHistoryOption.RoundTripper,
		},
	}
	getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

	err = jClient.EditBuild(jobName, o.id, o.displayName, o.description)
	return
}

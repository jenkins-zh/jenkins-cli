package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	cobra_ext "github.com/linuxsuren/cobra-extension/pkg"
	"strconv"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobArtifactOption is the options of job artifact command
type JobArtifactOption struct {
	cobra_ext.OutputOption
	common.Option
}

var jobArtifactOption JobArtifactOption

func init() {
	jobCmd.AddCommand(jobArtifactCmd)
	jobArtifactOption.SetFlagWithHeaders(jobArtifactCmd, "Name,Path,Size")
}

var jobArtifactCmd = &cobra.Command{
	Use:   "artifact <jobName> [buildID]",
	Short: i18n.T("Print the artifact list of target job"),
	Long:  i18n.T("Print the artifact list of target job"),
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		argLen := len(args)
		jobName := args[0]
		buildID := -1

		if argLen >= 2 {
			if buildID, err = strconv.Atoi(args[1]); err != nil {
				return
			}
		}

		jclient := &client.ArtifactClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobArtifactOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		var artifacts []client.Artifact
		if artifacts, err = jclient.List(jobName, buildID); err == nil {
			jobArtifactOption.Writer = cmd.OutOrStdout()
			err = jobArtifactOption.OutputV2(artifacts)
		}
		return
	},
}

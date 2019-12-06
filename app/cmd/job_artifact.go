package cmd

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"net/http"
	"strconv"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// JobArtifactOption is the options of job artifact command
type JobArtifactOption struct {
	OutputOption

	RoundTripper http.RoundTripper
}

var jobArtifactOption JobArtifactOption

func init() {
	jobCmd.AddCommand(jobArtifactCmd)
	jobArtifactOption.SetFlag(jobArtifactCmd)
}

var jobArtifactCmd = &cobra.Command{
	Use:   "artifact <jobName> [buildID]",
	Short: i18n.T("Print the artifact list of target job"),
	Long:  i18n.T("Print the artifact list of target job"),
	Run: func(cmd *cobra.Command, args []string) {
		argLen := len(args)
		if argLen == 0 {
			cmd.Help()
			return
		}

		var err error
		jobName := args[0]
		buildID := -1

		if argLen >= 2 {
			if buildID, err = strconv.Atoi(args[1]); err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		jclient := &client.ArtifactClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobArtifactOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		artifacts, err := jclient.List(jobName, buildID)
		if err == nil {
			var data []byte
			data, err = jobArtifactOption.Output(artifacts)
			if err == nil {
				cmd.Print(string(data))
			}
		}
		helper.CheckErr(cmd, err)
	},
}

// Output render data into byte array
func (o *JobArtifactOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		artifacts := obj.([]client.Artifact)
		buf := new(bytes.Buffer)

		table := util.CreateTable(buf)
		table.AddRow("id", "name", "path", "size")
		for _, artifact := range artifacts {
			table.AddRow(artifact.ID, artifact.Name, artifact.Path, fmt.Sprintf("%d", artifact.Size))
		}
		table.Render()
		err = nil
		data = buf.Bytes()
	}
	return
}

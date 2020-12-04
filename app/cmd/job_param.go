package cmd

import (
	"encoding/json"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	cobra_ext "github.com/linuxsuren/cobra-extension"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobParamOption is the job param option
type JobParamOption struct {
	cobra_ext.OutputOption

	Indent bool

	RoundTripper http.RoundTripper
}

var jobParamOption JobParamOption

func init() {
	jobCmd.AddCommand(jobParamCmd)
	jobParamCmd.Flags().BoolVarP(&jobParamOption.Indent, "indent", "", false, "Output with indent")
	jobParamOption.SetFlag(jobParamCmd)
}

var jobParamCmd = &cobra.Command{
	Use:   "param <jobName>",
	Short: i18n.T("Get parameters of the job of your Jenkins"),
	Long:  i18n.T("Get parameters of the job of your Jenkins"),
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobParamOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		job, err := jclient.GetJob(name)
		var data []byte
		if err == nil {
			proCount := len(job.Property)
			if proCount != 0 {
				for _, pro := range job.Property {
					if len(pro.ParameterDefinitions) == 0 {
						continue
					}

					if jobParamOption.Indent {
						data, err = json.MarshalIndent(pro.ParameterDefinitions, "", " ")
					} else {
						data, err = json.Marshal(pro.ParameterDefinitions)
					}
					break
				}
			}
		}
		if err == nil && len(data) > 0 {
			cmd.Println(string(data))
		}
		helper.CheckErr(cmd, err)
	},
}

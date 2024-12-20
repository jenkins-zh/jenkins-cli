package cmd

import (
	"encoding/json"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	cobra_ext "github.com/linuxsuren/cobra-extension/pkg"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobParamOption is the job param option
type JobParamOption struct {
	cobra_ext.OutputOption

	Indent bool
	Remove string
	Add    string

	RoundTripper http.RoundTripper
}

var jobParamOption JobParamOption

func init() {
	jobCmd.AddCommand(jobParamCmd)
	flags := jobParamCmd.Flags()
	flags.BoolVarP(&jobParamOption.Indent, "indent", "", false, "Output with indent")
	flags.StringVarP(&jobParamOption.Add, "add", "", "",
		`Add parameters into the Pipeline. Example data: [{"name":"name","value":"rick","desc":"this is a name"}]`)
	flags.StringVarP(&jobParamOption.Remove, "remove", "", "",
		`Remove parameters from the Pipeline. Example data: name,age`)
	jobParamOption.SetFlag(jobParamCmd)
}

var jobParamCmd = &cobra.Command{
	Use:   "param",
	Short: i18n.T("Get parameters of the job of your Jenkins"),
	Long:  i18n.T("Get parameters of the job of your Jenkins"),
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		name := args[0]
		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobParamOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		if jobParamOption.Remove != "" {
			if err = jclient.RemoveParameters(name, jobParamOption.Remove); err != nil {
				return
			}
		}

		if jobParamOption.Add != "" {
			if err = jclient.AddParameters(name, jobParamOption.Add); err != nil {
				return
			}
		}

		var job *client.Job
		var data []byte
		if job, err = jclient.GetJob(name); err == nil {
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
			if len(data) > 0 {
				cmd.Println(string(data))
			}
		}
		return
	},
}

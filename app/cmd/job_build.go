package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobBuildOption is the job build option
type JobBuildOption struct {
	common.BatchOption
	common.CommonOption

	Param      string
	ParamArray []string
}

var jobBuildOption JobBuildOption

// ResetJobBuildOption give it a clean option struct
func ResetJobBuildOption() {
	jobBuildOption = JobBuildOption{}
}

func init() {
	jobCmd.AddCommand(jobBuildCmd)
	jobBuildOption.SetFlag(jobBuildCmd)
	jobBuildCmd.Flags().StringVarP(&jobBuildOption.Param, "param", "", "",
		i18n.T("Parameters of the job which is JSON format"))
	jobBuildCmd.Flags().StringArrayVar(&jobBuildOption.ParamArray, "param-entry", nil,
		i18n.T("Parameters of the job which are the entry format, for example: --param-entry name=value"))
	jobBuildOption.BatchOption.Stdio = common.GetSystemStdio()
	jobBuildOption.CommonOption.Stdio = common.GetSystemStdio()
}

var jobBuildCmd = &cobra.Command{
	Use:   "build <jobName>",
	Short: i18n.T("Build the job of your Jenkins"),
	Long: i18n.T(`Build the job of your Jenkins.
You need to give the parameters if your pipeline has them. Learn more about it from https://jenkins.io/doc/book/pipeline/syntax/#parameters.`),
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if jobBuildOption.ParamArray == nil {
			return
		}

		paramDefs := make([]client.ParameterDefinition, 0)
		if jobBuildOption.Param != "" {
			if err = json.Unmarshal([]byte(jobBuildOption.Param), &paramDefs); err != nil {
				return
			}
		}

		for _, paramEntry := range jobBuildOption.ParamArray {
			if entryArray := strings.SplitN(paramEntry, "=", 2); len(entryArray) == 2 {
				paramDefs = append(paramDefs, client.ParameterDefinition{
					Name:  entryArray[0],
					Value: entryArray[1],
					Type:  "StringParameterDefinition",
				})
			}
		}

		var data []byte
		if data, err = json.Marshal(paramDefs); err == nil {
			jobBuildOption.Param = string(data)
		}
		return
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		name := args[0]

		if !jobBuildOption.Confirm(fmt.Sprintf("Are you sure to build job %s", name)) {
			return
		}

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobBuildOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		paramDefs := []client.ParameterDefinition{}
		hasParam := false

		var job *client.Job
		if jobBuildOption.Batch {
			if jobBuildOption.Param != "" {
				hasParam = true

				err = json.Unmarshal([]byte(jobBuildOption.Param), &paramDefs)
			}
		} else if job, err = jclient.GetJob(name); err == nil {
			proCount := len(job.Property)
			if proCount != 0 {
				for _, pro := range job.Property {
					if len(pro.ParameterDefinitions) == 0 {
						continue
					}

					var data []byte
					if data, err = json.MarshalIndent(pro.ParameterDefinitions, "", " "); err == nil {
						content := string(data)
						content, err = jobBuildOption.Editor(content, "Edit your pipeline script")
						if err == nil {
							err = json.Unmarshal([]byte(content), &paramDefs)
						}
					}
					hasParam = true
					break
				}
			}
		}

		if err == nil {
			if hasParam {
				err = jclient.BuildWithParams(name, paramDefs)
			} else {
				err = jclient.Build(name)
			}
		}
		return
	},
}

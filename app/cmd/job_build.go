package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	cobra_ext "github.com/linuxsuren/cobra-extension/pkg"
	"github.com/spf13/cobra"
)

// JobBuildOption is the job build option
type JobBuildOption struct {
	common.BatchOption
	common.Option
	cobra_ext.OutputOption

	Param      string
	ParamArray []string

	ParamFilePathArray []string

	Wait     bool
	WaitTime int
	Delay    int
	Cause    string
}

var jobBuildOption JobBuildOption

// ResetJobBuildOption give it a clean option struct
func ResetJobBuildOption() {
	jobBuildOption = JobBuildOption{}
}

func init() {
	jobCmd.AddCommand(jobBuildCmd)
	jobBuildCmd.Flags().BoolVarP(&jobBuildOption.Batch, "batch", "b", false, "Batch mode, no need confirm")
	jobBuildCmd.Flags().StringVarP(&jobBuildOption.Param, "param", "", "",
		i18n.T("Parameters of the job which is JSON format"))
	jobBuildCmd.Flags().StringArrayVar(&jobBuildOption.ParamArray, "param-entry", nil,
		i18n.T("Parameters of the job which are the entry format, for example: --param-entry name=value"))
	jobBuildCmd.Flags().StringArrayVar(&jobBuildOption.ParamFilePathArray, "param-file", nil,
		i18n.T("Parameters of the job which is file path, for example: --param-file name=filename"))
	jobBuildCmd.Flags().BoolVarP(&jobBuildOption.Wait, "wait", "", false,
		i18n.T("If you want to wait for the build ID from Jenkins. You need to install plugin pipeline-restful-api first"))
	jobBuildCmd.Flags().IntVarP(&jobBuildOption.WaitTime, "wait-timeout", "", 30,
		i18n.T("The timeout of seconds when you wait for the build ID"))
	jobBuildCmd.Flags().IntVarP(&jobBuildOption.Delay, "delay", "", 0,
		i18n.T("Delay when trigger a Jenkins job"))
	jobBuildCmd.Flags().StringVarP(&jobBuildOption.Cause, "cause", "", "triggered by jcli",
		i18n.T("The cause of a job build"))

	jobBuildOption.SetFlagWithHeaders(jobBuildCmd, "Number,URL")
	jobBuildOption.BatchOption.Stdio = common.GetSystemStdio()
	jobBuildOption.Option.Stdio = common.GetSystemStdio()
}

var jobBuildCmd = &cobra.Command{
	Use:   "build <jobName>",
	Short: i18n.T("Build the job of your Jenkins"),
	Long: i18n.T(`Build the job of your Jenkins.
You need to give the parameters if your pipeline has them. Learn more about it from https://jenkins.io/doc/book/pipeline/syntax/#parameters.`),
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(_ *cobra.Command, _ []string) (err error) {
		if jobBuildOption.ParamArray == nil && jobBuildOption.ParamFilePathArray == nil {
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
					Type:  client.StringParameterDefinition,
				})
			}
		}

		for _, filepathEntry := range jobBuildOption.ParamFilePathArray {
			if filepathArray := strings.SplitN(filepathEntry, "=", 2); len(filepathArray) == 2 {
				paramDefs = append(paramDefs, client.ParameterDefinition{
					Name:     filepathArray[0],
					Filepath: filepathArray[1],
					Type:     client.FileParameterDefinition,
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
			} else if jobBuildOption.Wait {
				var build client.IdentityBuild
				if build, err = jclient.BuildAndReturn(name, jobBuildOption.Cause, jobBuildOption.WaitTime, jobBuildOption.Delay); err == nil {
					jobBuildOption.Writer = cmd.OutOrStdout()
					err = jobBuildOption.OutputV2([1]client.SimpleJobBuild{build.Build.SimpleJobBuild})
				}
			} else {
				err = jclient.Build(name)
			}
		}
		return
	},
}

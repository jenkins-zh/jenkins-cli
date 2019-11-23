package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"log"
	"net/http"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobBuildOption is the job build option
type JobBuildOption struct {
	BatchOption

	Param      string
	ParamArray []string
	Debug      bool

	RoundTripper http.RoundTripper
}

var jobBuildOption JobBuildOption

func init() {
	jobCmd.AddCommand(jobBuildCmd)
	jobBuildCmd.Flags().BoolVarP(&jobBuildOption.Batch, "batch", "b", false,
		i18n.T("Batch mode, no need to confirm"))
	jobBuildCmd.Flags().StringVarP(&jobBuildOption.Param, "param", "", "",
		i18n.T("Parameters of the job which is JSON format"))
	jobBuildCmd.Flags().StringArrayVar(&jobBuildOption.ParamArray, "param-entry", nil,
		i18n.T("Parameters of the job which are the entry format, for example: --param-entry name=value"))
	jobBuildCmd.Flags().BoolVarP(&jobBuildOption.Debug, "verbose", "", false,
		i18n.T("Output the verbose"))
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
				})
			}
		}

		var data []byte
		if data, err = json.Marshal(paramDefs); err == nil {
			jobBuildOption.Param = string(data)
		}
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if !jobBuildOption.Batch && !jobBuildOption.Confirm(fmt.Sprintf("Are you sure to build job %s", name)) {
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

		if jobBuildOption.Batch {
			if jobBuildOption.Param != "" {
				hasParam = true

				if err := json.Unmarshal([]byte(jobBuildOption.Param), &paramDefs); err != nil {
					log.Fatal(err)
				}
			}
		} else if job, err := jclient.GetJob(name); err == nil {
			proCount := len(job.Property)
			if jobBuildOption.Debug {
				fmt.Println("Found properties ", proCount)
			}
			if proCount != 0 {
				for _, pro := range job.Property {
					if len(pro.ParameterDefinitions) == 0 {
						continue
					}

					if data, err := json.MarshalIndent(pro.ParameterDefinitions, "", " "); err == nil {
						content := string(data)
						prompt := &survey.Editor{
							Message:       "Edit your pipeline script",
							FileName:      "*.sh",
							Default:       content,
							HideDefault:   true,
							AppendDefault: true,
						}

						if err = survey.AskOne(prompt, &content); err != nil {
							log.Fatal(err)
						}

						if err = json.Unmarshal([]byte(content), &paramDefs); err != nil {
							log.Fatal(err)
						}
					}
					hasParam = true
					break
				}
			}
		} else {
			log.Fatal(err)
		}

		if hasParam {
			jclient.BuildWithParams(name, paramDefs)
		} else {
			if jobBuildOption.Debug {
				fmt.Println("Not params found")
			}
			jclient.Build(name)
		}
	},
}

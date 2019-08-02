package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobBuildOption struct {
	BatchOption

	Param string
	Debug bool
}

var jobBuildOption JobBuildOption

func init() {
	jobCmd.AddCommand(jobBuildCmd)
	jobBuildCmd.Flags().BoolVarP(&jobBuildOption.Batch, "batch", "b", false, "Batch mode, no need confirm")
	jobBuildCmd.Flags().StringVarP(&jobBuildOption.Param, "param", "", "", "Params of the job")
	jobBuildCmd.Flags().BoolVarP(&jobBuildOption.Debug, "verbose", "", false, "Output the verbose")
}

var jobBuildCmd = &cobra.Command{
	Use:   "build <jobName>",
	Short: "Build the job of your Jenkins",
	Long:  `Build the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		name := args[0]

		if !jobBuildOption.Batch && !jobBuildOption.Confirm(fmt.Sprintf("Are you sure to build job %s", name)) {
			return
		}

		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

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

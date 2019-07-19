package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AlecAivazis/survey"
	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobBuildOption struct {
	BatchOption
}

var jobBuildOption JobBuildOption

func init() {
	jobCmd.AddCommand(jobBuildCmd)
	jobBuildCmd.Flags().BoolVarP(&jobBuildOption.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

var jobBuildCmd = &cobra.Command{
	Use:   "build -n",
	Short: "Build the job of your Jenkins",
	Long:  `Build the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if jobOption.Name == "" {
			cmd.Help()
			return
		}

		if !jobBuildOption.Confirm(fmt.Sprintf("Are you sure to build job %s", jobOption.Name)) {
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
		if job, err := jclient.GetJob(jobOption.Name); err == nil {
			fmt.Println(job.Property)
			if len(job.Property) != 0 {
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
		}

		if hasParam {
			jclient.BuildWithParams(jobOption.Name, paramDefs)
		} else {
			jclient.Build(jobOption.Name)
		}
	},
}

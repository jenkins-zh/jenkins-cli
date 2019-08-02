package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type JobParamOption struct {
	OutputOption

	Indent bool
}

var jobParamOption JobParamOption

func init() {
	jobCmd.AddCommand(jobParamCmd)
	jobParamCmd.Flags().BoolVarP(&jobParamOption.Indent, "indent", "", false, "Output with indent")
	jobParamOption.SetFlag(jobParamCmd)
}

var jobParamCmd = &cobra.Command{
	Use:   "param <jobName>",
	Short: "Get param of the job of your Jenkins",
	Long:  `Get param of the job of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		name := args[0]
		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if job, err := jclient.GetJob(name); err == nil {
			proCount := len(job.Property)
			if jobBuildOption.Debug {
				fmt.Println("Found properties ", proCount)
			}
			if proCount != 0 {
				for _, pro := range job.Property {
					if len(pro.ParameterDefinitions) == 0 {
						continue
					}

					if jobParamOption.Indent {
						if data, err := json.MarshalIndent(pro.ParameterDefinitions, "", " "); err == nil {
							fmt.Println(string(data))
						} else {
							log.Fatal(err)
						}
					} else {
						if data, err := json.Marshal(pro.ParameterDefinitions); err == nil {
							fmt.Println(string(data))
						} else {
							log.Fatal(err)
						}
					}
					break
				}
			}
		} else {
			log.Fatal(err)
		}
	},
}

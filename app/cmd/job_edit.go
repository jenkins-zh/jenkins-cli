/*
MIT License

Copyright (c) 2019 Zhao Xiaojie

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"regexp"
)

// JobEditOption is the option for job create command
type JobEditOption struct {
	common.CommonOption

	Filename  string
	Script    string
	URL       string
	Sample    bool
	TrimSpace bool
	// Build flag indicates if trigger the Jenkins job after the action of edit
	Build bool
}

var jobEditOption JobEditOption

func init() {
	jobCmd.AddCommand(jobEditCmd)

	flags := jobEditCmd.Flags()
	flags.StringVarP(&jobEditOption.URL, "url", "", "",
		i18n.T("URL of the Jenkinsfile to files to use to replace pipeline"))
	flags.StringVarP(&jobEditOption.Filename, "filename", "f", "",
		i18n.T("Filename to files to use to replace pipeline"))
	flags.StringVarP(&jobEditOption.Script, "script", "s", "",
		i18n.T("Script to use to replace pipeline. Use script first if you give filename at the meantime."))
	flags.BoolVarP(&jobEditOption.Sample, "sample", "", false,
		i18n.T("Give it a sample Jenkinsfile if the target script is empty"))
	flags.BoolVarP(&jobEditOption.TrimSpace, "trim", "", true,
		i18n.T("If trim the leading and tailing white space"))
	flags.BoolVarP(&jobEditOption.Build, "build", "b", false,
		i18n.T("Indicates if trigger the Jenkins job after the action of edit"))

	jobEditOption.Stdio = common.GetSystemStdio()
}

var jobEditCmd = &cobra.Command{
	Use:   "edit",
	Short: i18n.T("Edit the job of your Jenkins"),
	Long: i18n.T(fmt.Sprintf(`Edit the job of your Jenkins. We only support to edit the pipeline job.
Official Pipeline syntax document is here https://jenkins.io/doc/book/pipeline/syntax/

%s`, common.GetEditorHelpText())),
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		name := args[0]
		var content string

		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobEditOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if content, err = jobEditOption.getPipeline(jclient, name); err == nil {
			if jobEditOption.TrimSpace {
				content = regexp.MustCompile("^\\W+|\\s+$").ReplaceAllString(content, "")
			}
			err = jclient.UpdatePipeline(name, content)
		}

		// just trigger the Jenkins job when edit it successfully
		if err == nil && jobEditOption.Build {
			err = jclient.Build(name)
		}
		return
	},
}

func (j *JobEditOption) getSampleJenkinsfile() string {
	return `pipeline {
    agent {
		    label 'master'
	  }
    options {
        timeout(time: 1, unit: 'HOURS') 
    }
    environment { 
        NAME = 'jack'
    }
    stages {
        stage('Example') {
            steps {
                echo 'Hello World'
                echo env.NAME
            }
        }
        stage('Parallel Stage') {
            failFast true
            parallel {
                stage('Brother A') {
                    steps {
                        echo "On Brother A"
                    }
                }
                stage('Brother B') {
                    steps {
                        echo "On Brother B"
                    }
                }
            }
		}
    }
    post { 
        always { 
            echo 'I will always say Hello again!'
        }
    }
}
`
}

//func getPipeline(name string) (script string, err error) {
func (j *JobEditOption) getPipeline(jClient *client.JobClient, name string) (script string, err error) {
	script = j.Script //we take the script from input firstly
	if script != "" {
		return
	}

	// take script from a file
	if script, err = j.getPipelineFromFile(); script != "" || err != nil {
		return
	}

	if script, err = j.getPipelineFromURL(jClient); script != "" || err != nil {
		return
	}

	var job *client.Pipeline
	if job, err = jClient.GetPipeline(name); err == nil {
		content := ""
		if job != nil {
			content = job.Script
		}

		// if the original script is empty, give it a sample script
		if content == "" && j.Sample {
			content = j.getSampleJenkinsfile()
		}

		j.EditFileName = fmt.Sprintf("Jenkinsfile.%s.groovy", base64.StdEncoding.EncodeToString([]byte(name)))
		script, err = j.Editor(content, "Edit your pipeline script")
	}
	return
}

func (j *JobEditOption) getPipelineFromFile() (script string, err error) {
	if j.Filename == "" {
		return
	}

	var data []byte
	if data, err = ioutil.ReadFile(j.Filename); err == nil {
		script = string(data)
	}
	return
}

func (j *JobEditOption) getPipelineFromURL(jClient *client.JobClient) (script string, err error) {
	if j.URL == "" {
		return
	}

	var resp *http.Response
	var body []byte
	httpClient := jClient.JenkinsCore.GetClient()
	if resp, err = httpClient.Get(j.URL); err == nil {
		if body, err = ioutil.ReadAll(resp.Body); err == nil {
			script = string(body)
		}
	}
	return
}

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
}

var jobEditOption JobEditOption

func init() {
	jobCmd.AddCommand(jobEditCmd)
	jobEditCmd.Flags().StringVarP(&jobEditOption.URL, "url", "", "",
		i18n.T("URL of the Jenkinsfile to files to use to replace pipeline"))
	jobEditCmd.Flags().StringVarP(&jobEditOption.Filename, "filename", "f", "",
		i18n.T("Filename to files to use to replace pipeline"))
	jobEditCmd.Flags().StringVarP(&jobEditOption.Script, "script", "s", "",
		i18n.T("Script to use to replace pipeline. Use script first if you give filename at the meantime."))
	jobEditCmd.Flags().BoolVarP(&jobEditOption.Sample, "sample", "", false,
		i18n.T("Give it a sample Jenkinsfile if the target script is empty"))
	jobEditCmd.Flags().BoolVarP(&jobEditOption.TrimSpace, "trim", "", true,
		i18n.T("If trim the leading and tailing white space"))
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
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		if content, err = jobEditOption.getPipeline(jclient, name); err == nil {
			if jobEditOption.TrimSpace {
				content = regexp.MustCompile("^\\W+|\\s+$").ReplaceAllString(content, "")
			}
			err = jclient.UpdatePipeline(name, content)
		}
		return
	},
}

func (j *JobEditOption) getSampleJenkinsfile() string {
	return `pipeline {
    agent {
		    label 'master'
	  }
    stages {
        stage('Example') {
            steps {
                echo 'Hello World'
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

package cmd

import (
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobEditOption is the option for job create command
type JobEditOption struct {
	CommonOption

	Filename string
	Script   string
	URL      string
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
	jobEditOption.Stdio = GetSystemStdio()
}

var jobEditCmd = &cobra.Command{
	Use:   "edit <jobName>",
	Short: i18n.T("Edit the job of your Jenkins"),
	Long:  i18n.T(`Edit the job of your Jenkins. We only support to edit the pipeline job.`),
	Args:  cobra.MinimumNArgs(1),
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
			err = jclient.UpdatePipeline(name, content)
		}
		return
	},
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

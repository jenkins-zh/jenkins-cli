package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
)

type jenkinsfileOption struct {
	RoundTripper http.RoundTripper
	output       string
	input        string
}

func newCenterJenkinsfileCommand() (cmd *cobra.Command) {
	opt := &jenkinsfileOption{}
	cmd = &cobra.Command{
		Use:     "convert",
		Short:   "Convert Jenkinsfile between JSON or groovy format",
		Example: `cat Jenkinsfile | jcli center convert`,
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}
	flags := cmd.Flags()
	flags.StringVarP(&opt.output, "output", "o", "json", "The expect format")
	flags.StringVarP(&opt.input, "input", "i", "", "The input file path")
	return
}

func (o *jenkinsfileOption) preRunE(cmd *cobra.Command, args []string) (err error) {
	if o.output != "json" && o.output != "jenkinsfile" {
		err = fmt.Errorf("not support format: %s", o.output)
	}
	return
}

func (o *jenkinsfileOption) runE(cmd *cobra.Command, args []string) (err error) {
	jClient := &core.Client{
		JenkinsCore: core.JenkinsCore{
			RoundTripper: o.RoundTripper,
		},
	}
	getCurrentJenkinsAndClientV2(&(jClient.JenkinsCore))

	var data []byte
	stat, _ := os.Stdin.Stat()
	if stat.Mode()&os.ModeCharDevice == 0 {
		if data, err = ioutil.ReadAll(cmd.InOrStdin()); err != nil {
			err = fmt.Errorf("failed to read from pipe")
			return
		}
	} else if o.input != "" {
		if data, err = ioutil.ReadFile(o.input); err != nil {
			err = fmt.Errorf("failed to read from file: %s, error: %v", o.input, err)
			return
		}
	} else {
		cmd.Println("data required from pipe")
		return
	}

	var result core.GenericResult
	switch o.output {
	case "jenkinsfile":
		if result, err = jClient.ToJenkinsfile(string(data)); err == nil && result.GetStatus() != "success" {
			err = fmt.Errorf("convert failed: %v", result.GetErrors())
		}
	case "json":
		if result, err = jClient.ToJSON(string(data)); err == nil && result.GetStatus() != "success" {
			err = fmt.Errorf("convert failed: %v", result.GetErrors())
		}
	}
	if err == nil {
		cmd.Println(result.GetResult())
	}
	return
}

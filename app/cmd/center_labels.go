package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"github.com/spf13/cobra"
	"net/http"
)

type labelOption struct {
	RoundTripper http.RoundTripper
}

func newCenterLabelCommand() (cmd *cobra.Command) {
	opt := &labelOption{}
	cmd = &cobra.Command{
		Use:   "labels",
		Short: "Print all the labels of the Jenkins agents",
		RunE:  opt.runE,
	}
	return
}

func (o *labelOption) runE(cmd *cobra.Command, args []string) (err error) {
	jClient := &core.Client{
		JenkinsCore: core.JenkinsCore{
			RoundTripper: o.RoundTripper,
		},
	}
	getCurrentJenkinsAndClientV2(&(jClient.JenkinsCore))
	var labelsRes *core.LabelsResponse
	if labelsRes, err = jClient.GetLabels(); err == nil {
		if labelsRes.Status == "ok" {
			for i := range labelsRes.Data {
				cmd.Println(labelsRes.Data[i].Label)
			}
		} else {
			err = fmt.Errorf("failed to get labels, status: %s", labelsRes.Status)
		}
	}
	return
}

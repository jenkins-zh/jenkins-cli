package center

import (
	"encoding/json"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/spf13/cobra"
)

func NewCenterIdentityCmd(client common.JenkinsClient) (cmd *cobra.Command) {
	opt := &CenterIdentityOption{
		JenkinsClient: client,
	}
	cmd = &cobra.Command{
		Use:   "identity",
		Short: i18n.T("Print the identity of current Jenkins"),
		Long:  i18n.T("Print the identity of current Jenkins"),
		RunE:  opt.RunE,
	}
	return
}

func (o *CenterIdentityOption) RunE(cmd *cobra.Command, _ []string) (err error) {
	jClient := &client.CoreClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: o.RoundTripper,
		},
	}
	o.JenkinsClient.GetCurrentJenkinsAndClient(&(jClient.JenkinsCore))

	var identity client.JenkinsIdentity
	var data []byte
	if identity, err = jClient.GetIdentity(); err == nil {
		if data, err = json.MarshalIndent(identity, "", " "); err == nil {
			cmd.Println(string(data))
		}
	}
	return
}

package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cascCmd)

	healthCheckRegister.Register(getCmdPath(cascCmd)+".*", &cascOptions)
}

// CASCOptions is the option of casc
type CASCOptions struct {
	common.CommonOption
}

var cascOptions CASCOptions

// Check do the health check of casc cmd
func (o *CASCOptions) Check() (err error) {
	opt := PluginOptions{
		CommonOption: common.CommonOption{RoundTripper: o.RoundTripper},
	}
	_, err = opt.FindPlugin("configuration-as-code")
	return
}

var cascCmd = &cobra.Command{
	Use:   "casc",
	Short: i18n.T("Configuration as Code"),
	Long:  i18n.T("Configuration as Code"),
	Annotations: map[string]string{
		common.Since: common.VersionSince0024,
	},
}

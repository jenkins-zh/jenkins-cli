package cmd

import (
	"fmt"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cascCmd)

	healthCheckRegister.Register(getCmdPath(cascCmd)+".*", &cascOptions)
}

// CASCOptions is the option of casc
type CASCOptions struct {
	CommonOption
}

var cascOptions CASCOptions

// Check do the health check of casc cmd
func (o *CASCOptions) Check() (err error) {
	jClient := &client.PluginManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: cascOptions.RoundTripper,
		},
	}
	getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

	support := false
	var plugins *client.InstalledPluginList
	if plugins, err = jClient.GetPlugins(1); err == nil {
		for _, plugin := range plugins.Plugins {
			if plugin.ShortName == "configuration-as-code" {
				support = true
				break
			}
		}
	}

	if !support {
		err = fmt.Errorf(i18n.T("lack of plugin configuration-as-code"))
	}
	return
}

var cascCmd = &cobra.Command{
	Use:   "casc",
	Short: i18n.T("Configuration as Code"),
	Long:  i18n.T("Configuration as Code"),
}

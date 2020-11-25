package config

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/spf13/cobra"
)

// NewConfigPluginListCmd create a command for list all jcli plugins
func NewConfigPluginListCmd(opt *common.Option) (cmd *cobra.Command) {
	configPluginListCmd := configPluginListCmd{
		Option: opt,
	}

	cmd = &cobra.Command{
		Use:               "list",
		Short:             "List all installed plugins",
		Long:              "List all installed plugins",
		RunE:              configPluginListCmd.RunE,
		ValidArgsFunction: common.NoFileCompletion,
		Annotations: map[string]string{
			common.Since: common.VersionSince0028,
		},
	}

	configPluginListCmd.SetFlagWithHeaders(cmd, "Use,Version,Installed,DownloadLink")
	return
}

// RunE is the main entry point of config plugin list command
func (c *configPluginListCmd) RunE(cmd *cobra.Command, args []string) (err error) {
	c.Writer = cmd.OutOrStdout()
	var plugins []plugin
	if plugins, err = findPlugins(); err == nil {
		err = c.OutputV2(plugins)
	}
	return
}

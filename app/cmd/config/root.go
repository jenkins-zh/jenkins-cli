package config

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	goPlugin "github.com/linuxsuren/go-cli-plugin/pkg/cmd"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

// NewConfigPluginCmd create a command as root of config plugin
func NewConfigPluginCmd(opt *common.Option) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "plugin",
		Short: i18n.T("Manage plugins for jcli"),
		Long: i18n.T(`Manage plugins for jcli
If you want to submit a plugin for jcli, please see also the following project.
https://github.com/jenkins-zh/jcli-plugins`),
		Annotations: map[string]string{
			common.Since: common.VersionSince0028,
		},
	}

	goPlugin.AppendPluginCmd(cmd, "jenkins-zh", "jcli-plugins")
	return
}

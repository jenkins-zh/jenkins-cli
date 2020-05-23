package config_plugin

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// NewConfigPluginUninstallCmd create a command to uninstall a plugin
func NewConfigPluginUninstallCmd(opt *common.CommonOption) (cmd *cobra.Command) {
	jcliPluginUninstallCmd := jcliPluginUninstallCmd{
		CommonOption: opt,
	}

	cmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Remove a plugin",
		Long:  "Remove a plugin",
		Args:  cobra.MinimumNArgs(1),
		RunE:  jcliPluginUninstallCmd.RunE,
		Annotations: map[string]string{
			common.Since: common.VersionSince0028,
		},
	}
	return
}

// RunE is the main entry point of command
func (c *jcliPluginUninstallCmd) RunE(cmd *cobra.Command, args []string) (err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	name := args[0]
	cachedMetadataFile := fmt.Sprintf("%s/.jenkins-cli/pluginss/%s.yaml", userHome, name)

	var data []byte
	if data, err = ioutil.ReadFile(cachedMetadataFile); err == nil {
		plugin := &plugin{}
		if err = yaml.Unmarshal(data, plugin); err == nil {
			mainFile := fmt.Sprintf("%s/.jenkins-cli/pluginss/%s", userHome, plugin.Main)

			os.Remove(cachedMetadataFile)
			os.Remove(mainFile)
		}
	}
	return
}

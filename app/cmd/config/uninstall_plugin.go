package config

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// NewConfigPluginUninstallCmd create a command to uninstall a plugin
func NewConfigPluginUninstallCmd(opt *common.Option) (cmd *cobra.Command) {
	jcliPluginUninstallCmd := jcliPluginUninstallCmd{
		Option: opt,
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
	cachedMetadataFile := common.GetJCLIPluginPath(userHome, name, false)

	var data []byte
	if data, err = ioutil.ReadFile(cachedMetadataFile); err == nil {
		plugin := &plugin{}
		if err = yaml.Unmarshal(data, plugin); err == nil {
			mainFile := common.GetJCLIPluginPath(userHome, plugin.Main, true)

			os.Remove(cachedMetadataFile)
			os.Remove(mainFile)
		}
	} else if os.IsNotExist(err) {
		err = nil
		cmd.Printf("plugin \"%s\" does not exists\n", name)
	}
	return
}

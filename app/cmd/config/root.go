package config

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	appCfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

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

	cmd.AddCommand(NewConfigPluginListCmd(opt),
		NewConfigPluginFetchCmd(opt),
		NewConfigPluginInstallCmd(opt),
		NewConfigPluginUninstallCmd(opt))
	return
}

func findPlugins() (plugins []plugin, err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	plugins = make([]plugin, 0)
	pluginsDir := fmt.Sprintf("%s/.jenkins-cli/plugins-repo/*.yaml", userHome)
	if files, err := filepath.Glob(pluginsDir); err == nil {
		for _, metaFile := range files {
			var data []byte
			plugin := plugin{}
			if data, err = ioutil.ReadFile(metaFile); err == nil {
				if err = yaml.Unmarshal(data, &plugin); err != nil {
					fmt.Println(err)
				} else {
					if plugin.Main == "" {
						plugin.Main = fmt.Sprintf("jcli-%s-plugin", plugin.Use)
					}

					if _, fileErr := os.Stat(common.GetJCLIPluginPath(userHome, plugin.Main, true)); !os.IsNotExist(fileErr) {
						plugin.Installed = true
					}
					plugins = append(plugins, plugin)
				}
			}
		}
	}
	return
}

// LoadPlugins loads the plugins
func LoadPlugins(cmd *cobra.Command) {
	var plugins []plugin
	var err error
	if plugins, err = findPlugins(); err != nil {
		cmd.PrintErrln("Cannot load plugins successfully")
		return
	}
	//cmd.Println("found plugins, count", len(plugins))

	for _, plugin := range plugins {
		// This function is used to setup the environment for the plugin and then
		// call the executable specified by the parameter 'main'
		callPluginExecutable := func(cmd *cobra.Command, main string, argv []string, out io.Writer) error {
			env := os.Environ()

			prog := exec.Command(main, argv...)
			prog.Env = env
			prog.Stdin = os.Stdin
			prog.Stdout = out
			prog.Stderr = os.Stderr
			if err := prog.Run(); err != nil {
				if eerr, ok := err.(*exec.ExitError); ok {
					os.Stderr.Write(eerr.Stderr)
					status := eerr.Sys().(syscall.WaitStatus)
					return pluginError{
						error: errors.Errorf("plugin %s exited with error", main),
						code:  status.ExitStatus(),
					}
				}
				return err
			}

			return nil
		}

		//cmd.Println("register plugin name", plugin.Use)
		c := &cobra.Command{
			Use:   plugin.Use,
			Short: plugin.Short,
			Long:  plugin.Long,
			RunE: func(cmd *cobra.Command, args []string) (err error) {
				var userHome string
				if userHome, err = homedir.Dir(); err != nil {
					return
				}

				pluginExec := common.GetJCLIPluginPath(userHome, plugin.Main, true)

				err = callPluginExecutable(cmd, pluginExec, args, cmd.OutOrStdout())
				return
			},
			// This passes all the flags to the subcommand.
			DisableFlagParsing: true,
		}
		c.Annotations = map[string]string{
			appCfg.ANNOTATION_CONFIG_LOAD: "disable",
		}
		cmd.AddCommand(c)
	}
}

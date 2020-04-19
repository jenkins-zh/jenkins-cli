package cmd

import (
	"fmt"
	"github.com/ghodss/yaml"
	. "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func init() {
	configCmd.AddCommand(configPluginCmd)

	configPluginCmd.AddCommand(NewConfigPluginListCmd(),
		NewConfigPluginFetchCmd(),
		NewConfigPluginInstallCmd())
}

var configPluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: i18n.T("plugin for jcli"),
	Long:  i18n.T("plugin for jcli"),
}

// NewConfigPluginListCmd create a command for list all jcli plugins
func NewConfigPluginListCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "list",
		Short: "list plugins",
		Long:  "list plugins",
		Run: func(cmd *cobra.Command, args []string) {
			for _, plugin := range findPlugins() {
				cmd.Println(plugin)
			}
		},
	}
	return
}

// NewConfigPluginFetchCmd create a command for fetching plugin metadata
func NewConfigPluginFetchCmd() (cmd *cobra.Command) {
	pluginFetchCmd := pluginFetchCmd{}

	cmd = &cobra.Command{
		Use:   "fetch",
		Short: "fetch metadata of plugins",
		Long:  "fetch metadata of plugins",
		RunE:  pluginFetchCmd.Run,
	}

	// add flags
	flags := cmd.Flags()
	flags.StringVarP(&pluginFetchCmd.PluginRepo, "plugin-repo", "", "https://github.com/jenkins-zh/jcli-plugins/",
		i18n.T("The plugin git repository URL"))
	flags.BoolVarP(&pluginFetchCmd.Reset, "reset", "", true,
		i18n.T("If you want to reset the git local repo when pulling it"))
	flags.StringVarP(&pluginFetchCmd.Username, "username", "u", "",
		i18n.T("The username of git repository"))
	flags.StringVarP(&pluginFetchCmd.Password, "password", "p", "",
		i18n.T("The password of git repository"))

	sshKeyFile := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
	flags.StringVarP(&pluginFetchCmd.SSHKeyFile, "ssh-key-file", "", sshKeyFile,
		i18n.T("SSH key file"))
	return
}

// NewConfigPluginInstallCmd create a command for fetching plugin metadata
func NewConfigPluginInstallCmd() (cmd *cobra.Command) {
	pluginInstallCmd := jcliPluginInstallCmd{}

	cmd = &cobra.Command{
		Use:   "install",
		Short: "install a jcli plugins",
		Long:  "install a jcli plugins",
		Args:  cobra.MinimumNArgs(1),
		RunE:  pluginInstallCmd.Run,
	}

	// add flags
	flags := cmd.Flags()
	flags.BoolVarP(&pluginInstallCmd.ShowProgress, "show-progress", "", true,
		i18n.T("If you want to show the progress of download"))
	return
}

type (
	plugin struct {
		Use          string
		Short        string
		Long         string
		Main         string
		Version      string
		DownloadLink string `yaml:"downloadLink"`
	}
	pluginError struct {
		error
		code int
	}
	pluginFetchCmd struct {
		PluginRepo string
		Reset      bool

		Username   string
		Password   string
		SSHKeyFile string

		output io.Writer
	}
	jcliPluginInstallCmd struct {
		RoundTripper http.RoundTripper
		ShowProgress bool

		output io.Writer
	}
)

// Run is the main entry point of plugin fetch command
func (c *pluginFetchCmd) Run(cmd *cobra.Command, args []string) (err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	pluginRepo := fmt.Sprintf("%s/.jenkins-cli/plugins-repo", userHome)
	c.output = cmd.OutOrStdout()

	var r *git.Repository
	if r, err = git.PlainOpen(pluginRepo); err == nil {
		var w *git.Worktree
		if w, err = r.Worktree(); err != nil {
			return
		}

		if c.Reset {
			if err = w.Reset(&git.ResetOptions{
				Mode: git.HardReset,
			}); err != nil {
				return
			}
		}

		err = w.Pull(c.getPullOptions())
		if err == git.NoErrAlreadyUpToDate {
			err = nil // consider it's ok
		}
	} else {
		cloneOptions := c.getCloneOptions()
		_, err = git.PlainClone(pluginRepo, false, cloneOptions)
	}
	return
}

func (c *pluginFetchCmd) getCloneOptions() (cloneOptions *git.CloneOptions) {
	cloneOptions = &git.CloneOptions{
		URL:      c.PluginRepo,
		Progress: c.output,
		Auth:     c.getAuth(),
	}
	return
}

func (c *pluginFetchCmd) getPullOptions() (pullOptions *git.PullOptions) {
	pullOptions = &git.PullOptions{
		RemoteName: "origin",
		Progress:   c.output,
		Auth:       c.getAuth(),
	}
	return
}

func (c *pluginFetchCmd) getAuth() (auth transport.AuthMethod) {
	if c.Username != "" {
		auth = &githttp.BasicAuth{
			Username: c.Username,
			Password: c.Password,
		}
	}

	if strings.HasPrefix(c.PluginRepo, "git@") {
		if sshKey, err := ioutil.ReadFile(c.SSHKeyFile); err == nil {
			signer, _ := ssh.ParsePrivateKey(sshKey)
			auth = &gitssh.PublicKeys{User: "git", Signer: signer}
		}
	}
	return
}

func findPlugins() (plugins []plugin) {
	var userHome string
	var err error
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	plugins = make([]plugin, 0)
	pluginsDir := fmt.Sprintf("%s/.jenkins-cli/plugins/*.yaml", userHome)
	if files, err := filepath.Glob(pluginsDir); err == nil {
		for _, metaFile := range files {
			var data []byte

			plugin := plugin{}
			if data, err = ioutil.ReadFile(metaFile); err == nil {
				if err = yaml.Unmarshal(data, &plugin); err != nil {
					fmt.Println(err)
				} else {
					plugins = append(plugins, plugin)
				}
			}
		}
	}
	return
}

func loadPlugins(cmd *cobra.Command) {
	plugins := findPlugins()
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
						error: errors.Errorf("plugin %q exited with error"),
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

				pluginExec := fmt.Sprintf("%s/.jenkins-cli/plugins/%s", userHome, plugin.Main)

				err = callPluginExecutable(cmd, pluginExec, args, cmd.OutOrStdout())
				return
			},
			// This passes all the flags to the subcommand.
			DisableFlagParsing: true,
		}
		c.Annotations = map[string]string{
			ANNOTATION_CONFIG_LOAD: "disable",
		}
		cmd.AddCommand(c)
	}
}

// Run main entry point for plugin install command
func (c *jcliPluginInstallCmd) Run(cmd *cobra.Command, args []string) (err error) {
	name := args[0]
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	var data []byte
	pluginsMetadataFile := fmt.Sprintf("%s/.jenkins-cli/plugins-repo/%s.yaml", userHome, name)
	if data, err = ioutil.ReadFile(pluginsMetadataFile); err == nil {
		plugin := plugin{}
		if err = yaml.Unmarshal(data, &plugin); err == nil {
			err = c.download(plugin)
		}
	}
	return
}

func (c *jcliPluginInstallCmd) download(plugin plugin) (err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	link := plugin.DownloadLink
	output := fmt.Sprintf("%s/.jenkins-cli/plugins/%s", userHome, plugin.Main)
	logger.Info("start to download plugin",
		zap.String("path", output), zap.String("link", link))

	downloader := util.HTTPDownloader{
		RoundTripper:   c.RoundTripper,
		TargetFilePath: output,
		URL:            link,
		ShowProgress:   c.ShowProgress,
	}
	err = downloader.DownloadFile()
	return
}

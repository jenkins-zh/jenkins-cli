package config

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	appCfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"path/filepath"
	"runtime"
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

// NewConfigPluginListCmd create a command for list all jcli plugins
func NewConfigPluginListCmd(opt *common.Option) (cmd *cobra.Command) {
	configPluginListCmd := configPluginListCmd{
		Option: opt,
	}

	cmd = &cobra.Command{
		Use:   "list",
		Short: "list all installed plugins",
		Long:  "list all installed plugins",
		RunE:  configPluginListCmd.RunE,
		Annotations: map[string]string{
			common.Since: common.VersionSince0028,
		},
	}

	configPluginListCmd.SetFlagWithHeaders(cmd, "Use,Version,DownloadLink")
	return
}

// RunE is the main entry point of config plugin list command
func (c *configPluginListCmd) RunE(cmd *cobra.Command, args []string) (err error) {
	c.Writer = cmd.OutOrStdout()
	err = c.OutputV2(findPlugins())
	return
}

// NewConfigPluginFetchCmd create a command for fetching plugin metadata
func NewConfigPluginFetchCmd(opt *common.Option) (cmd *cobra.Command) {
	pluginFetchCmd := jcliPluginFetchCmd{
		Option: opt,
	}

	cmd = &cobra.Command{
		Use:   "fetch",
		Short: "fetch metadata of plugins",
		Long: `fetch metadata of plugins
The official metadata git repository is https://github.com/jenkins-zh/jcli-plugins,
but you can change it by giving a command parameter.`,
		RunE: pluginFetchCmd.Run,
		Annotations: map[string]string{
			common.Since: common.VersionSince0028,
		},
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
func NewConfigPluginInstallCmd(opt *common.Option) (cmd *cobra.Command) {
	pluginInstallCmd := jcliPluginInstallCmd{
		Option: opt,
	}

	cmd = &cobra.Command{
		Use:   "install",
		Short: "install a jcli plugin",
		Long:  "install a jcli plugin",
		Args:  cobra.MinimumNArgs(1),
		RunE:  pluginInstallCmd.Run,
		Annotations: map[string]string{
			common.Since: common.VersionSince0028,
		},
	}

	// add flags
	flags := cmd.Flags()
	flags.BoolVarP(&pluginInstallCmd.ShowProgress, "show-progress", "", true,
		i18n.T("If you want to show the progress of download"))
	return
}

// Run is the main entry point of plugin fetch command
func (c *jcliPluginFetchCmd) Run(cmd *cobra.Command, args []string) (err error) {
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

func (c *jcliPluginFetchCmd) getCloneOptions() (cloneOptions *git.CloneOptions) {
	cloneOptions = &git.CloneOptions{
		URL:      c.PluginRepo,
		Progress: c.output,
		Auth:     c.getAuth(),
	}
	return
}

func (c *jcliPluginFetchCmd) getPullOptions() (pullOptions *git.PullOptions) {
	pullOptions = &git.PullOptions{
		RemoteName: "origin",
		Progress:   c.output,
		Auth:       c.getAuth(),
	}
	return
}

func (c *jcliPluginFetchCmd) getAuth() (auth transport.AuthMethod) {
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

				pluginExec := fmt.Sprintf("%s/.jenkins-cli/plugins/%s", userHome, plugin.Main)

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

// Run main entry point for plugin install command
func (c *jcliPluginInstallCmd) Run(cmd *cobra.Command, args []string) (err error) {
	name := args[0]
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	var data []byte
	pluginsMetadataFile := fmt.Sprintf("%s/.jenkins-cli/plugins-repo/%s.yaml", userHome, name)
	c.Logger.Info("read plugin metadata info", zap.String("path", pluginsMetadataFile))
	if data, err = ioutil.ReadFile(pluginsMetadataFile); err == nil {
		plugin := plugin{}
		if err = yaml.Unmarshal(data, &plugin); err == nil {
			err = c.download(plugin)
		}
	}

	if err == nil {
		cachedMetadataFile := fmt.Sprintf("%s/.jenkins-cli/pluginss/%s.yaml", userHome, name)
		err = ioutil.WriteFile(cachedMetadataFile, data, 0664)
	}
	return
}

func (c *jcliPluginInstallCmd) download(plugin plugin) (err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	link := c.getDownloadLink(plugin)
	output := fmt.Sprintf("%s/.jenkins-cli/plugins/%s.tar.gz", userHome, plugin.Main)
	c.Logger.Info("start to download plugin",
		zap.String("path", output), zap.String("link", link))

	downloader := util.HTTPDownloader{
		RoundTripper:   c.RoundTripper,
		TargetFilePath: output,
		URL:            link,
		ShowProgress:   c.ShowProgress,
	}
	if err = downloader.DownloadFile(); err == nil {
		c.Logger.Info("start to extract files")
		err = c.extractFiles(plugin, output)
	}
	return
}

func (c *jcliPluginInstallCmd) getDownloadLink(plugin plugin) (link string) {
	link = plugin.DownloadLink
	if link == "" {
		operationSystem := runtime.GOOS
		arch := runtime.GOARCH
		link = fmt.Sprintf("https://github.com/jenkins-zh/%s/releases/download/%s/%s-%s-%s.tar.gz",
			plugin.Main, plugin.Version, plugin.Main, operationSystem, arch)
	}
	return
}

func (c *jcliPluginInstallCmd) extractFiles(plugin plugin, tarFile string) (err error) {
	var f *os.File
	var gzf *gzip.Reader
	if f, err = os.Open(tarFile); err != nil {
		c.Logger.Error("open file error", zap.String("path", tarFile))
		return
	}
	defer f.Close()

	if gzf, err = gzip.NewReader(f); err != nil {
		c.Logger.Error("open tar file error", zap.String("path", tarFile))
		return
	}

	tarReader := tar.NewReader(gzf)
	var header *tar.Header
	for {
		if header, err = tarReader.Next(); err == io.EOF {
			c.Logger.Info("extracted all files")
			err = nil
			break
		} else if err != nil {
			c.Logger.Error("tar file reading error")
			break
		}
		name := header.Name

		switch header.Typeflag {
		case tar.TypeReg:
			if name != plugin.Main {
				c.Logger.Debug("ignore file", zap.String("name", name))
				continue
			}
			var targetFile *os.File
			if targetFile, err = os.OpenFile(fmt.Sprintf("%s/%s", filepath.Dir(tarFile), name),
				os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode)); err != nil {
				break
			}
			c.Logger.Info("extracting file", zap.String("path", targetFile.Name()))
			if _, err = io.Copy(targetFile, tarReader); err != nil {
				break
			}
			targetFile.Close()
		default:
			c.Logger.Debug("ignore this type from tar file",
				zap.Int32("type", int32(header.Typeflag)), zap.String("name", name))
		}
	}
	return
}

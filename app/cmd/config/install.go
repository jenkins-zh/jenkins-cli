package config

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// NewConfigPluginInstallCmd create a command for fetching plugin metadata
func NewConfigPluginInstallCmd(opt *common.Option) (cmd *cobra.Command) {
	pluginInstallCmd := jcliPluginInstallCmd{
		Option: opt,
	}

	cmd = &cobra.Command{
		Use:               "install",
		Short:             "install a jcli plugin",
		Long:              "install a jcli plugin",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: ValidPluginNames,
		RunE:              pluginInstallCmd.Run,
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
		cachedMetadataFile := common.GetJCLIPluginPath(userHome, name, true)
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

package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// NewSelfUpgradeCmd create a command for self upgrade
func NewSelfUpgradeCmd(client common.JenkinsClient, jenkinsConfigMgr common.JenkinsConfigMgr) (cmd *cobra.Command) {
	opt := &SelfUpgradeOption{}

	cmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade jcli it self",
		Long:  `Upgrade jcli it self`,
		RunE: opt.RunE,
		Annotations: map[string]string{
			common.Since: "v0.0.26",
		},
	}
	opt.addFlags(cmd.Flags())
	return
}

func (o *SelfUpgradeOption) addFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&o.ShowProgress, "show-progress", "", true,
		i18n.T("If you want to show the progress of download Jenkins CLI"))
}

func (o *SelfUpgradeOption) RunE(cmd *cobra.Command, args []string) (err error) {
	var version string
	if len(args) > 0 {
		version = args[0]
	}

	switch(version) {
	case "", "dev":
		version = "master"
	}

	output := "jcli.tar.gz"
	fileURL := fmt.Sprintf("https://cdn.jsdelivr.net/gh/jenkins-zh/jcli-repo@%s/jcli-%s-amd64.tar.gz",
		version, runtime.GOOS)

	downloader := util.HTTPDownloader{
		RoundTripper:   o.RoundTripper,
		TargetFilePath: output,
		URL:            fileURL,
		ShowProgress:   o.ShowProgress,
	}
	err = downloader.DownloadFile()

	// copy binary file into system path
	var targetPath string
	if targetPath, err = exec.LookPath("jcli"); err != nil {
		return
	}

	if err = o.extractFiles(output); err == nil {
		targetF, _ := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0644)
		sourceF, _ := os.Open("jcli")
		_, err = io.Copy(targetF, sourceF)
	}
	return
}

func (c *SelfUpgradeOption) extractFiles(tarFile string) (err error) {
	var f *os.File
	var gzf *gzip.Reader
	if f, err = os.Open(tarFile); err != nil {
		return
	}
	defer f.Close()

	if gzf, err = gzip.NewReader(f); err != nil {
		return
	}

	tarReader := tar.NewReader(gzf)
	var header *tar.Header
	for {
		if header, err = tarReader.Next(); err == io.EOF {
			err = nil
			break
		} else if err != nil {
			break
		}
		name := header.Name

		switch header.Typeflag {
		case tar.TypeReg:
			if name != "jcli" {
				continue
			}
			var targetFile *os.File
			if targetFile, err = os.OpenFile(fmt.Sprintf("%s/%s", filepath.Dir(tarFile), name),
				os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode)); err != nil {
				break
			}
			if _, err = io.Copy(targetFile, tarReader); err != nil {
				break
			}
			targetFile.Close()
		}
	}
	return
}

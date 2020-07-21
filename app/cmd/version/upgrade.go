package version

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/google/go-github/v29/github"
	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"io/ioutil"
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
		Short: "Upgrade jcli itself",
		Long: `Upgrade jcli itself
You can use any exists version to upgrade jcli itself. If there's no argument given, it will upgrade to the latest release.
You can upgrade to the latest developing version, please use it like: jcli version upgrade dev'`,
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

// RunE is the main point of current command
func (o *SelfUpgradeOption) RunE(cmd *cobra.Command, args []string) (err error) {
	var version string
	if len(args) > 0 {
		version = args[0]
	}

	// copy binary file into system path
	var targetPath string
	if targetPath, err = exec.LookPath("jcli"); err != nil {
		err = fmt.Errorf("cannot find Jenkins CLI from system path, error: %v", err)
		return
	}
	var targetF *os.File
	if targetF, err = os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return
	}

	// try to understand the version from user input
	switch version {
	case "dev":
		version = "master"
	case "":
		o.GitHubClient = github.NewClient(nil)
		ghClient := &client.GitHubReleaseClient{
			Client: o.GitHubClient,
		}
		if asset, assetErr := ghClient.GetLatestJCLIAsset(); assetErr == nil && asset != nil {
			version = asset.TagName
		} else {
			err = fmt.Errorf("cannot get the latest version, error: %s", assetErr)
			return
		}
	}

	// version review
	currentVersion := app.GetVersion()
	if currentVersion == version {
		cmd.Println("no need to upgrade Jenkins CLI")
		return
	}
	cmd.Println(fmt.Sprintf("prepare to upgrade to %s", version))

	// download the tar file of Jenkins CLI
	tmpDir := os.TempDir()
	output := fmt.Sprintf("%s/jcli.tar.gz", tmpDir)
	fileURL := fmt.Sprintf("https://cdn.jsdelivr.net/gh/jenkins-zh/jcli-repo@%s/jcli-%s-amd64.tar.gz",
		version, runtime.GOOS)
	defer func() {
		_ = os.RemoveAll(output)
	}()

	// make sure we count the download action
	go func() {
		o.downloadCount(version, runtime.GOOS)
	}()

	downloader := util.HTTPDownloader{
		RoundTripper:   o.RoundTripper,
		TargetFilePath: output,
		URL:            fileURL,
		ShowProgress:   o.ShowProgress,
	}
	if err = downloader.DownloadFile(); err != nil {
		err = fmt.Errorf("cannot download Jenkins CLI from %s, error: %v", fileURL, err)
		return
	}

	if err = o.extractFiles(output); err == nil {
		sourceFile := fmt.Sprintf("%s/jcli", filepath.Dir(output))
		sourceF, _ := os.Open(sourceFile)
		if _, err = io.Copy(targetF, sourceF); err != nil {
			err = fmt.Errorf("cannot copy Jenkins CLI from %s to %s, error: %v", sourceFile, targetPath, err)
		}
	} else {
		err = fmt.Errorf("cannot extract Jenkins CLI from tar file, error: %v", err)
	}
	return
}

func (o *SelfUpgradeOption) downloadCount(version string, arch string) {
	countURL := fmt.Sprintf("https: //github.com/jenkins-zh/jenkins-cli/releases/download/v%s/jcli-%s-amd64.tar.gz",
		version, arch)

	if tempDir, err := ioutil.TempDir(".", "download-count"); err == nil {
		tempFile := tempDir + "/jcli.tar.gz"
		defer func() {
			_ = os.RemoveAll(tempDir)
		}()

		downloader := util.HTTPDownloader{
			RoundTripper:   o.RoundTripper,
			TargetFilePath: tempFile,
			URL:            countURL,
		}
		// we don't care about the result, just for counting
		_ = downloader.DownloadFile()
	}
}

func (o *SelfUpgradeOption) extractFiles(tarFile string) (err error) {
	var f *os.File
	var gzf *gzip.Reader
	if f, err = os.Open(tarFile); err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()

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
			_ = targetFile.Close()
		}
	}
	return
}

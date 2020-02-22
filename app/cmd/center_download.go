package cmd

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterDownloadOption as the options of download command
type CenterDownloadOption struct {
	LTS     bool
	Mirror  string
	Version string

	Output       string
	ShowProgress bool

	RoundTripper http.RoundTripper
}

var centerDownloadOption CenterDownloadOption

func init() {
	centerCmd.AddCommand(centerDownloadCmd)
	centerDownloadCmd.Flags().BoolVarP(&centerDownloadOption.LTS, "lts", "", true,
		i18n.T("If you want to download Jenkins as LTS"))
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Version, "war-version", "", "",
		i18n.T("Version of the Jenkins which you want to download"))
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Mirror, "mirror", "m", "default",
		i18n.T("The mirror site of Jenkins"))
	centerDownloadCmd.Flags().BoolVarP(&centerDownloadOption.ShowProgress, "progress", "p", true,
		i18n.T("If you want to show the download progress"))
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Output, "output", "o", "jenkins.war",
		i18n.T("The file of output"))
}

var centerDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: i18n.T("Download Jenkins"),
	Long:  i18n.T(`Download Jenkins from a mirror site. You can get more mirror sites from https://jenkins-zh.cn/tutorial/management/mirror/`),
	RunE: func(cmd *cobra.Command, _ []string) error {
		return centerDownloadOption.DownloadJenkins()
	},
}

// DownloadJenkins download the Jenkins
func (c *CenterDownloadOption) DownloadJenkins() (err error) {
	parentDir := filepath.Dir(c.Output)
	if err = os.MkdirAll(parentDir, os.FileMode(0755)); err != nil {
		return
	}

	mirrorSite := getMirror(c.Mirror)
	if mirrorSite == "" {
		err = fmt.Errorf("cannot found Jenkins mirror by: %s", c.Mirror)
		return
	}

	jClient := &client.UpdateCenterManager{
		MirrorSite: mirrorSite,
		JenkinsCore: client.JenkinsCore{
			RoundTripper: c.RoundTripper,
		},
		LTS:          c.LTS,
		Version:      c.Version,
		Output:       c.Output,
		ShowProgress: c.ShowProgress,
	}

	jenkinsWarURL := jClient.GetJenkinsWarURL()
	logger.Info("prepare to download jenkins.war", zap.String("URL", jenkinsWarURL))

	err = jClient.DownloadJenkins()
	return
}

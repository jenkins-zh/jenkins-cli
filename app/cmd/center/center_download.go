package center

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"os"
	"path/filepath"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func NewCenterDownloadCmd(client common.JenkinsClient, jenkinsConfigMgr common.JenkinsConfigMgr) (cmd *cobra.Command) {
	opt := &CenterDownloadOption{
		JenkinsClient:    client,
		JenkinsConfigMgr: jenkinsConfigMgr,
	}
	cmd = &cobra.Command{
		Use:   "download",
		Short: i18n.T("Download jenkins.war"),
		Long: i18n.T(`Download jenkins.war from a mirror site.
You can get more mirror sites from https://jenkins-zh.cn/tutorial/management/mirror/
If you want to download different formulas of jenkins.war, please visit the following project
https://github.com/jenkins-zh/docker-zh`),
		RunE:    opt.RunE,
		PreRunE: opt.PreRunE,
	}

	cmd.Flags().BoolVarP(&opt.LTS, "lts", "", true,
		i18n.T("If you want to download Jenkins as LTS"))
	cmd.Flags().StringVarP(&opt.Version, "war-version", "", "",
		i18n.T("Version of the Jenkins which you want to download"))
	cmd.Flags().StringVarP(&opt.Mirror, "mirror", "m", "default",
		i18n.T("The mirror site of Jenkins"))
	cmd.Flags().BoolVarP(&opt.ShowProgress, "progress", "p", true,
		i18n.T("If you want to show the download progress"))
	cmd.Flags().StringVarP(&opt.Output, "output", "o", "",
		i18n.T("The file of output"))
	cmd.Flags().StringVarP(&opt.Formula, "formula", "", "",
		i18n.T("The formula of jenkins.war, only support zh currently"))
	return
}

func (o *CenterDownloadOption) RunE(cmd *cobra.Command, _ []string) error {
	return o.DownloadJenkins()
}

func (o *CenterDownloadOption) PreRunE(cmd *cobra.Command, args []string) (err error) {
	if o.Output != "" {
		return
	}

	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	o.Output = fmt.Sprintf("%s/.jenkins-cli/cache/%s/jenkins.war", userHome, o.Version)
	return
}

// DownloadJenkins download the Jenkins
func (c *CenterDownloadOption) DownloadJenkins() (err error) {
	parentDir := filepath.Dir(c.Output)
	if err = os.MkdirAll(parentDir, os.FileMode(0755)); err != nil {
		return
	}

	mirrorSite := c.JenkinsConfigMgr.GetMirror(c.Mirror)
	if mirrorSite == "" {
		err = fmt.Errorf("cannot found Jenkins mirror by: %s", c.Mirror)
		return
	}

	jClient := &client.UpdateCenterManager{
		MirrorSite: mirrorSite,
		JenkinsCore: client.JenkinsCore{
			RoundTripper: c.RoundTripper,
		},
		Formula:      c.Formula,
		LTS:          c.LTS,
		Version:      c.Version,
		Output:       c.Output,
		ShowProgress: c.ShowProgress,
	}

	jenkinsWarURL := jClient.GetJenkinsWarURL()
	c.Logger.Info("prepare to download jenkins.war", zap.String("URL", jenkinsWarURL))

	err = jClient.DownloadJenkins()
	return
}

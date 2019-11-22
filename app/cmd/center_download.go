package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"
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
	Run: func(cmd *cobra.Command, _ []string) {
		mirrorSite := getMirror(centerDownloadOption.Mirror)
		if mirrorSite == "" {
			cmd.PrintErrln("cannot found Jenkins mirror by:", centerDownloadOption.Mirror)
			return
		}

		jClient := &client.UpdateCenterManager{
			MirrorSite: mirrorSite,
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerDownloadOption.RoundTripper,
			},
			LTS:          centerDownloadOption.LTS,
			Version:      centerDownloadOption.Version,
			Output:       centerDownloadOption.Output,
			ShowProgress: centerDownloadOption.ShowProgress,
		}

		err := jClient.DownloadJenkins()
		helper.CheckErr(cmd, err)
	},
}

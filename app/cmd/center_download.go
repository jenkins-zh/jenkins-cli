package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterDownloadOption as the options of download command
type CenterDownloadOption struct {
	LTS    bool
	Mirror string

	Output       string
	ShowProgress bool

	RoundTripper http.RoundTripper
}

var centerDownloadOption CenterDownloadOption

func init() {
	centerCmd.AddCommand(centerDownloadCmd)
	centerDownloadCmd.Flags().BoolVarP(&centerDownloadOption.LTS, "lts", "", true, "If you want to download Jenkins as LTS")
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Mirror, "mirror", "m", "default", "The mirror site of Jenkins")
	centerDownloadCmd.Flags().BoolVarP(&centerDownloadOption.ShowProgress, "progress", "p", true, "If you want to show the download progress")
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Output, "output", "o", "jenkins.war", "The file of output")
}

var centerDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Jenkins",
	Long:  `Download Jenkins from a mirror site. You can get more mirror sites from https://jenkins-zh.cn/tutorial/management/mirror/`,
	Run: func(cmd *cobra.Command, _ []string) {
		config := getConfig()
		mirrorSite := centerDownloadOption.getMirrorSite(config)
		if mirrorSite == "" {
			cmd.PrintErrln("cannot found Jenkins mirror by:", centerDownloadOption.Mirror)
			return
		}

		jclient := &client.UpdateCenterManager{
			MirrorSite: mirrorSite,
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerDownloadOption.RoundTripper,
			},
		}

		err := jclient.DownloadJenkins(centerDownloadOption.LTS, centerDownloadOption.ShowProgress,
			centerDownloadOption.Output)
		helper.CheckErr(cmd, err)
	},
}

func (c *CenterDownloadOption) getMirrorSite(config *Config) (site string) {
	mirrors := getMirrors()
	for _, mirror := range mirrors {
		if mirror.Name == c.Mirror {
			site = mirror.URL
			return
		}
	}
	return
}

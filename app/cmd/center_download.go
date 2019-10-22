package cmd

import (
	"log"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterDownloadOption as the options of download command
type CenterDownloadOption struct {
	LTS    bool
	Output string
	ShowProgress bool

	RoundTripper http.RoundTripper
}

var centerDownloadOption CenterDownloadOption

func init() {
	centerCmd.AddCommand(centerDownloadCmd)
	centerDownloadCmd.Flags().BoolVarP(&centerDownloadOption.LTS, "lts", "", true, "If you want to download Jenkins as LTS")
	centerDownloadCmd.Flags().BoolVarP(&centerDownloadOption.ShowProgress, "progress", "", true, "If you want to show the download progress")
	centerDownloadCmd.Flags().StringVarP(&centerDownloadOption.Output, "output", "o", "jenkins.war", "The file of output")
}

var centerDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Jenkins",
	Long:  `Download Jenkins`,
	Run: func(_ *cobra.Command, _ []string) {
		jclient := &client.UpdateCenterManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerDownloadOption.RoundTripper,
			},
		}

		if err := jclient.DownloadJenkins(centerDownloadOption.LTS, centerDownloadOption.ShowProgress,
			centerDownloadOption.Output); err != nil {
			log.Fatal(err)
		}
	},
}

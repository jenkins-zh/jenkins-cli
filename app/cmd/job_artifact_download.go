package cmd

import (
	"fmt"
	appCfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg/net"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/helper"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobArtifactDownloadOption is the options of job artifact download command
type JobArtifactDownloadOption struct {
	ID           string
	ShowProgress bool
	DownloadDir  string

	Jenkins      *appCfg.JenkinsServer
	RoundTripper http.RoundTripper
}

var jobArtifactDownloadOption JobArtifactDownloadOption

func init() {
	jobArtifactCmd.AddCommand(jobArtifactDownloadCmd)
	jobArtifactDownloadCmd.Flags().StringVarP(&jobArtifactDownloadOption.ID, "id", "i", "",
		i18n.T("ID of the job artifact"))
	jobArtifactDownloadCmd.Flags().BoolVarP(&jobArtifactDownloadOption.ShowProgress, "progress", "", true,
		i18n.T("Whether show the progress"))
	jobArtifactDownloadCmd.Flags().StringVarP(&jobArtifactDownloadOption.DownloadDir, "download-dir", "", "",
		i18n.T("The directory which artifact will be downloaded"))
}

var jobArtifactDownloadCmd = &cobra.Command{
	Use:   "download <jobName> [buildID]",
	Short: i18n.T("Download the artifact of target job"),
	Long:  i18n.T(`Download the artifact of target job`),
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		argLen := len(args)
		var err error
		jobName := args[0]
		buildID := -1

		if argLen >= 2 {
			if buildID, err = strconv.Atoi(args[1]); err != nil {
				cmd.PrintErrln(err)
				return
			}
		}

		jclient := &client.ArtifactClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobArtifactDownloadOption.RoundTripper,
			},
		}
		jobArtifactDownloadOption.Jenkins = getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))

		artifacts, err := jclient.List(jobName, buildID)
		if err == nil {
			for _, artifact := range artifacts {
				if jobArtifactDownloadOption.ID != "" && jobArtifactDownloadOption.ID != artifact.ID {
					continue
				}

				downloadPath := filepath.Join(jobArtifactDownloadOption.DownloadDir, artifact.Name)
				err = jobArtifactDownloadOption.download(artifact.URL, downloadPath)
				if err != nil {
					break
				}
			}
		}
		helper.CheckErr(cmd, err)
	},
}

func (j *JobArtifactDownloadOption) download(artifactURL, fileName string) (err error) {
	jenkinsURL, _ := url.Parse(j.Jenkins.URL)
	targetURL := fmt.Sprintf("%s%s", j.Jenkins.URL, strings.TrimPrefix(artifactURL, jenkinsURL.Path))
	fmt.Println("start to download from", targetURL)
	//downloader := httpdownloader.HTTPDownloader{
	//	RoundTripper:   j.RoundTripper,
	//	TargetFilePath: fileName,
	//	URL:            targetURL,
	//	UserName:       j.Jenkins.UserName,
	//	Password:       j.Jenkins.Token,
	//	Proxy:          j.Jenkins.Proxy,
	//	ProxyAuth:      j.Jenkins.ProxyAuth,
	//	ShowProgress:   j.ShowProgress,
	//	Thread:         10,
	//}
	//err = downloader.DownloadFile()

	download := &httpdownloader.MultiThreadDownloader{}
	download.WithBasicAuth(j.Jenkins.UserName, j.Jenkins.Token)
	download.WithShowProgress(j.ShowProgress)
	download.WithKeepParts(true)
	download.WithRoundTripper(j.RoundTripper)
	err = download.Download(targetURL, fileName, runtime.NumCPU())
	return
}

package cmd

import (
	"fmt"
	"strconv"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// JobArtifactDownloadOption is the options of job artifact download command
type JobArtifactDownloadOption struct {
	ID string

	Jenkins *JenkinsServer
	RoundTripper http.RoundTripper
}

var jobArtifactDownloadOption JobArtifactDownloadOption

func init() {
	jobArtifactCmd.AddCommand(jobArtifactDownloadCmd)
	jobArtifactDownloadCmd.Flags().StringVarP(&jobArtifactDownloadOption.ID, "id", "i", "", "ID of the job artifact")
}

var jobArtifactDownloadCmd = &cobra.Command{
	Use:   "download <jobName> [buildID]",
	Short: "Download the artifact of target job",
	Long:  `Download the artifact of target job`,
	Run: func(cmd *cobra.Command, args []string) {
		argLen := len(args)
		if argLen == 0 {
			cmd.Help()
			return
		}

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
		jobArtifactDownloadOption.Jenkins = getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if artifacts, err := jclient.List(jobName, buildID); err == nil {
			for _, artifact := range artifacts {
				if jobArtifactDownloadOption.ID != "" && jobArtifactDownloadOption.ID != artifact.ID {
					continue
				}

				err = jobArtifactDownloadOption.download(artifact.URL, artifact.Name)
				if err != nil {
					cmd.PrintErrln(err)
				}
			}
		} else {
			cmd.PrintErrln(err)
		}
	},
}

func (j *JobArtifactDownloadOption) download(url, fileName string) (err error) {
	downloader := util.HTTPDownloader{
		RoundTripper:   j.RoundTripper,
		TargetFilePath: fileName,
		URL:            fmt.Sprintf("%s/%s", j.Jenkins.URL, url),
		UserName: j.Jenkins.UserName,
		Password: j.Jenkins.Token,
		Proxy: j.Jenkins.Proxy,
		ProxyAuth: j.Jenkins.ProxyAuth,
		ShowProgress:   true,
	}
	err = downloader.DownloadFile()
	return
}

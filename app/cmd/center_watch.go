package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"net/http"
	"time"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// CenterWatchOption as the options of watch command
type CenterWatchOption struct {
	common.WatchOption
	UtilNeedRestart     bool
	UtilInstallComplete bool

	RoundTripper  http.RoundTripper
	CeneterStatus string
}

var centerWatchOption CenterWatchOption

func init() {
	centerCmd.AddCommand(centerWatchCmd)
	centerWatchCmd.Flags().BoolVarP(&centerWatchOption.UtilNeedRestart, "util-need-restart", "", false,
		i18n.T("The watch will be continue util Jenkins needs restart"))
	centerWatchCmd.Flags().BoolVarP(&centerWatchOption.UtilInstallComplete, "util-install-complete", "", false,
		i18n.T("The watch will be continue util all Jenkins plugins installation is completed"))
	centerWatchOption.SetFlag(centerWatchCmd)
}

var centerWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch your update center status",
	Long:  `Watch your update center status`,
	Run: func(cmd *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		printJenkinsStatus(jenkins, cmd, centerWatchOption.RoundTripper)

		for ; centerWatchOption.Count >= 0; centerWatchOption.Count-- {
			if status, err := printUpdateCenter(jenkins, cmd, centerOption.RoundTripper); err != nil {
				cmd.PrintErr(err)
				break
			} else if (centerWatchOption.UtilNeedRestart && status.RestartRequiredForCompletion) ||
				(centerWatchOption.UtilInstallComplete && allPluginsCompleted(status)) {
				return
			}

			time.Sleep(time.Duration(centerOption.Interval) * time.Second)
		}
	},
}

func allPluginsCompleted(status *client.UpdateCenter) (completed bool) {
	if status == nil || status.Jobs == nil {
		return
	}

	for _, job := range status.Jobs {
		if job.Type == "InstallationJob" && !job.Status.Success {
			return
		}
	}
	completed = true
	return
}

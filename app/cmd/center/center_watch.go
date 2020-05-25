package center

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"time"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func NewCenterWatchCmd(jenkinsClient common.JenkinsClient, centerOption *CenterOption) (cmd *cobra.Command) {
	opt := &CenterWatchOption{
		JenkinsClient: jenkinsClient,
		CenterOption:  centerOption,
	}
	cmd = &cobra.Command{
		Use:   "watch",
		Short: "Watch your update center status",
		Long:  `Watch your update center status`,
		Run:   opt.RunE,
	}

	cmd.Flags().BoolVarP(&opt.UtilNeedRestart, "util-need-restart", "", false,
		i18n.T("The watch will be continue util Jenkins needs restart"))
	cmd.Flags().BoolVarP(&opt.UtilInstallComplete, "util-install-complete", "", false,
		i18n.T("The watch will be continue util all Jenkins plugins installation is completed"))
	opt.SetFlag(cmd)
	return
}

func (o *CenterWatchOption) RunE(cmd *cobra.Command, _ []string) {
	jenkins := o.JenkinsClient.GetCurrentJenkinsFromOptions()
	o.printJenkinsStatus(jenkins, cmd, o.RoundTripper)

	for ; o.Count >= 0; o.Count-- {
		if status, err := o.printUpdateCenter(jenkins, cmd, o.RoundTripper); err != nil {
			cmd.PrintErr(err)
			break
		} else if (o.UtilNeedRestart && status.RestartRequiredForCompletion) ||
			(o.UtilInstallComplete && o.allPluginsCompleted(status)) {
			return
		}

		time.Sleep(time.Duration(o.Interval) * time.Second)
	}
}

func (o *CenterWatchOption) allPluginsCompleted(status *client.UpdateCenter) (completed bool) {
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

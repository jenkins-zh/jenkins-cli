package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerLaunchOption option for config list command
type ComputerLaunchOption struct {
	CommonOption
}

var computerLaunchOption ComputerLaunchOption

func init() {
	computerCmd.AddCommand(computerLaunchCmd)
}

var computerLaunchCmd = &cobra.Command{
	Use:     "launch <name>",
	Aliases: []string{"start"},
	Short:   i18n.T("Launch the agent of your Jenkins"),
	Long:    i18n.T("Launch the agent of your Jenkins"),
	Args:    cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) (err error) {
		jClient := &client.ComputerClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: computerLaunchOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		err = jClient.Launch(args[0])
		return
		//	/computer/nginx/toggleOffline
		//	json: {"offlineMessage": "sdd", "Jenkins-Crumb": "c1482407700d5edfeca3d78315fe7ca33ba89caaaf55ae6b3c6f351fcc2f5470"}

		//	/computer/nginx/changeOfflineCause
		//json: {"offlineMessage": "被 jenkins : sdd断开连接", "Jenkins-Crumb": "c1482407700d5edfeca3d78315fe7ca33ba89caaaf55ae6b3c6f351fcc2f5470"}
	},
}

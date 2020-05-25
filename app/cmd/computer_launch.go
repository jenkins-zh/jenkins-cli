package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/center"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	. "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"go.uber.org/zap"
	"io/ioutil"
	"os"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// ComputerLaunchOption option for config list command
type ComputerLaunchOption struct {
	common.CommonOption

	Type         string
	ShowProgress bool

	/** share info between inner functions */
	ComputerClient *client.ComputerClient
	CurrentJenkins *JenkinsServer
	Output         string
}

var computerLaunchOption ComputerLaunchOption

func init() {
	computerCmd.AddCommand(computerLaunchCmd)
	computerLaunchCmd.Flags().StringVarP(&computerLaunchOption.Type, "type", "", "",
		i18n.T("The type of agent, include jnlp"))
	computerLaunchCmd.Flags().BoolVarP(&computerLaunchOption.ShowProgress, "show-progress", "", true,
		i18n.T("Show the progress of downloading agent.jar"))
}

var computerLaunchCmd = &cobra.Command{
	Use:     "launch",
	Aliases: []string{"start"},
	Short:   i18n.T("Launch the agent of your Jenkins"),
	Long:    i18n.T("Launch the agent of your Jenkins"),
	Args:    cobra.MinimumNArgs(1),
	Example: `jcli agent launch agent-name
jcli agent launch agent-name --type jnlp`,
	PreRunE: func(_ *cobra.Command, args []string) (err error) {
		computerLaunchOption.ComputerClient, computerLaunchOption.CurrentJenkins =
			GetComputerClient(computerLaunchOption.CommonOption)

		if computerLaunchOption.Type != "jnlp" {
			return
		}

		var f *os.File
		if f, err = ioutil.TempFile("/tmp", "agent.jar"); err == nil {
			computerLaunchOption.Output = f.Name()
			downloader := util.HTTPDownloader{
				RoundTripper:   computerLaunchOption.RoundTripper,
				TargetFilePath: computerLaunchOption.Output,
				URL:            fmt.Sprintf("%s/jnlpJars/agent.jar", computerLaunchOption.ComputerClient.URL),
				ShowProgress:   computerLaunchOption.ShowProgress,
			}
			err = downloader.DownloadFile()
		}
		return
	},
	RunE: func(_ *cobra.Command, args []string) (err error) {
		name := args[0]
		switch computerLaunchOption.Type {
		case "":
			err = computerLaunchOption.Launch(name)
		case "jnlp":
			err = computerLaunchOption.LaunchJnlp(name)
		}
		return
	},
}

// Launch start a normal agent
func (o *ComputerLaunchOption) Launch(name string) (err error) {
	err = o.ComputerClient.Launch(name)
	return
}

// LaunchJnlp start a JNLP agent
func (o *ComputerLaunchOption) LaunchJnlp(name string) (err error) {
	var secret string
	if secret, err = o.ComputerClient.GetSecret(name); err == nil {
		var binary string
		binary, err = util.LookPath("java", center.centerStartOption.LookPathContext)
		if err == nil {
			env := os.Environ()
			agentArgs := []string{"java", "-jar", computerLaunchOption.Output,
				"-jnlpUrl", fmt.Sprintf("%s/computer/%s/slave-agent.jnlp", o.ComputerClient.URL, name),
				"-secret", secret, "-workDir", "/tmp"}

			logger.Debug("start a jnlp agent", zap.Any("command", agentArgs))

			err = util.Exec(binary, agentArgs, env, o.SystemCallExec)
		}
	}
	return
}

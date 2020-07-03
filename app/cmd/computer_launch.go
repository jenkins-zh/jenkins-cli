package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	. "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"go.uber.org/zap"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

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

const (
	AGENT_JNLP = "jnlp"
)

var computerLaunchOption ComputerLaunchOption

func init() {
	computerCmd.AddCommand(computerLaunchCmd)
	computerLaunchCmd.Flags().StringVarP(&computerLaunchOption.Type, "type", "", AGENT_JNLP,
		i18n.T("The type of agent, include jnlp"))
	computerLaunchCmd.Flags().BoolVarP(&computerLaunchOption.ShowProgress, "show-progress", "", true,
		i18n.T("Show the progress of downloading agent.jar"))

	healthCheckRegister.Register(getCmdPath(computerLaunchCmd), &computerLaunchOption)
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

		if computerLaunchOption.Type != AGENT_JNLP {
			return
		}

		var f *os.File
		if f, err = ioutil.TempFile("/tmp", "agent.jar"); err == nil {
			computerLaunchOption.Output = f.Name()
			agentURL := fmt.Sprintf("%s/jnlpJars/agent.jar", computerLaunchOption.ComputerClient.URL)
			logger.Debug("start to download agent.jar", zap.String("url", agentURL))
			logger.Debug("proxy setting", zap.String("sever", computerLaunchOption.CurrentJenkins.Proxy),
				zap.String("auth", computerLaunchOption.CurrentJenkins.ProxyAuth))

			downloader := util.HTTPDownloader{
				RoundTripper:   computerLaunchOption.RoundTripper,
				TargetFilePath: computerLaunchOption.Output,
				URL:            agentURL,
				Proxy:          computerLaunchOption.CurrentJenkins.Proxy,
				ProxyAuth:      computerLaunchOption.CurrentJenkins.ProxyAuth,
				ShowProgress:   computerLaunchOption.ShowProgress,
			}
			err = downloader.DownloadFile()
		}
		return
	},
	RunE: func(_ *cobra.Command, args []string) (err error) {
		name := args[0]
		logger.Info("prepare to start agent", zap.String("name", name), zap.String("type", computerLaunchOption.Type))

		switch computerLaunchOption.Type {
		case "":
			err = computerLaunchOption.Launch(name)
		case AGENT_JNLP:
			err = computerLaunchOption.LaunchJnlp(name)
		default:
			err = fmt.Errorf("unsupported agent type %s", computerLaunchOption.Type)
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
		logger.Info("get agent secret", zap.String("value", secret))

		var binary string
		binary, err = util.LookPath("java", centerStartOption.LookPathContext)
		if err == nil {
			env := os.Environ()
			agentArgs := []string{"java", "-jar", computerLaunchOption.Output,
				"-jnlpUrl", fmt.Sprintf("%s/computer/%s/slave-agent.jnlp", o.ComputerClient.URL, name),
				"-secret", secret, "-workDir", "/tmp"}

			if o.CurrentJenkins.ProxyAuth != "" {
				proxyAuth := strings.SplitN(o.CurrentJenkins.ProxyAuth, ":", 2)

				proxyURL, _ := url.Parse(o.CurrentJenkins.Proxy)
				if len(proxyAuth) == 2 {
					env = append(env, fmt.Sprintf("http_proxy=http://%s:%s@%s", url.QueryEscape(proxyAuth[0]), url.QueryEscape(proxyAuth[1]), proxyURL.Host))
				}
			}

			logger.Debug("start a jnlp agent", zap.Any("command", agentArgs))

			err = util.Exec(binary, agentArgs, env, o.SystemCallExec)
		}
	}
	return
}

// Check do the health check of casc cmd
func (o *ComputerLaunchOption) Check() (err error) {
	opt := PluginOptions{
		CommonOption: common.CommonOption{RoundTripper: o.RoundTripper},
	}
	_, err = opt.FindPlugin("pipeline-restful-api")
	return
}

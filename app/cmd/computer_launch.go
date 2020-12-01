package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	appCfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

// ComputerLaunchOption option for config list command
type ComputerLaunchOption struct {
	common.Option

	Type         string
	ShowProgress bool

	/** share info between inner functions */
	ComputerClient *client.ComputerClient
	CurrentJenkins *appCfg.JenkinsServer
	Output         string

	Mode      string
	Remove    bool
	Restart   string
	Detach    bool
	AgentType string
}

const (
	// AgentJNLP is the agent type of jnlp
	AgentJNLP = "jnlp"
)

var computerLaunchOption ComputerLaunchOption

func init() {
	computerCmd.AddCommand(computerLaunchCmd)
	flags := computerLaunchCmd.Flags()

	flags.StringVarP(&computerLaunchOption.Mode, "mode", "m", "java",
		i18n.T("Mode of launching Jenkins, you can choose: java, docker"))
	flags.BoolVarP(&computerLaunchOption.Remove, "remove", "", false,
		i18n.T("Automatically remove the container when it exits"))
	flags.BoolVarP(&computerLaunchOption.Detach, "detach", "d", false,
		i18n.T("Run container in background and print container ID"))
	flags.StringVarP(&computerLaunchOption.Restart, "restart", "", "no",
		i18n.T("Restart policy to apply when a container exits"))
	flags.StringVarP(&computerLaunchOption.Type, "type", "", AgentJNLP,
		i18n.T("The type of agent, include jnlp"))
	flags.StringVarP(&computerLaunchOption.AgentType, "agent-type", "", "",
		i18n.T("The type of agent, include generic, maven, python. See also https://github.com/jenkinsci/jnlp-agents"))
	flags.BoolVarP(&computerLaunchOption.ShowProgress, "show-progress", "", true,
		i18n.T("Show the progress of downloading agent.jar"))

	if err := computerLaunchCmd.RegisterFlagCompletionFunc("restart", common.ArrayCompletion("no", "always")); err != nil {
		pluginCmd.PrintErrln(err)
	}
	if err := computerLaunchCmd.RegisterFlagCompletionFunc("mode", common.ArrayCompletion("java", "docker")); err != nil {
		pluginCmd.PrintErrln(err)
	}
	if err := computerLaunchCmd.RegisterFlagCompletionFunc("agent-type", common.ArrayCompletion(
		"generic", "maven", "python", "node", "ruby", "docker")); err != nil {
		pluginCmd.PrintErrln(err)
	}

	healthCheckRegister.Register(getCmdPath(computerLaunchCmd), &computerLaunchOption)
}

var computerLaunchCmd = &cobra.Command{
	Use:               "launch",
	Aliases:           []string{"start"},
	Short:             i18n.T("Launch the agent of your Jenkins"),
	Long:              i18n.T("Launch the agent of your Jenkins"),
	ValidArgsFunction: ValidAgentNames,
	Args:              cobra.MinimumNArgs(1),
	Example: `jcli agent launch agent-name
jcli agent launch agent-name --type jnlp`,
	PreRunE: func(_ *cobra.Command, args []string) (err error) {
		computerLaunchOption.ComputerClient, computerLaunchOption.CurrentJenkins =
			GetComputerClient(computerLaunchOption.Option)

		if computerLaunchOption.Type != AgentJNLP {
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
		case AgentJNLP:
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

		if o.Mode == "java" {
			var binary string
			binary, err = util.LookPath("java", centerStartOption.LookPathContext)
			if err == nil {
				env := os.Environ()
				agentArgs := []string{"java", "-jar", computerLaunchOption.Output,
					"-jnlpUrl", fmt.Sprintf("%s/computer/%s/slave-agent.jnlp", o.ComputerClient.URL, name),
					"-secret", secret, "-workDir", "/tmp"}

				if o.CurrentJenkins.ProxyAuth != "" {
					proxyURL, _ := url.Parse(o.CurrentJenkins.Proxy)
					agentArgs = append(agentArgs, "-proxyCredentials", o.CurrentJenkins.ProxyAuth)

					proxyAuth := strings.SplitN(o.CurrentJenkins.ProxyAuth, ":", 2)
					if len(proxyAuth) == 2 {
						env = append(env, fmt.Sprintf("http_proxy=http://%s:%s@%s", url.QueryEscape(proxyAuth[0]), url.QueryEscape(proxyAuth[1]), proxyURL.Host))
					}
				}

				logger.Debug("start a jnlp agent", zap.Any("command", strings.Join(agentArgs, " ")))

				err = util.Exec(binary, agentArgs, env, o.SystemCallExec)
			}
		} else if o.Mode == "docker" {
			var binary string
			binary, err = util.LookPath("docker", centerStartOption.LookPathContext)
			if err == nil {
				env := os.Environ()
				agentArgs := []string{"docker", "run", "--restart", o.Restart}

				if o.Remove {
					agentArgs = append(agentArgs, "--rm")
				}

				if o.Detach {
					agentArgs = append(agentArgs, "--detach")
				}

				var agentImage string
				switch o.AgentType {
				case "generic":
					agentImage = "jenkins/slave:4.0.1-1-alpine"
					agentArgs = append(agentArgs, []string{agentImage, "java", "-jar", "/usr/share/jenkins/agent.jar",
						"-jnlpUrl", fmt.Sprintf("%s/computer/%s/slave-agent.jnlp", o.ComputerClient.URL, name),
						"-secret", secret, "-workDir", "/tmp"}...)
				case "maven", "python", "docker", "node", "ruby":
					agentImage = fmt.Sprintf("jenkins/jnlp-agent-%s", o.AgentType)
					agentArgs = append(agentArgs, []string{agentImage, "-url", o.ComputerClient.URL, secret, name}...)
				}

				if o.CurrentJenkins.ProxyAuth != "" {
					proxyURL, _ := url.Parse(o.CurrentJenkins.Proxy)
					agentArgs = append(agentArgs, "-proxyCredentials", o.CurrentJenkins.ProxyAuth)

					proxyAuth := strings.SplitN(o.CurrentJenkins.ProxyAuth, ":", 2)
					if len(proxyAuth) == 2 {
						env = append(env, fmt.Sprintf("http_proxy=http://%s:%s@%s", url.QueryEscape(proxyAuth[0]), url.QueryEscape(proxyAuth[1]), proxyURL.Host))
					}
				}

				logger.Debug("start a jnlp agent", zap.Any("command", strings.Join(agentArgs, " ")))

				err = util.Exec(binary, agentArgs, env, o.SystemCallExec)
			}
		} else {
			err = fmt.Errorf("not support mode: %s\n", o.Mode)
		}
	}
	return
}

// Check do the health check of casc cmd
func (o *ComputerLaunchOption) Check() (err error) {
	opt := PluginOptions{
		Option: common.Option{RoundTripper: o.RoundTripper},
	}
	_, err = opt.FindPlugin("pipeline-restful-api")
	return
}

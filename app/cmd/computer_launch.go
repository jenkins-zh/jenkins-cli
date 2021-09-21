package cmd

import (
	"errors"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	appCfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/jenkins-zh/jenkins-client/pkg/computer"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
	"strings"
)

// ComputerLaunchOption option for config list command
type ComputerLaunchOption struct {
	common.Option

	Type         string
	ShowProgress bool

	/** share info between inner functions */
	ComputerClient *computer.Client
	CurrentJenkins *appCfg.JenkinsServer
	Output         string

	Mode                 LaunchMode
	Remove               bool
	Restart              string
	Detach               bool
	AgentType            JNLPAgentImage
	AgentImageTag        string
	GeneralAgentImageTag string
	CustomImage          string
}

const (
	// AgentJNLP is the agent type of jnlp
	AgentJNLP = "jnlp"
)

// LaunchMode represents Jenkins agent launch mode
type LaunchMode string

const (
	// LaunchModeJava represents java launch mode
	LaunchModeJava LaunchMode = "java"
	// LaunchModeDocker represents docker launch mode
	LaunchModeDocker LaunchMode = "docker"
)

// Equal determine if they are same
func (l LaunchMode) Equal(mode string) bool {
	return string(l) == mode
}

// String returns the string of the mode
func (l LaunchMode) String() string {
	return string(l)
}

// All returns all launch modes
func (l LaunchMode) All() []string {
	return []string{LaunchModeDocker.String(), LaunchModeJava.String()}
}

// Set give a appropriate value
func (l *LaunchMode) Set(s string) (err error) {
	switch s {
	case LaunchModeDocker.String():
		*l = LaunchModeDocker
	case LaunchModeJava.String():
		*l = LaunchModeJava
	default:
		err = fmt.Errorf("invalid launch mode: %s", s)
	}
	return
}

// Type returns the type of current struct
func (l LaunchMode) Type() string {
	return "LaunchMode"
}

// JNLPAgentImage is the type of Jenkins JNLP agent image
type JNLPAgentImage string

const (
	// GenericAgentImage represents JNLP agent image for generic
	GenericAgentImage JNLPAgentImage = "generic"
	// GolangAgentImage represents JNLP agent image for golang
	GolangAgentImage JNLPAgentImage = "golang"
	// MavenAgentImage represents JNLP agent image for maven
	MavenAgentImage JNLPAgentImage = "maven"
	// PythonAgentImage represents JNLP agent image for python
	PythonAgentImage JNLPAgentImage = "python"
	// NodeAgentImage represents JNLP agent image for node
	NodeAgentImage JNLPAgentImage = "node"
	// RubyAgentImage represents JNLP agent image for ruby
	RubyAgentImage JNLPAgentImage = "ruby"
	// DockerAgentImage represents JNLP agent image for docker
	DockerAgentImage JNLPAgentImage = "docker"
	// TerraformAgentImage represents JNLP agent image for terraform
	TerraformAgentImage JNLPAgentImage = "terraform"
	// CustomAgentImage represents JNLP agent image for custom
	CustomAgentImage JNLPAgentImage = "custom"
)

// All returns all the supported image list
func (i JNLPAgentImage) All() []string {
	return []string{GenericAgentImage.String(), GolangAgentImage.String(), MavenAgentImage.String(),
		PythonAgentImage.String(), NodeAgentImage.String(), RubyAgentImage.String(), DockerAgentImage.String(),
		TerraformAgentImage.String(), CustomAgentImage.String()}
}

// Type returns the type of current struct
func (i JNLPAgentImage) Type() string {
	return "JNLPAgentImage"
}

// String returns the string format of JNLPAgentImage
func (i JNLPAgentImage) String() string {
	return string(i)
}

// Set give a appropriate value
func (i *JNLPAgentImage) Set(s string) (err error) {
	switch s {
	case GenericAgentImage.String():
		*i = GenericAgentImage
	case GolangAgentImage.String():
		*i = GolangAgentImage
	case MavenAgentImage.String():
		*i = MavenAgentImage
	case PythonAgentImage.String():
		*i = PythonAgentImage
	case NodeAgentImage.String():
		*i = NodeAgentImage
	case RubyAgentImage.String():
		*i = RubyAgentImage
	case DockerAgentImage.String():
		*i = DockerAgentImage
	case TerraformAgentImage.String():
		*i = TerraformAgentImage
	case CustomAgentImage.String():
		*i = CustomAgentImage
	default:
		err = fmt.Errorf("invalid JNLP agent image: %s", s)
	}
	return
}

var computerLaunchOption ComputerLaunchOption

func init() {
	computerCmd.AddCommand(computerLaunchCmd)
	flags := computerLaunchCmd.Flags()

	flags.VarP(&computerLaunchOption.Mode, "mode", "m",
		i18n.T(fmt.Sprintf("Mode of launching Jenkins, you can choose: %v", LaunchMode.All(""))))
	flags.BoolVarP(&computerLaunchOption.Remove, "remove", "", false,
		i18n.T("Automatically remove the container when it exits"))
	flags.BoolVarP(&computerLaunchOption.Detach, "detach", "d", false,
		i18n.T("Run container in background and print container ID"))
	flags.StringVarP(&computerLaunchOption.Restart, "restart", "", "no",
		i18n.T("Restart policy to apply when a container exits"))
	flags.StringVarP(&computerLaunchOption.CustomImage, "custom-image", "", "",
		i18n.T(fmt.Sprintf("The custom docker image of Jenkins agent. It works only you set --agent-type=%s", CustomAgentImage.String())))
	flags.StringVarP(&computerLaunchOption.Type, "type", "", AgentJNLP,
		i18n.T("The type of agent, include jnlp"))
	flags.VarP(&computerLaunchOption.AgentType, "agent-type", "",
		i18n.T(fmt.Sprintf("The type of agent, include %v. See also https://github.com/jenkinsci/jnlp-agents", JNLPAgentImage.All(""))))
	flags.StringVarP(&computerLaunchOption.AgentImageTag, "agent-image-tag", "", "latest",
		i18n.T("The Jenkins agent image tag. See also https://github.com/jenkinsci/jnlp-agents"))
	flags.StringVarP(&computerLaunchOption.GeneralAgentImageTag, "general-agent-image-tag", "", "4.0.1-1-alpine",
		i18n.T("The tag of jenkins/slave. See also "))
	flags.BoolVarP(&computerLaunchOption.ShowProgress, "show-progress", "", true,
		i18n.T("Show the progress of downloading agent.jar"))

	if err := computerLaunchCmd.RegisterFlagCompletionFunc("restart", common.ArrayCompletion("no", "always")); err != nil {
		pluginCmd.PrintErrln(err)
	}
	if err := computerLaunchCmd.RegisterFlagCompletionFunc("mode", common.ArrayCompletion(LaunchMode.All("")...)); err != nil {
		pluginCmd.PrintErrln(err)
	}
	if err := computerLaunchCmd.RegisterFlagCompletionFunc("agent-type", common.ArrayCompletion(JNLPAgentImage.All("")...)); err != nil {
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

		if computerLaunchOption.Type != AgentJNLP || LaunchModeDocker == computerLaunchOption.Mode {
			return
		}

		var f *os.File
		tmpPath := "/tmp"
		if runtime.GOOS == "windows" {
			userHome, homeErr := homedir.Dir()
			if homeErr == nil {
				tmpPath, _ = ioutil.TempDir(userHome, "/tmp")
			}
		}
		if f, err = ioutil.TempFile(tmpPath, "agent.jar"); err == nil {
			computerLaunchOption.Output = f.Name()
			agentURL := fmt.Sprintf("%s/jnlpJars/agent.jar", computerLaunchOption.ComputerClient.URL)
			logger.Debug("start to download agent.jar", zap.String("url", agentURL))
			logger.Debug("proxy setting", zap.String("sever", computerLaunchOption.CurrentJenkins.Proxy),
				zap.String("auth", computerLaunchOption.CurrentJenkins.ProxyAuth))

			downloader := httpdownloader.HTTPDownloader{
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
	if secret, err = o.ComputerClient.GetSecret(name); err != nil {
		return
	}
	logger.Info("get agent secret", zap.String("value", secret))

	switch o.Mode {
	case LaunchModeJava, "":
		var binary string
		binary, err = util.LookPath("java", centerStartOption.LookPathContext)
		if err == nil {
			env := os.Environ()
			agentArgs := []string{"java", "-jar", computerLaunchOption.Output,
				"-jnlpUrl", fmt.Sprintf("%s/computer/%s/slave-agent.jnlp", o.ComputerClient.URL, name),
				"-secret", secret, "-workDir", computer.GetDefaultAgentWorkDir()}

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
	case LaunchModeDocker:
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
			case GolangAgentImage, MavenAgentImage, PythonAgentImage,
				DockerAgentImage, NodeAgentImage, RubyAgentImage, TerraformAgentImage:
				agentImage = fmt.Sprintf("jenkins/jnlp-agent-%s:%s", o.AgentType, o.AgentImageTag)
			case GenericAgentImage:
				agentImage = fmt.Sprintf("jenkins/inbound-agent:%s", o.AgentImageTag)
			case CustomAgentImage:
				if o.CustomImage == "" {
					err = errors.New("--custom-image cannot be empty if you choose custom agent type")
					return
				}
				agentImage = o.CustomImage
			}
			agentArgs = append(agentArgs, []string{agentImage, "-url", o.ComputerClient.URL, secret, name}...)

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
	default:
		err = fmt.Errorf("not support mode: %s", o.Mode)
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

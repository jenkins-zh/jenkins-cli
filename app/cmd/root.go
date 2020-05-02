package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	. "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/health"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger

// RootOptions is a global option for whole cli
type RootOptions struct {
	ConfigFile string
	ConfigLoad bool
	Jenkins    string
	Debug      bool

	URL                string
	Username           string
	Token              string
	InsecureSkipVerify bool
	Proxy              string
	ProxyAuth          string
	ProxyDisable       bool

	Doctor    bool
	StartTime time.Time
	EndTime   time.Time

	CommonOption *common.CommonOption

	LoggerLevel string
}

var healthCheckRegister = &health.CheckRegister{
	Member: make(map[string]health.CommandHealth, 0),
}

var rootCmd = &cobra.Command{
	Use:   "jcli",
	Short: i18n.T("Jenkins CLI written by golang which could help you with your multiple Jenkins"),
	Long: `Jenkins CLI written by golang which could help you with your multiple Jenkins,

We'd love to hear your feedback at https://github.com/jenkins-zh/jenkins-cli/issues`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		rootOptions.StartTime = time.Now()
		if logger, err = util.InitLogger(rootOptions.LoggerLevel); err == nil {
			(&configOptions).Logger = logger
			client.SetLogger(logger)
		} else {
			return
		}

		rootOptions.ConfigLoad = !("false" == os.Getenv("JCLI_CONFIG_LOAD"))
		if rootOptions.ConfigLoad && needReadConfig(cmd) {
			if rootOptions.ConfigFile == "" {
				rootOptions.ConfigFile = os.Getenv("JCLI_CONFIG")
			}

			logger.Debug("read config file", zap.String("path", rootOptions.ConfigFile))
			if rootOptions.ConfigFile == "" {
				err = loadDefaultConfig()
			} else {
				err = loadConfig(rootOptions.ConfigFile)
			}
		} else {
			logger.Debug("ignore loading config", zap.String("cmd", cmd.Name()))
		}

		if err == nil {
			config = getConfig()
			if config != nil {
				// set Header Accept-Language
				client.SetLanguage(config.Language)
			}

			err = rootOptions.RunDiagnose(cmd)
		}
		return
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		cmdPath := getCmdPath(cmd)

		// calculate the time
		rootOptions.EndTime = time.Now()

		logger.Debug("done with command", zap.String("command", cmdPath),
			zap.Float64("duration", rootOptions.EndTime.Sub(rootOptions.StartTime).Seconds()))
	},
	BashCompletionFunction: jcliBashCompletionFunc,
}

func needReadConfig(cmd *cobra.Command) bool {
	ignoreConfigLoad := []string{
		"config.generate",
		//"center.start", // relay on the config when find a mirror
		"cwp",
		"version",
		"completion",
		"doc",
	}
	configPath := getCmdPath(cmd)

	for _, item := range ignoreConfigLoad {
		if item == configPath {
			return false
		}
	}

	// allow sub-commands give their decisions
	if cmd.Annotations != nil {
		if disable, ok := cmd.Annotations[ANNOTATION_CONFIG_LOAD]; ok {
			return disable != "disable"
		}
	}
	return true
}

// RunDiagnose run the diagnose for a specific command
func (o *RootOptions) RunDiagnose(cmd *cobra.Command) (err error) {
	if !o.Doctor {
		return
	}
	path := getCmdPath(cmd)
	logger.Debug("start to run diagnose", zap.String("path", path))

	for k, v := range healthCheckRegister.Member {
		if ok, _ := regexp.MatchString(k, path); ok {
			err = v.Check()
			break
		}

	}
	return
}

// Execute will execute the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootOptions RootOptions

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootOptions.ConfigFile, "configFile", "", "",
		i18n.T("An alternative config file"))
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.ConfigLoad, "config-load", "", true,
		i18n.T("If load a default config file"))
	rootCmd.PersistentFlags().StringVarP(&rootOptions.Jenkins, "jenkins", "j", "",
		i18n.T("Select a Jenkins server for this time"))
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Debug, "debug", "", false, "Print the output into debug.html")
	rootCmd.PersistentFlags().StringVarP(&rootOptions.LoggerLevel, "logger-level", "", "warn",
		"Logger level which could be: debug, info, warn, error")
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Doctor, "doctor", "", false,
		i18n.T("Run the diagnose for current command"))

	rootCmd.PersistentFlags().StringVarP(&rootOptions.URL, "url", "", "",
		i18n.T("The URL of Jenkins"))
	rootCmd.PersistentFlags().StringVarP(&rootOptions.Username, "username", "", "",
		i18n.T("The username of Jenkins"))
	rootCmd.PersistentFlags().StringVarP(&rootOptions.Token, "token", "", "",
		i18n.T("The token of Jenkins"))
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.InsecureSkipVerify, "insecureSkipVerify", "", true,
		i18n.T("If skip insecure skip verify"))
	rootCmd.PersistentFlags().StringVarP(&rootOptions.Proxy, "proxy", "", "",
		i18n.T("The proxy of connection to Jenkins"))
	rootCmd.PersistentFlags().StringVarP(&rootOptions.ProxyAuth, "proxy-auth", "", "",
		i18n.T("The auth of proxy of connection to Jenkins"))
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.ProxyDisable, "proxy-disable", "", false,
		i18n.T("Disable proxy setting"))

	rootCmd.SetOut(os.Stdout)

	loadPlugins(rootCmd)

	// add sub-commands
	NewShutdownCmd(&rootOptions)
}

// GetRootOptions returns the root options
func GetRootOptions() *RootOptions {
	return &rootOptions
}

// GetRootCommand returns the root cmd
func GetRootCommand() *cobra.Command {
	return rootCmd
}

func getCurrentJenkinsFromOptions() (jenkinsServer *JenkinsServer) {
	jenkinsOpt := rootOptions.Jenkins

	if jenkinsOpt == "" {
		jenkinsServer = getCurrentJenkins()
	} else {
		jenkinsServer = findJenkinsByName(jenkinsOpt)
	}

	// take URL from options if it's not empty
	if jenkinsServer == nil && rootOptions.URL != "" {
		jenkinsServer = &JenkinsServer{}
	}

	if jenkinsServer != nil {
		if rootOptions.URL != "" {
			jenkinsServer.URL = rootOptions.URL
		}

		if rootOptions.Username != "" {
			jenkinsServer.UserName = rootOptions.Username
		}

		if rootOptions.Token != "" {
			jenkinsServer.Token = rootOptions.Token
		}

		if rootOptions.Proxy != "" {
			jenkinsServer.Proxy = rootOptions.Proxy
		}

		if rootOptions.ProxyAuth != "" {
			jenkinsServer.ProxyAuth = rootOptions.ProxyAuth
		}

		if rootOptions.ProxyDisable {
			jenkinsServer.Proxy = ""
			jenkinsServer.ProxyAuth = ""
		}
	}
	return
}

// Deprecated, please use getCurrentJenkinsFromOptions instead of it
func getCurrentJenkinsFromOptionsOrDie() (jenkinsServer *JenkinsServer) {
	if jenkinsServer = getCurrentJenkinsFromOptions(); jenkinsServer == nil {
		log.Fatal("Cannot found Jenkins by", rootOptions.Jenkins)
	}
	return
}

func getCmdPath(cmd *cobra.Command) string {
	current := cmd.Use
	if cmd.HasParent() {
		parentName := getCmdPath(cmd.Parent())
		if parentName == "" {
			return current
		}

		return fmt.Sprintf("%s.%s", parentName, current)
	}
	// don't need the name of root cmd
	return ""
}

func executePreCmd(cmd *cobra.Command, _ []string, writer io.Writer) (err error) {
	config := getConfig()
	if config == nil {
		err = fmt.Errorf("cannot find config file")
		return
	}

	path := getCmdPath(cmd)
	for _, hook := range config.PreHooks {
		if path != hook.Path {
			continue
		}

		if err = execute(hook.Command, writer); err != nil {
			return
		}
	}
	return
}

func executePostCmd(cmd *cobra.Command, _ []string, writer io.Writer) (err error) {
	config := getConfig()
	if config == nil {
		err = fmt.Errorf("cannot find config file")
		return
	}

	path := getCmdPath(cmd)
	for _, hook := range config.PostHooks {
		if path != hook.Path {
			continue
		}

		if err = execute(hook.Command, writer); err != nil {
			return
		}
	}
	return
}

func execute(command string, writer io.Writer) (err error) {
	array := strings.Split(command, " ")
	cmd := exec.Command(array[0], array[1:]...)
	if err = cmd.Start(); err == nil {
		if err = cmd.Wait(); err == nil {
			var data []byte
			if data, err = cmd.Output(); err == nil {
				_, _ = writer.Write(data)
			}
		}
	}
	return
}

// Deprecated, please replace this with getCurrentJenkinsAndClient
func getCurrentJenkinsAndClientOrDie(jclient *client.JenkinsCore) (jenkins *JenkinsServer) {
	jenkins = getCurrentJenkinsFromOptionsOrDie()
	jclient.URL = jenkins.URL
	jclient.UserName = jenkins.UserName
	jclient.Token = jenkins.Token
	jclient.Proxy = jenkins.Proxy
	jclient.ProxyAuth = jenkins.ProxyAuth
	return
}

func getCurrentJenkinsAndClient(jClient *client.JenkinsCore) (jenkins *JenkinsServer) {
	if jenkins = getCurrentJenkinsFromOptions(); jenkins != nil {
		jClient.URL = jenkins.URL
		jClient.UserName = jenkins.UserName
		jClient.Token = jenkins.Token
		jClient.Proxy = jenkins.Proxy
		jClient.ProxyAuth = jenkins.ProxyAuth
		jClient.InsecureSkipVerify = jenkins.InsecureSkipVerify
	}
	return
}

const (
	jcliBashCompletionFunc = `__plugin_name_parse_get()
{
    local jcli_output out
    if jcli_output=$(jcli plugin list --filter HasUpdate=true --no-headers --filter ShortName="$1" 2>/dev/null); then
        out=($(echo "${jcli_output}" | awk '{print $1}'))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__jcli_get_plugin_name()
{
    __plugin_name_parse_get
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__config_name_parse_get()
{
    local jcli_output out
    if jcli_output=$(jcli config list --no-headers 2>/dev/null); then
        out=($(echo "${jcli_output}" | awk '{print $1}'))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__jcli_get_config_name()
{
    __config_name_parse_get
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__job_name_parse_get()
{
    local jcli_output out
    if jcli_output=$(jcli job search --columns URL --no-headers "$cur" 2>/dev/null); then
        out=($(echo "${jcli_output}"))
        COMPREPLY=( ${out} )
    fi
}

__jcli_get_job_name()
{
    __job_name_parse_get
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__computer_name_parse_get()
{
    local jcli_output out
    if jcli_output=$(jcli computer list --no-headers --columns DisplayName 2>/dev/null); then
        out=($(echo "${jcli_output}"))
        COMPREPLY=( ${out} )
    fi
}

__jcli_get_computer_name()
{
    __computer_name_parse_get
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__jcli_custom_func() {
    case ${last_command} in
        jcli_plugin_upgrade | jcli_plugin_uninstall)
            __jcli_get_plugin_name
            return
            ;;
        jcli_open | jcli_config_select | jcli_config_remove | jcli_shell)
            __jcli_get_config_name
            return
            ;;
        jcli_computer_delete | jcli_computer_delete | jcli_computer_launch)
            __jcli_get_computer_name
            return
            ;;
        jcli_job_build | jcli_job_stop | jcli_job_log | jcli_job_delete | jcli_job_history | jcli_job_artifact | jcli_job_input)
            __jcli_get_job_name
            return
            ;;
        *)
            ;;
    esac
}
`
)

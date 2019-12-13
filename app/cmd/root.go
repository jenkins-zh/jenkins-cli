package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/health"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger

// RootOptions is a global option for whole cli
type RootOptions struct {
	ConfigFile string
	Jenkins    string
	Version    bool
	Debug      bool

	Doctor bool

	LoggerLevel string
}

var healthCheckRegister = &health.CheckRegister{
	Member: make(map[string]health.CommandHealth, 0),
}

var rootCmd = &cobra.Command{
	Use:   "jcli",
	Short: i18n.T("jcli is a tool which could help you with your multiple Jenkins"),
	Long: `jcli is Jenkins CLI which could help with your multiple Jenkins,
Manage your Jenkins and your pipelines
More information could found at https://jenkins-zh.cn`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if logger, err = util.InitLogger(rootOptions.LoggerLevel); err == nil {
			client.SetLogger(logger)
		} else {
			return
		}

		if rootOptions.ConfigFile == "" {
			rootOptions.ConfigFile = os.Getenv("JCLI_CONFIG")
		}

		logger.Debug("read config file", zap.String("path", rootOptions.ConfigFile))
		if rootOptions.Version && cmd.Flags().NFlag() == 1 {
			return
		}

		if rootOptions.ConfigFile == "" {
			if err := loadDefaultConfig(); err != nil {
				configLoadErrorHandle(err)
			}
		} else {
			if err := loadConfig(rootOptions.ConfigFile); err != nil {
				configLoadErrorHandle(err)
			}
		}

		// set Header Accept-Language
		config = getConfig()
		if config != nil {
			client.SetLanguage(config.Language)
		}

		err = rootOptions.RunDiagnose(cmd)
		return
	},
	BashCompletionFunction: jcliBashCompletionFunc,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(i18n.T("Jenkins CLI (jcli) manage your Jenkins"))
		if rootOptions.Version {
			cmd.Printf("Version: %s\n", app.GetVersion())
			cmd.Printf("Commit: %s\n", app.GetCommit())
		}
		if rootOptions.Jenkins != "" {
			current := getCurrentJenkinsFromOptionsOrDie()
			if current != nil {
				cmd.Println("Current Jenkins is:", current.Name)
			} else {
				cmd.Println("Cannot found the configuration")
			}
		}
	},
}

// RunDiagnose run the diagnose for a specific command
func (o *RootOptions) RunDiagnose(cmd *cobra.Command) (err error) {
	if !o.Doctor {
		return
	}
	path := getCmdPath(cmd)

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
	rootCmd.PersistentFlags().StringVarP(&rootOptions.Jenkins, "jenkins", "j", "",
		i18n.T("Select a Jenkins server for this time"))
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Debug, "debug", "", false, "Print the output into debug.html")
	rootCmd.PersistentFlags().StringVarP(&rootOptions.LoggerLevel, "logger-level", "", "warn",
		"Logger level which could be: debug, info, warn, error")
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Doctor, "doctor", "", false,
		i18n.T("Run the diagnose for current command"))
	rootCmd.Flags().BoolVarP(&rootOptions.Version, "version", "v", false,
		i18n.T("Print the version of Jenkins CLI"))
	rootCmd.SetOut(os.Stdout)
}

func configLoadErrorHandle(err error) {
	if os.IsNotExist(err) {
		log.Printf("No config file found.")
		return
	}

	log.Fatalf("Config file is invalid: %v", err)
}

func getCurrentJenkinsFromOptions() (jenkinsServer *JenkinsServer) {
	jenkinsOpt := rootOptions.Jenkins

	if jenkinsOpt == "" {
		jenkinsServer = getCurrentJenkins()
	} else {
		jenkinsServer = findJenkinsByName(jenkinsOpt)
	}
	return
}

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
		err = fmt.Errorf("Cannot find config file")
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
	cmd.Stdout = writer
	err = cmd.Run()
	return
}

const (
	jcliBashCompletionFunc = `__plugin_name_parse_get()
{
    local jcli_output out
    if jcli_output=$(jcli plugin list --filter hasUpdate --no-headers --filter name="$1" 2>/dev/null); then
        out=($(echo "${jcli_output}" | awk '{print $2}'))
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
        out=($(echo "${jcli_output}" | awk '{print $2}'))
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
    if jcli_output=$(jcli job search -o path "$cur" 2>/dev/null); then
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

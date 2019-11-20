package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

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

	LoggerLevel string
}

var rootCmd = &cobra.Command{
	Use:   "jcli",
	Short: "jcli is a tool which could help you with your multiple Jenkins",
	Long: `jcli is Jenkins CLI which could help with your multiple Jenkins,
				  Manage your Jenkins and your pipelines
				  More information could found at https://jenkins-zh.cn`,
	BashCompletionFunction: jcliBashCompletionFunc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if logger, err = util.InitLogger(rootOptions.LoggerLevel); err != nil {
			cmd.PrintErrln(err)
		} else {
			client.SetLogger(logger)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Jenkins CLI (jcli) manage your Jenkins")
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

// Execute will exectue the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootOptions RootOptions

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&rootOptions.ConfigFile, "configFile", "", "", "An alternative config file")
	rootCmd.PersistentFlags().StringVarP(&rootOptions.Jenkins, "jenkins", "j", "", "Select a Jenkins server for this time")
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Version, "version", "v", false, "Print the version of Jenkins CLI")
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Debug, "debug", "", false, "Print the output into debug.html")
	rootCmd.PersistentFlags().StringVarP(&rootOptions.LoggerLevel, "logger-level", "", "warn",
		"Logger level which could be: debug, info, warn, error")
	rootCmd.SetOut(os.Stdout)
}

func initConfig() {
	if rootOptions.Version && rootCmd.Flags().NFlag() == 1 {
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
		log.Fatal("Cannot found Jenkins by", rootOptions.Jenkins) // TODO not accurate
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

__jcli_custom_func() {
    case ${last_command} in
        jcli_plugin_upgrade)
            __jcli_get_plugin_name
            return
            ;;
        *)
            ;;
    esac
}
`
)

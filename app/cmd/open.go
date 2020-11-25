package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	appCfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// OpenOption is the open cmd option
type OpenOption struct {
	common.Option
	common.InteractiveOption

	Browser string

	Config bool
}

var openOption OpenOption

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().BoolVarP(&openOption.Config, "config", "c", false,
		i18n.T("Open the configuration page of Jenkins"))
	openCmd.Flags().StringVarP(&openOption.Browser, "browser", "b", "",
		i18n.T("Open Jenkins with a specific browser"))
	openOption.SetFlag(openCmd)
	openOption.Stdio = common.GetSystemStdio()

	err := openCmd.RegisterFlagCompletionFunc("browser", func(cmd *cobra.Command, args []string, toComplete string) (strings []string, directive cobra.ShellCompDirective) {
		return []string{"Google-Chrome", "Safari", "Microsoft-Edge", "Firefox"}, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		rootCmd.Println(err)
	}
}

var openCmd = &cobra.Command{
	Use:               "open",
	Short:             i18n.T("Open your Jenkins with a browser"),
	Long:              i18n.T(`Open your Jenkins with a browser`),
	ValidArgsFunction: ValidJenkinsAndDataNames,
	Example: `jcli open -n [config name]
Open Jenkins with a specific browser is useful in some use cases. For example, one browser has a proxy setting.
There are two ways to achieve this:
jcli open --browser "Google-Chrome"
JCLI_BROWSER="Google Chrome" jcli open`,
	PreRun: func(_ *cobra.Command, _ []string) {
		if openOption.Browser == "" {
			openOption.Browser = os.Getenv("JCLI_BROWSER")
		}
	},
	RunE: openOption.run,
}

func (o *OpenOption) run(_ *cobra.Command, args []string) (err error) {
	var jenkins *appCfg.JenkinsServer

	var configName string
	if len(args) > 0 {
		configName = args[0]
	}

	if configName == "" && o.Interactive {
		jenkinsNames := getJenkinsNames()
		configName, err = o.Select(jenkinsNames,
			i18n.T("Choose a Jenkins which you want to open:"), "")
	}

	jenkinsName, external := o.parseName(configName)
	logger.Info("open jenkins",
		zap.String("jenkins name", jenkinsName),
		zap.String("external", external))

	if err == nil {
		if jenkinsName != "" {
			jenkins = findJenkinsByName(jenkinsName)
		} else {
			jenkins = getCurrentJenkins()
		}

		if jenkins != nil && jenkins.URL != "" {
			url := jenkins.URL
			if o.Config {
				url = fmt.Sprintf("%s/configure", url)
			} else if external != "" {
				url = jenkins.Data[external]
			}
			err = o.smartOpen(url)
		} else {
			err = fmt.Errorf("no URL found with Jenkins %s", jenkinsName)
		}
	}
	return
}

// smartOpen can open with or without specific protocol
func (o *OpenOption) smartOpen(url string) (err error) {
	if strings.HasPrefix(url, "ssh://") {
		var ssh string
		if ssh, err = exec.LookPath("ssh"); err == nil {
			err = syscall.Exec(ssh, []string{"ssh", strings.TrimLeft(url, "ssh://")}, os.Environ())
		}
	} else {
		browser := o.Browser
		err = util.Open(url, browser, o.ExecContext)
	}
	return
}

// parseName the string expect likes name or name.external
func (o *OpenOption) parseName(configName string) (jenkins, external string) {
	array := strings.SplitN(configName, ".", 2)
	if len(array) > 0 {
		jenkins = array[0]
	}
	if len(array) > 1 {
		external = array[1]
	}
	return
}

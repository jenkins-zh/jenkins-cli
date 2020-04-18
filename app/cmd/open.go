package cmd

import (
	"fmt"
	. "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"os"
)

// OpenOption is the open cmd option
type OpenOption struct {
	CommonOption
	InteractiveOption

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
	openOption.Stdio = GetSystemStdio()

	err := openCmd.RegisterFlagCompletionFunc("browser", func(cmd *cobra.Command, args []string, toComplete string) (strings []string, directive cobra.ShellCompDirective) {
		return []string{"Google-Chrome", "Safari", "Microsoft-Edge", "Firefox"}, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		rootCmd.Println(err)
	}
}

var openCmd = &cobra.Command{
	Use:     "open",
	Short:   i18n.T("Open your Jenkins with a browser"),
	Long:    i18n.T(`Open your Jenkins with a browser`),
	Example: `jcli open -n [config name]`,
	PreRun: func(_ *cobra.Command, _ []string) {
		if openOption.Browser == "" {
			openOption.Browser = os.Getenv("BROWSER")
		}
	},
	RunE: func(_ *cobra.Command, args []string) (err error) {
		var jenkins *JenkinsServer

		var configName string
		if len(args) > 0 {
			configName = args[0]
		}

		if configName == "" && openOption.Interactive {
			jenkinsNames := getJenkinsNames()
			configName, err = openOption.Select(jenkinsNames,
				i18n.T("Choose a Jenkins which you want to open:"), "")
		}

		if err == nil {
			if configName != "" {
				jenkins = findJenkinsByName(configName)
			} else {
				jenkins = getCurrentJenkins()
			}

			if jenkins != nil && jenkins.URL != "" {
				url := jenkins.URL
				if openOption.Config {
					url = fmt.Sprintf("%s/configure", url)
				}
				browser := openOption.Browser
				err = util.Open(url, browser, openOption.ExecContext)
			} else {
				err = fmt.Errorf("no URL found with Jenkins %s", configName)
			}
		}
		return
	},
}

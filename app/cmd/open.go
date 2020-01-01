package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// OpenOption is the open cmd option
type OpenOption struct {
	CommonOption
	InteractiveOption

	Config bool
}

var openOption OpenOption

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().BoolVarP(&openOption.Config, "config", "c", false,
		i18n.T("Open the configuration page of Jenkins"))
	openOption.SetFlag(openCmd)
}

var openCmd = &cobra.Command{
	Use:     "open",
	Short:   i18n.T("Open your Jenkins with a browser"),
	Long:    i18n.T(`Open your Jenkins with a browser`),
	Example: `jcli open -n [config name]`,
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
				err = util.Open(url, openOption.ExecContext)
			} else {
				err = fmt.Errorf("no URL found with Jenkins %s", configName)
			}
		}
		return
	},
}

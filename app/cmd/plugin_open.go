package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"os"

	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// PluginOpenOption is the option of plugin open cmd
type PluginOpenOption struct {
	ExecContext util.ExecContext

	Browser string
}

var pluginOpenOption PluginOpenOption

func init() {
	pluginCmd.AddCommand(pluginOpenCmd)
	pluginOpenCmd.Flags().StringVarP(&pluginOpenOption.Browser, "browser", "b", "",
		i18n.T("Open Jenkins with a specific browser"))
}

var pluginOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open update center server in browser",
	Long:  `Open update center server in browser`,
	PreRun: func(_ *cobra.Command, _ []string) {
		if pluginOpenOption.Browser == "" {
			pluginOpenOption.Browser = os.Getenv("BROWSER")
		}
	},
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		jenkins := getCurrentJenkinsFromOptions()
		if jenkins == nil {
			err = fmt.Errorf("cannot found Jenkins by %s", rootOptions.Jenkins)
			return
		}

		if jenkins.URL != "" {
			browser := pluginOpenOption.Browser
			err = util.Open(fmt.Sprintf("%s/pluginManager", jenkins.URL), browser, pluginOpenOption.ExecContext)
		} else {
			err = fmt.Errorf("no URL fond from %s", jenkins.Name)
		}
		return
	},
}

package cmd

import (
	"fmt"
	"os"

	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// PluginOpenOption is the option of plugin open cmd
type PluginOpenOption struct {
	ExecContext util.ExecContext
}

var pluginOpenOption PluginOpenOption

func init() {
	pluginCmd.AddCommand(pluginOpenCmd)
}

var pluginOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open update center server in browser",
	Long:  `Open update center server in browser`,
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()

		if jenkins.URL != "" {
			browser := os.Getenv("BROWSER")
			err = util.Open(fmt.Sprintf("%s/pluginManager", jenkins.URL), browser, pluginOpenOption.ExecContext)
		} else {
			err = fmt.Errorf("no URL fond from %s", jenkins.Name)
		}
		return
	},
}

package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"os"

	"github.com/jenkins-zh/jenkins-cli/util"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// PluginRunOptions for the plugin run command
type PluginRunOptions struct {
	common.Option

	DebugOutput bool
}

var pluginRunOptions PluginRunOptions

func init() {
	pluginCmd.AddCommand(pluginRunCmd)
	pluginRunCmd.Flags().BoolVar(&pluginRunOptions.DebugOutput, "debug-output", false,
		i18n.T("If you want the maven output the debug info"))
}

var pluginRunCmd = &cobra.Command{
	Use:   "run",
	Short: i18n.T("Run the Jenkins plugin project"),
	Long: i18n.T(`Run the Jenkins plugin project
The default behaviour is "mvn hpi:run"`),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		binary, err := util.LookPath("mvn", pluginRunOptions.LookPathContext)
		if err == nil {
			env := os.Environ()

			mvnArgs := []string{"mvn", "hpi:run", "-Dhpi.prefix=/", "-Djetty.port=8080"}
			if pluginRunOptions.DebugOutput {
				mvnArgs = append(mvnArgs, "-X")
			}
			err = util.Exec(binary, mvnArgs, env, pluginRunOptions.SystemCallExec)
		}
		return
	},
	Annotations: map[string]string{
		common.Since: "v0.0.31",
	},
}

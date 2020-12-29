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

	CleanHome   bool
	DebugOutput bool
}

var pluginRunOptions PluginRunOptions

func init() {
	pluginCmd.AddCommand(pluginRunCmd)

	flags := pluginRunCmd.Flags()
	flags.BoolVar(&pluginRunOptions.DebugOutput, "debug-output", false,
		i18n.T("If you want the maven output the debug info"))
	flags.BoolVarP(&pluginRunOptions.CleanHome, "clean-home", "", false,
		i18n.T("If you want to clean the JENKINS_HOME before start it"))
}

var pluginRunCmd = &cobra.Command{
	Use:   "run",
	Short: i18n.T("Run the Jenkins plugin project"),
	Long: i18n.T(`Run the Jenkins plugin project
The default behaviour is "mvn hpi:run"`),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if pluginRunOptions.CleanHome {
			err = os.RemoveAll("work")
		}
		return
	},
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		binary, err := util.LookPath("mvn", pluginRunOptions.LookPathContext)
		if err == nil {
			env := os.Environ()

			mvnArgs := []string{"mvn"}
			if pluginRunOptions.DebugOutput {
				mvnArgs = append(mvnArgs, "-X")
			}
			if pluginRunOptions.CleanHome {
				mvnArgs = append(mvnArgs, "clean")
			}
			mvnArgs = append(mvnArgs, []string{"hpi:run", "-Dhpi.prefix=/", "-Djetty.port=8080"}...)
			err = util.Exec(binary, mvnArgs, env, pluginRunOptions.SystemCallExec)
		}
		return
	},
	Annotations: map[string]string{
		common.Since: "v0.0.31",
	},
}

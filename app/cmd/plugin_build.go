package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"os"

	"github.com/jenkins-zh/jenkins-cli/util"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// PluginBuildOptions for the plugin build command
type PluginBuildOptions struct {
	common.CommonOption

	DebugOutput bool
}

var pluginBuildOptions PluginBuildOptions

func init() {
	pluginCmd.AddCommand(pluginBuildCmd)
	pluginBuildCmd.Flags().BoolVar(&pluginBuildOptions.DebugOutput, "debug-output", false,
		i18n.T("If you want the maven output the debug info"))
}

var pluginBuildCmd = &cobra.Command{
	Use:   "build",
	Short: i18n.T("Build the Jenkins plugin project"),
	Long: i18n.T(`Build the Jenkins plugin project
The default behaviour is "mvn clean package -DskipTests -Dmaven.test.skip"`),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		binary, err := util.LookPath("mvn", pluginBuildOptions.LookPathContext)
		if err == nil {
			env := os.Environ()

			mvnArgs := []string{"mvn", "clean", "package", "-DskipTests", "-Dmaven.test.skip"}
			if pluginBuildOptions.DebugOutput {
				mvnArgs = append(mvnArgs, "-X")
			}
			err = util.Exec(binary, mvnArgs, env, pluginBuildOptions.SystemCallExec)
		}
		return
	},
	Annotations: map[string]string{
		common.Since: "v0.0.27",
	},
}

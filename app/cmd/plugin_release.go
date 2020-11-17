package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"os"

	"github.com/jenkins-zh/jenkins-cli/util"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// PluginReleaseOptions for the plugin create command
type PluginReleaseOptions struct {
	common.Option

	Batch       bool
	Prepare     bool
	Perform     bool
	SkipTests   bool
	DebugOutput bool
}

var pluginReleaseOptions PluginReleaseOptions

func init() {
	pluginCmd.AddCommand(pluginReleaseCmd)
	pluginReleaseCmd.Flags().BoolVar(&pluginReleaseOptions.DebugOutput, "debug-output", false,
		i18n.T("If you want the maven output the debug info"))
	pluginReleaseCmd.Flags().BoolVar(&pluginReleaseOptions.Prepare, "prepare", true,
		i18n.T("Add mvn command release:prepare"))
	pluginReleaseCmd.Flags().BoolVar(&pluginReleaseOptions.Perform, "perform", true,
		i18n.T("Add mvn command release:perform"))
	pluginReleaseCmd.Flags().BoolVar(&pluginReleaseOptions.SkipTests, "skip-tests", true,
		i18n.T("Skip running tests"))
	pluginReleaseCmd.Flags().BoolVar(&pluginReleaseOptions.Batch, "batch", true,
		i18n.T("Run in non-interactive (batch)"))
}

var pluginReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: i18n.T("Release current plugin project"),
	Long:  i18n.T("Release current plugin project"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		binary, err := util.LookPath("mvn", pluginReleaseOptions.LookPathContext)
		if err == nil {
			env := os.Environ()

			mvnArgs := []string{"mvn"}

			if pluginReleaseOptions.SkipTests {
				mvnArgs = append(mvnArgs, `-Darguments="-DskipTests"`)
			}

			if pluginReleaseOptions.Prepare {
				mvnArgs = append(mvnArgs, "release:prepare")
			}

			if pluginReleaseOptions.Perform {
				mvnArgs = append(mvnArgs, "release:perform")
			}

			if pluginReleaseOptions.Batch {
				mvnArgs = append(mvnArgs, "-B")
			}

			if pluginReleaseOptions.DebugOutput {
				mvnArgs = append(mvnArgs, "-X")
			}
			err = util.Exec(binary, mvnArgs, env, pluginReleaseOptions.SystemCallExec)
		}
		return
	},
	Annotations: map[string]string{
		common.Since: common.VersionSince0024,
	},
}

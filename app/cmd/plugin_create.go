package cmd

import (
	"os"

	"github.com/jenkins-zh/jenkins-cli/util"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// PluginCreateOptions for the plugin create command
type PluginCreateOptions struct {
	CommonOption

	DebugOutput bool
}

var pluginCreateOptions PluginCreateOptions

func init() {
	pluginCmd.AddCommand(pluginCreateCmd)
	pluginCreateCmd.Flags().BoolVar(&pluginCreateOptions.DebugOutput, "debug-output", false,
		i18n.T("If you want the maven output the debug info"))
}

var pluginCreateCmd = &cobra.Command{
	Use:   "create",
	Short: i18n.T("Create a plugin project from the archetypes"),
	Long: i18n.T(`Create a plugin project from the archetypes
Plugin tutorial is here https://jenkins.io/doc/developer/tutorial/`),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		binary, err := util.LookPath("mvn", pluginCreateOptions.LookPathContext)
		if err == nil {
			env := os.Environ()

			mvnArgs := []string{"mvn", "archetype:generate", "-U", `-Dfilter=io.jenkins.archetypes:`}
			if pluginCreateOptions.DebugOutput {
				mvnArgs = append(mvnArgs, "-X")
			}
			err = util.Exec(binary, mvnArgs, env, pluginCreateOptions.SystemCallExec)
		}
		return
	},
	Annotations: map[string]string{
		since: "v0.0.23",
	},
}

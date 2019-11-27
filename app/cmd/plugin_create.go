package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// PluginCreateOptions for the plugin create command
type PluginCreateOptions struct {
	Debug bool
}

var pluginCreateOptions PluginCreateOptions

func init() {
	pluginCmd.AddCommand(pluginCreateCmd)
	pluginCreateCmd.Flags().BoolVar(&pluginCreateOptions.Debug, "debug-output", false,
		i18n.T("If you want the maven output the debug info"))
}

var pluginCreateCmd = &cobra.Command{
	Use:   "create",
	Short: i18n.T("Create a plugin project from the archetypes"),
	Long: i18n.T(`Create a plugin project from the archetypes
Plugin tutorial is here https://jenkins.io/doc/developer/tutorial/`),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		binary, err := exec.LookPath("mvn")
		if err == nil {
			env := os.Environ()

			mvnArgs := []string{"mvn", "archetype:generate", "-U", `-Dfilter=io.jenkins.archetypes:`}
			if pluginCreateOptions.Debug {
				mvnArgs = append(mvnArgs, "-X")
			}
			err = syscall.Exec(binary, mvnArgs, env)
		}
		return
	},
}

package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

// CompletionOptions is the option of completion command
type CompletionOptions struct {
	Type string
}

// ShellTypes contains all types of shell
var ShellTypes = []string{
	"zsh", "bash", "powerShell",
}

var completionOptions CompletionOptions

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringVarP(&completionOptions.Type, "type", "", "bash",
		i18n.T(fmt.Sprintf("Generate different types of shell which are %v", ShellTypes)))
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: i18n.T("Generate bash completion scripts"),
	Long:  i18n.T("Generate bash completion scripts"),
	Example: `# Installing bash completion on macOS using homebrew
	## If running Bash 3.2 included with macOS
	brew install bash-completion
	## or, if running Bash 4.1+
	brew install bash-completion@2
	## you may need add the completion to your completion directory
	jcli completion > $(brew --prefix)/etc/bash_completion.d/jcli
	## If you get trouble, please visit https://github.com/jenkins-zh/jenkins-cli/issues/83.
    
    # Load the jcli completion code for zsh[1] into the current shell
	source <(jcli completion --type zsh)
	# Set the jcli completion code for zsh[1] to autoload on startup
	jcli completion --type zsh > "${fpath[1]}/_jcli"`,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		shellType := completionOptions.Type
		switch shellType {
		case "zsh":
			err = rootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "powerShell":
			err = rootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
		case "bash":
			err = rootCmd.GenBashCompletion(cmd.OutOrStdout())
		default:
			err = fmt.Errorf("unknown shell type %s", shellType)
		}
		return
	},
}

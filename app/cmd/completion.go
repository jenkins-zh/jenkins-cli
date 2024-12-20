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
	"zsh", "bash", "powerShell", "fish",
}

var completionOptions CompletionOptions

func init() {
	rootCmd.AddCommand(completionCmd)

	flags := completionCmd.Flags()
	flags.StringVarP(&completionOptions.Type, "type", "", "bash",
		i18n.T(fmt.Sprintf("Generate different types of shell which are %v", ShellTypes)))

	err := completionCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) (
		i []string, directive cobra.ShellCompDirective) {
		return ShellTypes, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		completionCmd.PrintErrf("register flag type for sub-command doc failed %#v\n", err)
	}
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: i18n.T("Generate shell completion scripts"),
	Long: i18n.T(`Generate shell completion scripts
Normally you don't need to do more extra work to have this feature if you've installed jcli by brew`),
	Example: `  # Installing bash completion on macOS using homebrew
  ## If running Bash 3.2 included with macOS
  brew install bash-completion
  ## or, if running Bash 4.1+
  brew install bash-completion@2
  ## If jcli is installed via homebrew, this should start working immediately.
  ## If you've installed via other means, you may need add the completion to your completion directory
  jcli completion --type bash > $(brew --prefix)/etc/bash_completion.d/jcli


  # Installing bash completion on Linux
  ## If bash-completion is not installed on Linux, please install the 'bash-completion' package
  ## via your distribution's package manager.
  ## Load the jcli completion code for bash into the current shell
  source <(jcli completion --type bash)
  ## Write bash completion code to a file and source if from .bash_profile
  jcli completion --type bash > ~/.jenkins-cli/completion.bash.inc
  printf "
  # jcli shell completion
  source '$HOME/.jenkins-cli/completion.bash.inc'
  " >> $HOME/.bash_profile
  source $HOME/.bash_profile

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
		case "fish":
			err = rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		default:
			err = fmt.Errorf("unknown shell type %s", shellType)
		}
		return
	},
}

package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: i18n.T("Genereate bash completion scripts"),
	Long:  i18n.T("Genereate bash completion scripts"),
	Example: `# Installing bash completion on macOS using homebrew
	## If running Bash 3.2 included with macOS
	brew install bash-completion
	## or, if running Bash 4.1+
	brew install bash-completion@2
	## you may need add the completion to your completion directory
	jcli completion > $(brew --prefix)/etc/bash_completion.d/jcli
	## If you get trouble, please visit https://github.com/jenkins-zh/jenkins-cli/issues/83.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return rootCmd.GenBashCompletion(cmd.OutOrStdout())
	},
}

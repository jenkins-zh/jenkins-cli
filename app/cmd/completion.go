package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Genereate bash completion scripts",
	Long:  `Genereate bash completion scripts`,
	Example: `# Installing bash completion on macOS using homebrew
	## If running Bash 3.2 included with macOS
	brew install bash-completion
	## or, if running Bash 4.1+
	brew install bash-completion@2
	## you may need add the completion to your completion directory
	kubectl completion bash > $(brew --prefix)/etc/bash_completion.d/kubectl
	## If you get trouble, please visit https://github.com/jenkins-zh/jenkins-cli/issues/83.`,
	Run: func(cmd *cobra.Command, _ []string) {
		rootCmd.GenBashCompletion(cmd.OutOrStdout())
	},
}

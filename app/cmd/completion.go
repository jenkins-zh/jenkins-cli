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
	Long:  `To load completion run:
	jcli completion >> ~.bash_completion

If you get trouble, please visit https://github.com/jenkins-zh/jenkins-cli/issues/83.
`,
	Run: func(cmd *cobra.Command, _ []string) {
		rootCmd.GenBashCompletion(cmd.OutOrStdout())
	},
}

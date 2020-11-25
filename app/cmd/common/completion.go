package common

import "github.com/spf13/cobra"

// NoFileCompletion avoid completion with files
func NoFileCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

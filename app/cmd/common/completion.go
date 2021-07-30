package common

import "github.com/spf13/cobra"

// NoFileCompletion avoid completion with files
func NoFileCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// ArrayCompletion return a completion  which base on an array
func ArrayCompletion(array ...string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return array, cobra.ShellCompDirectiveNoFileComp
	}
}

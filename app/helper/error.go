package helper

import "github.com/spf13/cobra"

// CheckErr print a friendly error message
func CheckErr(cmd *cobra.Command, err error) {
	switch {
	case err == nil:
		return
	default:
		cmd.PrintErrln(err)
	}
}

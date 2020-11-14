package common

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// ExistsRegularFile returns a function to check if target file is a regular file
func ExistsRegularFile(flagName string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) (err error) {
		flag := cmd.Flag(flagName)
		val := flag.Value.String()

		switch val {
		case "":
			err = fmt.Errorf("argument '%s' cannot be empty", flagName)
		default:
			if _, err = os.Stat(val); os.IsNotExist(err) {
				err = fmt.Errorf("'%s' is not a regular file", val)
			}
		}
		return
	}
}

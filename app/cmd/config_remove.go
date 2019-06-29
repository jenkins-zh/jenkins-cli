package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configRemoveCmd)
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a Jenkins config",
	Long:  `Remove a Jenkins config`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("You need to give a name")
		}

		target := args[0]
		if err := removeJenkins(target); err != nil {
			log.Fatal(err)
		}
	},
}

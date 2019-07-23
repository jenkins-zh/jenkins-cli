package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	pluginCmd.AddCommand(pluginOpenCmd)
}

var pluginOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Openout update center server",
	Long:  `Openout update center server`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()

		if jenkins.URL != "" {
			open(fmt.Sprintf("%s/pluginManager", jenkins.URL))
		} else {
			log.Fatal(fmt.Sprintf("No URL fond from %s", jenkins.Name))
		}
	},
}

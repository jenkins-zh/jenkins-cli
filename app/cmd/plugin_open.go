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
	Short: "Open update center server in browser",
	Long:  `Open update center server in browser`,
	Run: func(_ *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()

		if jenkins.URL != "" {
			open(fmt.Sprintf("%s/pluginManager", jenkins.URL))
		} else {
			log.Fatal(fmt.Sprintf("No URL fond from %s", jenkins.Name))
		}
	},
}

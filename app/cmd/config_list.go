package cmd

import (
	"fmt"
	"os"

	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configListCmd)
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Jenkins config items",
	Long:  `List all Jenkins config items`,
	Run: func(_ *cobra.Command, _ []string) {
		current := getCurrentJenkins()

		table := util.CreateTable(os.Stdout)
		table.AddRow("number", "name", "url", "description")
		for i, jenkins := range getConfig().JenkinsServers {
			name := jenkins.Name
			if name == current.Name {
				name = fmt.Sprintf("*%s", name)
			}
			if len(jenkins.Description) > 15 {
				table.AddRow(fmt.Sprintf("%d", i), name, jenkins.URL, jenkins.Description[0:15])
			} else {
				table.AddRow(fmt.Sprintf("%d", i), name, jenkins.URL, jenkins.Description)
			}
		}
		table.Render()
	},
}

package cmd

import (
	"bytes"
	"fmt"

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
	Run: func(cmd *cobra.Command, _ []string) {
		current := getCurrentJenkins()

		buf := new(bytes.Buffer)
		table := util.CreateTable(buf)
		table.AddRow("number", "name", "url", "description")
		for i, jenkins := range getConfig().JenkinsServers {
			name := jenkins.Name
			if name == current.Name {
				name = fmt.Sprintf("*%s", name)
			}
			if len(jenkins.Description) > 15 {
				jenkins.Description = jenkins.Description[0:15]
			}
			table.AddRow(fmt.Sprintf("%d", i), name, jenkins.URL, jenkins.Description)
		}
		table.Render()
		cmd.Print(string(buf.Bytes()))
	},
}

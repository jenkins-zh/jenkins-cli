package cmd

import (
	"bytes"
	"fmt"

	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// ConfigListOption option for config list command
type ConfigListOption struct {
	OutputOption
}

var configListOption ConfigListOption

func init() {
	configCmd.AddCommand(configListCmd)
	configListOption.SetFlag(configListCmd)
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("List all Jenkins config items"),
	Long:  i18n.T("List all Jenkins config items"),
	Run: func(cmd *cobra.Command, _ []string) {
		current := getCurrentJenkins()

		data, err := configListOption.Output(current)
		cmd.Print(string(data))
		helper.CheckErr(cmd, err)
	},
}

// Output render data into byte array as a table format
func (o *ConfigListOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil && o.Format == TableOutputFormat {
		current := obj.(*JenkinsServer)

		buf := new(bytes.Buffer)
		table := util.CreateTableWithHeader(buf, o.WithoutHeaders)
		table.AddHeader("number", "name", "url", "description")
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
		err = nil
		data = buf.Bytes()
	}
	return
}

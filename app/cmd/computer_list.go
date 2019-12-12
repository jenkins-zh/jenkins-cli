package cmd

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// ComputerListOption option for config list command
type ComputerListOption struct {
	OutputOption

	RoundTripper http.RoundTripper
}

var computerListOption ComputerListOption

func init() {
	computerCmd.AddCommand(computerListCmd)
	computerListOption.SetFlag(computerListCmd)
}

var computerListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("List all Jenkins agents"),
	Long:  i18n.T("List all Jenkins agents"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.ComputerClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: computerListOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		var computers client.ComputerList
		if computers, err = jClient.List(); err == nil {
			var data []byte
			data, err = computerListOption.Output(computers.Computer)
			if err == nil && len(data) > 0 {
				cmd.Print(string(data))
			}
		}
		return
	},
}

// Output render data into byte array as a table format
func (o *ComputerListOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil && o.Format == TableOutputFormat {
		computers := obj.([]client.Computer)

		buf := new(bytes.Buffer)
		table := util.CreateTableWithHeader(buf, o.WithoutHeaders)
		table.AddHeader("number", "name", "executors", "description", "offline")
		for i, computer := range computers {
			table.AddRow(fmt.Sprintf("%d", i), computer.DisplayName,
				fmt.Sprintf("%d", computer.NumExecutors), computer.Description,
				colorOffline(computer.Offline))
		}
		table.Render()
		err = nil
		data = buf.Bytes()
	}
	return
}

func colorOffline(offline bool) string {
	if offline {
		return util.ColorWarning("yes")
	}
	return "no"
}

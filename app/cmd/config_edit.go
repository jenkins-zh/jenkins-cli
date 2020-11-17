package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"io/ioutil"

	"github.com/spf13/cobra"
)

// ConfigEditOption is the option for edit config command
type ConfigEditOption struct {
	common.Option
}

var configEditOption ConfigEditOption

func init() {
	configCmd.AddCommand(configEditCmd)
	configEditOption.Stdio = common.GetSystemStdio()
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: i18n.T("Edit a Jenkins config"),
	Long: i18n.T(fmt.Sprintf(`Edit a Jenkins config
%s`, common.GetEditorHelpText())),
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		current := getCurrentJenkinsFromOptions()
		configPath := configOptions.ConfigFileLocation

		var data []byte
		if data, err = ioutil.ReadFile(configPath); err == nil {
			content := string(data)
			//Help:          fmt.Sprintf("Config file path: %s", configPath),
			configEditOption.EditFileName = ".jenkins-cli.yaml"
			content, err = configEditOption.Editor(content, fmt.Sprintf("Edit config item %s", current.Name))
			if err == nil {
				err = ioutil.WriteFile(configPath, []byte(content), 0644)
			}
		}
		return
	},
}

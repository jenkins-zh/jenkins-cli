package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	cascCmd.AddCommand(cascOpenCmd)
}

var cascOpenCmd = &cobra.Command{
	Use:   "open",
	Short: i18n.T("Open Configuration as Code page in browser"),
	Long:  i18n.T("Open Configuration as Code page in browser"),
	RunE: func(_ *cobra.Command, _ []string) error {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		return Open(fmt.Sprintf("%s/configuration-as-code", jenkins.URL), exec.Command)
	},
}

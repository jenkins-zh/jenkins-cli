package cmd

import (
	"fmt"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// CASCOpenOption is the option of casc open cmd
type CASCOpenOption struct {
	ExecContext util.ExecContext
}

var cascOpenOption CASCOpenOption

func init() {
	cascCmd.AddCommand(cascOpenCmd)
}

var cascOpenCmd = &cobra.Command{
	Use:   "open",
	Short: i18n.T("Open Configuration as Code page in browser"),
	Long:  i18n.T("Open Configuration as Code page in browser"),
	RunE: func(_ *cobra.Command, _ []string) error {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		return util.Open(fmt.Sprintf("%s/configuration-as-code", jenkins.URL), cascOpenOption.ExecContext)
	},
}

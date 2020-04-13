package cmd

import (
	"fmt"
	"os"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// CASCOpenOption is the option of casc open cmd
type CASCOpenOption struct {
	ExecContext util.ExecContext

	Browser string
}

var cascOpenOption CASCOpenOption

func init() {
	cascCmd.AddCommand(cascOpenCmd)
	cascOpenCmd.Flags().StringVarP(&cascOpenOption.Browser, "browser", "b", "",
		i18n.T("Open Jenkins with a specific browser"))
}

var cascOpenCmd = &cobra.Command{
	Use:   "open",
	Short: i18n.T("Open Configuration as Code page in browser"),
	Long:  i18n.T("Open Configuration as Code page in browser"),
	PreRun: func(_ *cobra.Command, _ []string) {
		if cascOpenOption.Browser == "" {
			cascOpenOption.Browser = os.Getenv("BROWSER")
		}
	},
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		jenkins := getCurrentJenkinsFromOptions()
		if jenkins == nil {
			err = fmt.Errorf("cannot found Jenkins by %s", rootOptions.Jenkins)
			return
		}

		browser := cascOpenOption.Browser
		err = util.Open(fmt.Sprintf("%s/configuration-as-code", jenkins.URL), browser, cascOpenOption.ExecContext)
		return
	},
	Annotations: map[string]string{
		since: "v0.0.24",
	},
}

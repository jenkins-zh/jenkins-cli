package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// PluginCheckoutOption is the option for plugin checkout command
type PluginCheckoutOption struct {
	RoundTripper http.RoundTripper
	// Timeout is the timeout setting for check Jenkins update-center
	Timeout int64
}

var pluginCheckoutOption PluginCheckoutOption

func init() {
	pluginCmd.AddCommand(pluginCheckCmd)

	flags := pluginCheckCmd.Flags()
	flags.Int64VarP(&pluginCheckoutOption.Timeout, "timeout", "", 30,
		"Timeout in second setting for checking Jenkins update-center")
}

var pluginCheckCmd = &cobra.Command{
	Use:   "check",
	Short: i18n.T("Check update center server"),
	Long:  i18n.T(`Check update center server`),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginCheckoutOption.RoundTripper,
				Timeout:      time.Duration(pluginCheckoutOption.Timeout) * time.Second,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		err = jClient.CheckUpdate(func(response *http.Response) {
			code := response.StatusCode
			if code != 200 {
				contentData, _ := ioutil.ReadAll(response.Body)
				cmd.PrintErrln(fmt.Sprintf("response code is %d, content: %s", code, string(contentData)))
			}
		})
		return
	},
}

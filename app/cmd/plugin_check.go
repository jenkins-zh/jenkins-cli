package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// PluginCheckoutOption is the option for plugin checkout command
type PluginCheckoutOption struct {
	RoundTripper http.RoundTripper
}

var pluginCheckoutOption PluginCheckoutOption

func init() {
	pluginCmd.AddCommand(pluginCheckCmd)
}

var pluginCheckCmd = &cobra.Command{
	Use:   "check",
	Short: i18n.T("Check update center server"),
	Long:  i18n.T(`Check update center server`),
	Run: func(cmd *cobra.Command, _ []string) {
		jClient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginCheckoutOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClientOrDie(&(jClient.JenkinsCore))

		err := jClient.CheckUpdate(func(response *http.Response) {
			code := response.StatusCode
			if code != 200 {
				contentData, _ := ioutil.ReadAll(response.Body)
				cmd.PrintErrln(fmt.Sprintf("response code is %d, content: %s", code, string(contentData)))
			}
		})
		helper.CheckErr(cmd, err)
	},
}

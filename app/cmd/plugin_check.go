package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	pluginCmd.AddCommand(pluginCheckCmd)
}

var pluginCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checkout update center server",
	Long:  `Checkout update center server`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		jclient.CheckUpdate(func(response *http.Response) {
			code := response.StatusCode
			if code == 200 {
				fmt.Println("update site updated.")
			} else {
				contentData, _ := ioutil.ReadAll(response.Body)
				log.Fatal(fmt.Sprintf("response code is %d, content: %s",
					code, string(contentData)))
			}
		})
	},
}

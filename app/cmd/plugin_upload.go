package cmd

import (
	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	pluginCmd.AddCommand(pluginUploadCmd)
}

var pluginUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload the plugin from local to your Jenkins",
	Long:  `Upload the plugin from local to your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		jclient.Upload()
	},
}

package cmd

import (
	"log"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
)

func init() {
	pluginCmd.AddCommand(pluginUninstallCmd)
}

var pluginUninstallCmd = &cobra.Command{
	Use:   "uninstall [pluginName]",
	Short: "Uninstall the plugins",
	Long:  `Uninstall the plugins`,
	Run: func(cmd *cobra.Command, args []string) {
		var pluginName string
		if len(args) == 0 {
			cmd.Help()
			return
		}

		pluginName = args[0]

		jenkins := getCurrentJenkins()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if err := jclient.UninstallPlugin(pluginName); err != nil {
			log.Fatal(err)
		}
	},
}

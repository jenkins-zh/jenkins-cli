package cmd

import (
	"github.com/sladyn98/Jenkins/jenkins-cli/client"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type PluginNewOption struct {
	OutputOption

	New []string

	RoundTripper http.RoundTripper
}

var pluginNewOption PluginNewOption

func init() {
	pluginCmd.AddCommand(pluginNewCmd)
	pluginListCmd.Flags().StringArrayVarP(&pluginNewOption.New,"new", "", []string{}, "List of new plugins")
}


var pluginNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Print all the new plugins",
	Long:  `Print all the new plugins which are available for installation`,

	Run: func(cmd *cobra.Command, _ []string) {
		jclient := &client.PluginAPI{}

		newPluginList,err := jclient.NewPlugins()

		if err!=nil {
			/*
			  Print the list of plugins
			*/
		} else {
				log.Fatal(err)
		}
	},
}

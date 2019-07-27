package cmd

import (
	"fmt"
	"log"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type QueueOption struct {
	OutputOption
}

var queueOption QueueOption

func init() {
	rootCmd.AddCommand(queueCmd)
	queueCmd.PersistentFlags().StringVarP(&queueOption.Format, "output", "o", "json", "Format the output")
}

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Print the queue of your Jenkins",
	Long:  `Print the queue of your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		jclient := &client.QueueClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if status, err := jclient.Get(); err == nil {
			var data []byte
			if data, err = Format(status, queueOption.Format); err == nil {
				fmt.Printf("%s\n", string(data))
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

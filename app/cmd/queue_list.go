package cmd

import (
	"fmt"
	"log"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// QueueListOption represents the option of queue list command
type QueueListOption struct {
	OutputOption
}

var queueListOption QueueListOption

func init() {
	queueCmd.AddCommand(queueListCmd)
	queueListCmd.Flags().StringVarP(&queueListOption.Format, "output", "o", "json", "Format the output")
}

var queueListCmd = &cobra.Command{
	Use:   "list",
	Short: "Print the queue of your Jenkins",
	Long:  `Print the queue of your Jenkins`,
	Run: func(_ *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.QueueClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if status, err := jclient.Get(); err == nil {
			var data []byte
			if data, err = Format(status, queueListOption.Format); err == nil {
				fmt.Printf("%s\n", string(data))
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

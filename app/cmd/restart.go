package cmd

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

type RestartOptions struct {
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

//curl -X POST http://localhost:8080/jenkins/safeRestart

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart your Jenkins",
	Long:  `Restart your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		crumb, config := getCrumb()
		api := fmt.Sprintf("%s/safeRestart", config.URL)

		req, err := http.NewRequest("POST", api, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("Accept", "*/*")
		req.SetBasicAuth(config.UserName, config.Token)
		req.Header.Add(crumb.CrumbRequestField, crumb.Crumb)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		if response, err := client.Do(req); err == nil {
			if response.StatusCode != 200 {
				if data, err := ioutil.ReadAll(response.Body); err == nil {
					fmt.Println(string(data))
				} else {
					log.Fatal(err)
				}
			}
		} else {
			log.Fatal(err)
		}
	},
}

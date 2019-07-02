package cmd

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/AlecAivazis/survey"
	"github.com/spf13/cobra"
)

// RestartOption holds the options for restart cmd
type RestartOption struct {
	BatchOption
}

var restartOption RestartOption

func init() {
	rootCmd.AddCommand(restartCmd)
	restartCmd.Flags().BoolVarP(&restartOption.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart your Jenkins",
	Long:  `Restart your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		crumb, config := getCrumb()

		if !restartOption.Batch {
			confirm := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Are you sure to restart Jenkins %s?", config.URL),
			}
			survey.AskOne(prompt, &confirm)
			if !confirm {
				return
			}
		}

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

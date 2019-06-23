package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

//curl -k -u $JENKINS_USER:$JENKINS_TOKEN $JENKINS_URL/crumbIssuer/api/json -s

// Start contains the command line options
type CrumbIssuerOptions struct {
	Upload bool
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "crumb",
	Short: "Print crumbIssuer of Jenkins",
	Long:  `Print crumbIssuer of Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

type CrumbIssuer struct {
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

func getCrumb() (CrumbIssuer, JenkinsServer) {
	config := getCurrentJenkins()

	jenkinsRoot := config.URL
	api := fmt.Sprintf("%s/crumbIssuer/api/json", jenkinsRoot)

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.UserName, config.Token)

	var crumbIssuer CrumbIssuer
	client := &http.Client{}
	if response, err := client.Do(req); err == nil {
		if data, err := ioutil.ReadAll(response.Body); err == nil {
			fmt.Println(string(data))
			json.Unmarshal(data, &crumbIssuer)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	return crumbIssuer, config
}

func init() {
}

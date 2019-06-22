package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// Start contains the command line options
type CrumbIssuerOptions struct {
	Upload bool
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "crumb",
	Short: "Print the version number of Hugo",
	Long:  `Manage the plugin of Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		//curl -k -u $JENKINS_USER:$JENKINS_TOKEN $JENKINS_URL/crumbIssuer/api/json -s
	},
}

type CrumbIssuer struct {
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

func getCrumb() (CrumbIssuer, Config) {
	config := getConfig()
	fmt.Println(config)

	jenkinsRoot := config.JenkinsServers[0].URL
	api := fmt.Sprintf("%s/crumbIssuer/api/json", jenkinsRoot)

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.JenkinsServers[0].UserName, config.JenkinsServers[0].Token)

	var crumbIssuer CrumbIssuer
	client := &http.Client{}
	if response, err := client.Do(req); err == nil {
		if data, err := ioutil.ReadAll(response.Body); err == nil {
			fmt.Println("crumbe success")
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

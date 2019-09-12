package cmd

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

// CrumbIssuerOptions contains the command line options
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
	Run: func(_ *cobra.Command, _ []string) {
		crumb, _ := getCrumb()
		fmt.Printf("%s=%s", crumb.CrumbRequestField, crumb.Crumb)
	},
}

// CrumbIssuer represents Jenkins crumb
type CrumbIssuer struct {
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

func getCrumb() (CrumbIssuer, *JenkinsServer) {
	config := getCurrentJenkins()

	jenkinsRoot := config.URL
	api := fmt.Sprintf("%s/crumbIssuer/api/json", jenkinsRoot)

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.UserName, config.Token)

	var crumbIssuer CrumbIssuer
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if config.ProxyAuth != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(config.ProxyAuth))
		req.Header.Add("Proxy-Authorization", basicAuth)

		tr.ProxyConnectHeader = http.Header{}
		tr.ProxyConnectHeader.Add("Proxy-Authorization", basicAuth)

		if proxyURL, err := url.Parse(config.Proxy); err == nil {
			tr.Proxy = http.ProxyURL(proxyURL)
		}
	}
	client := &http.Client{Transport: tr}
	if response, err := client.Do(req); err == nil {
		if data, err := ioutil.ReadAll(response.Body); err == nil {
			if response.StatusCode == 200 {
				json.Unmarshal(data, &crumbIssuer)
			} else {
				fmt.Println("get curmb error")
				log.Fatal(string(data))
			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	return crumbIssuer, config
}

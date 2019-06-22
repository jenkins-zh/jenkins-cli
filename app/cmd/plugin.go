package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Start contains the command line options
type PluginOptions struct {
	Upload bool
}

func init() {
	rootCmd.AddCommand(pluginCmd)
}

var author PluginOptions

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Print the version number of Hugo",
	Long:  `Manage the plugin of Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")

		fmt.Printf("upload: %v\n", author.Upload)
		if author.Upload {
			crumb, config := getCrumb()

			fmt.Println("crumb", crumb)

			jenkinsRoot := getConfig().JenkinsServers[0].URL
			api := fmt.Sprintf("%s/pluginManager/uploadPlugin", jenkinsRoot)

			path, _ := os.Getwd()
			path += "/target/alauda-devops-sync.hpi"
			fmt.Println("target path", path)
			extraParams := map[string]string{}
			request, err := newfileUploadRequest(api, extraParams, "@name", path)
			if err != nil {
				log.Fatal(err)
			}
			request.SetBasicAuth(config.JenkinsServers[0].UserName, config.JenkinsServers[0].Token)
			request.Header.Add("Accept", "*/*")
			request.Header.Add(crumb.CrumbRequestField, crumb.Crumb)
			fmt.Println(request.Header)
			if err == nil {
				client := &http.Client{}
				var response *http.Response
				response, err = client.Do(request)
				if err != nil {
					fmt.Println(err)
				} else if response.StatusCode != 200 {
					var data []byte
					if data, err = ioutil.ReadAll(response.Body); err == nil {
						fmt.Println(string(data))
					} else {
						log.Fatal(err)
					}
				}
			} else {
				log.Fatal(err)
			}
		}
	},
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&author.Upload, "upload", false, "Upload plugin to your Jenkins server")
	viper.BindPFlag("upload", rootCmd.PersistentFlags().Lookup("upload"))
}

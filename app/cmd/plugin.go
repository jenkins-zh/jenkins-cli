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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PluginOptions contains the command line options
type PluginOptions struct {
	Upload bool
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.PersistentFlags().BoolVarP(&author.Upload, "upload", "u", false, "Upload plugin to your Jenkins server")
	viper.BindPFlag("upload", pluginCmd.PersistentFlags().Lookup("upload"))
}

var author PluginOptions

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage the plugins of Jenkins",
	Long:  `Manage the plugins of Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if author.Upload {
			crumb, config := getCrumb()

			api := fmt.Sprintf("%s/pluginManager/uploadPlugin", config.URL)

			path, _ := os.Getwd()
			dirName := filepath.Base(path)
			dirName = strings.Replace(dirName, "-plugin", "", -1)
			path += fmt.Sprintf("/target/%s.hpi", dirName)
			fmt.Println("target path", path)
			extraParams := map[string]string{}
			request, err := newfileUploadRequest(api, extraParams, "@name", path)
			if err != nil {
				log.Fatal(err)
			}
			request.SetBasicAuth(config.UserName, config.Token)
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

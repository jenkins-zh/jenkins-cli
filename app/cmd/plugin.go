package cmd

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosuri/uiprogress"
	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PluginOptions contains the command line options
type PluginOptions struct {
	Upload      bool
	CheckUpdate bool
	Open        bool
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.PersistentFlags().BoolVarP(&pluginOpt.Upload, "upload", "u", false, "Upload plugin to your Jenkins server")
	pluginCmd.PersistentFlags().BoolVarP(&pluginOpt.CheckUpdate, "check", "c", false, "Checkout update center server")
	pluginCmd.PersistentFlags().BoolVarP(&pluginOpt.Open, "open", "o", false, "Open the browse with the address of plugin manager")
	viper.BindPFlag("upload", pluginCmd.PersistentFlags().Lookup("upload"))
}

var pluginOpt PluginOptions

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage the plugins of Jenkins",
	Long:  `Manage the plugins of Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		if pluginOpt.Upload {
			crumb, config := getCrumb()

			api := fmt.Sprintf("%s/pluginManager/uploadPlugin", config.URL)

			path, _ := os.Getwd()
			dirName := filepath.Base(path)
			dirName = strings.Replace(dirName, "-plugin", "", -1)
			path += fmt.Sprintf("/target/%s.hpi", dirName)
			extraParams := map[string]string{}
			request, err := newfileUploadRequest(api, extraParams, "@name", path)
			if err != nil {
				log.Fatal(err)
			}
			request.SetBasicAuth(config.UserName, config.Token)
			request.Header.Add("Accept", "*/*")
			request.Header.Add(crumb.CrumbRequestField, crumb.Crumb)
			if err == nil {
				tr := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
				client := &http.Client{Transport: tr}
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

		if pluginOpt.CheckUpdate {
			jclient := &client.PluginManager{}

			crumb, config := getCrumb()
			jclient.CrumbRequestField = crumb.CrumbRequestField
			jclient.Crumb = crumb.Crumb
			jclient.URL = config.URL
			jclient.UserName = config.UserName
			jclient.Token = config.Token
			jclient.CheckUpdate(func(response *http.Response) {
				code := response.StatusCode
				if code == 200 {
					fmt.Println("update site updated.")
				} else {
					contentData, _ := ioutil.ReadAll(response.Body)
					log.Fatal(fmt.Sprintf("response code is %d, content: ",
						code, string(contentData)))
				}
			})
		}

		if pluginOpt.Open {
			jenkins := getCurrentJenkins()
			if jenkins.URL != "" {
				open(fmt.Sprintf("%s/pluginManager", jenkins.URL))
			} else {
				log.Fatal(fmt.Sprintf("No URL fond from %s", jenkins.Name))
			}
		}
	},
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var total float64
	if stat, err := file.Stat(); err != nil {
		panic(err)
	} else {
		total = float64(stat.Size())
	}
	defer file.Close()

	// body := &bytes.Buffer{}
	body := &ProgressIndicator{
		Total: total,
	}
	body.Init()
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

type ProgressIndicator struct {
	bytes.Buffer
	Total float64
	count float64
	bar   *uiprogress.Bar
}

func (i *ProgressIndicator) Init() {
	uiprogress.Start()             // start rendering
	i.bar = uiprogress.AddBar(100) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	i.bar.AppendCompleted()
	// i.bar.PrependElapsed()
}

func (i *ProgressIndicator) Write(p []byte) (n int, err error) {
	n, err = i.Buffer.Write(p)
	return
}

func (i *ProgressIndicator) Read(p []byte) (n int, err error) {
	n, err = i.Buffer.Read(p)
	i.count += float64(n)
	i.bar.Set((int)(i.count * 100 / i.Total))
	return
}

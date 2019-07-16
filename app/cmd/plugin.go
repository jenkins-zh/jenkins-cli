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
	"github.com/linuxsuren/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// PluginOptions contains the command line options
type PluginOptions struct {
	OutputOption

	Upload      bool
	CheckUpdate bool
	Open        bool
	List        bool

	Install   []string
	Uninstall string

	Filter []string
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.Flags().BoolVarP(&pluginOpt.Upload, "upload", "u", false, "Upload plugin to your Jenkins server")
	pluginCmd.Flags().BoolVarP(&pluginOpt.CheckUpdate, "check", "c", false, "Checkout update center server")
	pluginCmd.Flags().BoolVarP(&pluginOpt.Open, "open", "o", false, "Open the browse with the address of plugin manager")
	pluginCmd.Flags().BoolVarP(&pluginOpt.List, "list", "l", false, "Print all the plugins which are installed")
	pluginCmd.Flags().StringVarP(&pluginOpt.Format, "format", "", TableOutputFormat, "Format the output")
	pluginCmd.Flags().StringArrayVarP(&pluginOpt.Install, "install", "", []string{}, "Install a plugin by shortName")
	pluginCmd.Flags().StringVarP(&pluginOpt.Uninstall, "uninstall", "", "", "Uninstall a plugin by shortName")
	pluginCmd.Flags().StringArrayVarP(&pluginOpt.Filter, "filter", "", []string{}, "Filter for the list, like: active, hasUpdate, downgradable, enable, name=foo")
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

		config := getCurrentJenkins()
		jenkins := getCurrentJenkins()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = config.UserName
		jclient.Token = config.Token
		if pluginOpt.CheckUpdate {
			jclient.CheckUpdate(func(response *http.Response) {
				code := response.StatusCode
				if code == 200 {
					fmt.Println("update site updated.")
				} else {
					contentData, _ := ioutil.ReadAll(response.Body)
					log.Fatal(fmt.Sprintf("response code is %d, content: %s",
						code, string(contentData)))
				}
			})
		}

		if pluginOpt.Open {
			if jenkins.URL != "" {
				open(fmt.Sprintf("%s/pluginManager", jenkins.URL))
			} else {
				log.Fatal(fmt.Sprintf("No URL fond from %s", jenkins.Name))
			}
		}

		if pluginOpt.List {
			var (
				filter       bool
				hasUpdate    bool
				downgradable bool
				enable       bool
				active       bool
				pluginName   string
			)
			if pluginOpt.Filter != nil {
				filter = true
				for _, f := range pluginOpt.Filter {
					switch f {
					case "hasUpdate":
						hasUpdate = true
					case "downgradable":
						downgradable = true
					case "enable":
						enable = true
					case "active":
						active = true
					case "name":
						downgradable = true
					}

					if strings.HasPrefix(f, "name=") {
						pluginName = strings.TrimPrefix(f, "name=")
					}
				}
			}

			if plugins, err := jclient.GetPlugins(); err == nil {
				filteredPlugins := make([]client.InstalledPlugin, 0)
				for _, plugin := range plugins.Plugins {
					if filter {
						if hasUpdate && !plugin.HasUpdate {
							continue
						}

						if downgradable && !plugin.Downgradable {
							continue
						}

						if enable && !plugin.Enable {
							continue
						}

						if active && !plugin.Active {
							continue
						}

						if pluginName != "" && !strings.Contains(plugin.ShortName, pluginName) {
							continue
						}

						filteredPlugins = append(filteredPlugins, plugin)
					}
				}

				if data, err := pluginOpt.Output(filteredPlugins); err == nil {
					if len(data) > 0 {
						fmt.Println(string(data))
					}
				} else {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}

		if pluginOpt.Install != nil && len(pluginOpt.Install) > 0 {
			if err := jclient.InstallPlugin(pluginOpt.Install); err != nil {
				log.Fatal(err)
			}
		}

		if pluginOpt.Uninstall != "" {
			if err := jclient.UninstallPlugin(pluginOpt.Uninstall); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func (o *PluginOptions) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		pluginList := obj.([]client.InstalledPlugin)
		table := util.CreateTable(os.Stdout)
		table.AddRow("number", "name", "version", "update")
		for i, plugin := range pluginList {
			table.AddRow(fmt.Sprintf("%d", i), plugin.ShortName, plugin.Version, fmt.Sprintf("%v", plugin.HasUpdate))
		}
		table.Render()
		err = nil
		data = []byte{}
	}
	return
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

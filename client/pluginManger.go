package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/util"
)

type PluginManager struct {
	JenkinsCore
}

type Plugin struct {
	Active       bool
	Enabled      bool
	Bundled      bool
	Downgradable bool
	Deleted      bool
}

// PluginList represent a list of plugins
type InstalledPluginList struct {
	Plugins []InstalledPlugin
}

type AvailablePluginList struct {
	Data   []AvailablePlugin
	Status string
}

type AvailablePlugin struct {
	Plugin

	// for the available list
	Name      string
	Installed bool
	Website   string
	Title     string
}

// InstalledPlugin represent the installed plugin from Jenkins
type InstalledPlugin struct {
	Plugin

	Enable             bool
	ShortName          string
	LongName           string
	Version            string
	URL                string
	HasUpdate          bool
	Pinned             bool
	RequiredCoreVesion string
	MinimumJavaVersion string
	SupportDynamicLoad string
	BackVersion        string
}

// CheckUpdate fetch the lastest plugins from update center site
func (p *PluginManager) CheckUpdate(handle func(*http.Response)) {
	api := fmt.Sprintf("%s/pluginManager/checkUpdatesServer", p.URL)
	req, err := http.NewRequest("POST", api, nil)
	if err == nil {
		p.AuthHandle(req)
	} else {
		log.Fatal(err)
	}

	if err = p.CrumbHandle(req); err != nil {
		log.Fatal(err)
	}

	client := p.GetClient()
	if response, err := client.Do(req); err == nil {
		p.handleCheck(handle)(response)
	} else {
		log.Fatal(err)
	}
}

// GetAvailablePlugins get the aviable plugins from Jenkins
func (p *PluginManager) GetAvailablePlugins() (pluginList *AvailablePluginList, err error) {
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = p.Request("GET", "/pluginManager/plugins", nil, nil); err == nil {
		if statusCode == 200 {
			pluginList = &AvailablePluginList{}
			err = json.Unmarshal(data, pluginList)
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
			if p.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
	}
	return
}

// GetPlugins get installed plugins
func (p *PluginManager) GetPlugins() (pluginList *InstalledPluginList, err error) {
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = p.Request("GET", "/pluginManager/api/json?depth=1", nil, nil); err == nil {
		if statusCode == 200 {
			pluginList = &InstalledPluginList{}
			err = json.Unmarshal(data, pluginList)
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
			if p.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
	}
	return
}

func getPluginsInstallQuery(names []string) string {
	pluginNames := make([]string, 0)
	for _, name := range names {
		if name == "" {
			continue
		}
		pluginNames = append(pluginNames, fmt.Sprintf("plugin.%s=", name))
	}
	return strings.Join(pluginNames, "&")
}

// InstallPlugin install a plugin by name
func (p *PluginManager) InstallPlugin(names []string) (err error) {
	api := fmt.Sprintf("%s/pluginManager/install?%s", p.URL, getPluginsInstallQuery(names))
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("POST", api, nil)
	if err == nil {
		if err = p.AuthHandle(req); err != nil {
			log.Fatal(err)
		}
	} else {
		return
	}

	client := p.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			fmt.Println("install succeed.")
		} else if code == 400 {
			if errMsg, ok := response.Header["X-Error"]; ok {
				for _, msg := range errMsg {
					fmt.Println(msg)
				}
			} else {
				fmt.Println("Cannot found plugins", names)
			}
		} else {
			fmt.Println(response.Header)
			fmt.Println("status code", code)
			if err == nil && p.Debug && len(data) > 0 {
				ioutil.WriteFile("debug.html", data, 0664)
			} else if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal(err)
	}
	return
}

// UninstallPlugin uninstall a plugin by name
func (p *PluginManager) UninstallPlugin(name string) (err error) {
	api := fmt.Sprintf("/pluginManager/plugin/%s/doUninstall", name)
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = p.Request("POST", api, nil, nil); err == nil {
		if statusCode != 200 {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
			if p.Debug {
				ioutil.WriteFile("debug.html", data, 0664)
			}
		}
	}
	return
}

// Upload will upload a file from local filesystem into Jenkins
func (p *PluginManager) Upload(pluginFile string) {
	api := fmt.Sprintf("%s/pluginManager/uploadPlugin", p.URL)
	extraParams := map[string]string{}
	request, err := newfileUploadRequest(api, extraParams, "@name", pluginFile)
	if err != nil {
		log.Fatal(err)
	}

	if err == nil {
		p.AuthHandle(request)
	} else {
		return
	}

	client := p.GetClient()
	var response *http.Response
	response, err = client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else if response.StatusCode != 200 {
		fmt.Println("StatusCode", response.StatusCode)
		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil && p.Debug {
			ioutil.WriteFile("debug.html", data, 0664)
		} else {
			log.Fatal(err)
		}
	}
}

func (p *PluginManager) handleCheck(handle func(*http.Response)) func(*http.Response) {
	if handle == nil {
		handle = func(*http.Response) {
			// Do nothing, just for avoid nil exception
		}
	}
	return handle
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

	bytesBuffer := &bytes.Buffer{}
	progressWriter := &util.ProgressIndicator{
		Total:  total,
		Writer: bytesBuffer,
		Reader: bytesBuffer,
		Title:  "Uploading",
	}
	progressWriter.Init()
	writer := multipart.NewWriter(bytesBuffer)
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

	req, err := http.NewRequest("POST", uri, progressWriter)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

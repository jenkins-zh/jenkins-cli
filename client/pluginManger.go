package client

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

	"github.com/jenkins-zh/jenkins-cli/util"
)

// PluginManager is the client of plugin manager
type PluginManager struct {
	JenkinsCore

	ShowProgress bool
}

// Plugin represents a plugin of Jenkins
type Plugin struct {
	Active       bool
	Enabled      bool
	Bundled      bool
	Downgradable bool
	Deleted      bool
}

// InstalledPluginList represent a list of plugins
type InstalledPluginList struct {
	Plugins []InstalledPlugin
}

// AvailablePluginList represents a list of available plugins
type AvailablePluginList struct {
	Data   []AvailablePlugin
	Status string
}

// AvailablePlugin represetns a available plugin
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
	Dependencies       []Dependence
}

// Dependence represent the plugin's package dependence
type Dependence struct {
	Optional  bool
	ShortName string
	Version   string
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
	err = p.RequestWithData("GET", "/pluginManager/plugins", nil, nil, 200, &pluginList)
	return
}

// GetPlugins get installed plugins
func (p *PluginManager) GetPlugins() (pluginList *InstalledPluginList, err error) {
	err = p.RequestWithData("GET", "/pluginManager/api/json?depth=2", nil, nil, 200, &pluginList)
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
	api := fmt.Sprintf("/pluginManager/install?%s", getPluginsInstallQuery(names))
	_, err = p.RequestWithoutData("POST", api, nil, nil, 200)

	// TODO needs to consider the following cases
	// code == 400 {
	// 	if errMsg, ok := response.Header["X-Error"]; ok {
	// 		for _, msg := range errMsg {
	// 			fmt.Println(msg)
	// 		}
	// 	} else {
	// 		fmt.Println("Cannot found plugins", names)
	// 	}
	// }
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
func (p *PluginManager) Upload(pluginFile string) (err error) {
	api := fmt.Sprintf("%s/pluginManager/uploadPlugin", p.URL)
	extraParams := map[string]string{}
	var request *http.Request
	if request, err = p.newfileUploadRequest(api, extraParams, "@name", pluginFile); err != nil {
		return
	}

	p.AuthHandle(request)

	client := p.GetClient()
	var response *http.Response
	if response, err = client.Do(request); err != nil {
		return
	} else if response.StatusCode != 200 {
		err = fmt.Errorf("StatusCode: %d", response.StatusCode)
		if data, readErr := ioutil.ReadAll(response.Body); readErr == nil && p.Debug {
			ioutil.WriteFile("debug.html", data, 0664)
		}
	}
	return err
}

func (p *PluginManager) handleCheck(handle func(*http.Response)) func(*http.Response) {
	if handle == nil {
		handle = func(*http.Response) {
			// Do nothing, just for avoid nil exception
		}
	}
	return handle
}

func (p *PluginManager) newfileUploadRequest(uri string, params map[string]string, paramName, path string) (req *http.Request, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}

	var total float64
	var stat os.FileInfo
	if stat, err = file.Stat(); err != nil {
		return
	}
	total = float64(stat.Size())
	defer file.Close()

	bytesBuffer := &bytes.Buffer{}
	writer := multipart.NewWriter(bytesBuffer)

	var part io.Writer
	part, err = writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return
	}

	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return
	}

	var progressWriter *util.ProgressIndicator
	if p.ShowProgress {
		progressWriter = &util.ProgressIndicator{
			Total:  total,
			Writer: bytesBuffer,
			Reader: bytesBuffer,
			Title:  "Uploading",
		}
		progressWriter.Init()
		req, err = http.NewRequest("POST", uri, progressWriter)
	} else {
		req, err = http.NewRequest("POST", uri, bytesBuffer)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return
}

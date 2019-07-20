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

	"github.com/gosuri/uiprogress"
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

type PluginDependency struct {
	Name     string `json:"name"`
	Implied  bool   `json:"implied"`
	Optional bool   `json:"optional"`
	Title    string `json:"title"`
	Version  string `json:"version"`
}

type PluginInfo struct {
	BuildDate         string             `json:"buildDate"`
	Dependencies      []PluginDependency `json:"dependencies"`
	Excerpt           string             `json:"excerpt"`
	FirstRelease      string             `json:"firstRelease"`
	Gav               string             `json:"gav"`
	Name              string             `json:"name"`
	PreviousTimestamp string             `json:"previousTimestamp"`
	PreviousVersion   string             `json:"previousVersion"`
	ReleaseTimestamp  string             `json:"releaseTimestamp"`
	RequireCore       string             `json:"RequireCore"`
	Title             string             `json:"title"`
	URL               string             `json:"url"`
	Version           string             `json:"version"`
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

func (p *PluginManager) GetAvailablePlugins() (pluginList *AvailablePluginList, err error) {
	api := fmt.Sprintf("%s/pluginManager/plugins", p.URL)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("GET", api, nil)
	if err == nil {
		p.AuthHandle(req)
	} else {
		return
	}

	client := p.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			if err == nil {
				pluginList = &AvailablePluginList{}
				err = json.Unmarshal(data, pluginList)
			}
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (p *PluginManager) GetPlugins() (pluginList *InstalledPluginList, err error) {
	api := fmt.Sprintf("%s/pluginManager/api/json?pretty=true&depth=1", p.URL)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("GET", api, nil)
	if err == nil {
		p.AuthHandle(req)
	} else {
		return
	}

	client := p.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			if err == nil {
				pluginList = &InstalledPluginList{}
				err = json.Unmarshal(data, pluginList)
			}
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

// InstallPlugin install a plugin by name
func (p *PluginManager) InstallPlugin(names []string) (err error) {
	for i, name := range names {
		names[i] = fmt.Sprintf("plugin.%s", name)
	}
	api := fmt.Sprintf("%s/pluginManager/install?%s", p.URL, strings.Join(names, "=&"))
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
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

// UninstallPlugin uninstall a plugin by name
func (p *PluginManager) UninstallPlugin(name string) (err error) {
	api := fmt.Sprintf("%s/pluginManager/plugin/%s/uninstall", p.URL, name)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("POST", api, nil)
	if err == nil {
		p.AuthHandle(req)
	} else {
		return
	}

	client := p.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			fmt.Println("uninstall succeed.")
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (p *PluginManager) Upload() {
	api := fmt.Sprintf("%s/pluginManager/uploadPlugin", p.URL)

	path, _ := os.Getwd()
	dirName := filepath.Base(path)
	dirName = strings.Replace(dirName, "-plugin", "", -1)
	path += fmt.Sprintf("/target/%s.hpi", dirName)
	extraParams := map[string]string{}
	request, err := newfileUploadRequest(api, extraParams, "@name", path)
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
		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			fmt.Println(string(data))
		} else {
			log.Fatal(err)
		}
	}
}

func (p *PluginManager) handleCheck(handle func(*http.Response)) func(*http.Response) {
	if handle == nil {
		handle = func(*http.Response) {}
	}
	return handle
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

type PluginDownloader struct {
	dependencyMap map[string]string
}

// DownloadPlugins will download those plugins from update center
func (d *PluginDownloader) DownloadPlugins(names []string) {
	d.dependencyMap = make(map[string]string)
	fmt.Println("Start to collect plugin dependencies...")
	plugins := make([]PluginInfo, 0)
	for _, name := range names {
		plugins = append(plugins, d.collectDependencies(name)...)
	}

	fmt.Printf("Ready to download plugins, total: %d.\n", len(plugins))
	for i, plugin := range plugins {
		fmt.Printf("Start to download plugin %s, version: %s, number: %d\n",
			plugin.Name, plugin.Version, i)

		d.download(plugin.URL, plugin.Name)
	}
}

func (d *PluginDownloader) download(url string, name string) {
	if resp, err := http.Get(url); err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			fmt.Println(err)
		} else {
			if err = ioutil.WriteFile(fmt.Sprintf("%s.hpi", name), body, 0644); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (d *PluginDownloader) getPlugin(name string) (plugin *PluginInfo, err error) {
	resp, err := http.Get("https://plugins.jenkins.io/api/plugin/" + name)
	if err != nil {
		return plugin, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	plugin = &PluginInfo{}
	err = json.Unmarshal(body, plugin)
	if err != nil {
		log.Println("error when unmarshal:", string(body))
	}
	return
}

func (d *PluginDownloader) collectDependencies(pluginName string) (plugins []PluginInfo) {
	plugin, err := d.getPlugin(pluginName)
	if err != nil {
		log.Println("can't get the plugin by name:", pluginName)
		panic(err)
	}

	plugins = make([]PluginInfo, 0)
	plugins = append(plugins, *plugin)

	for _, dependency := range plugin.Dependencies {
		if dependency.Optional {
			continue
		}
		if _, ok := d.dependencyMap[dependency.Name]; !ok {
			d.dependencyMap[dependency.Name] = dependency.Version

			plugins = append(plugins, d.collectDependencies(dependency.Name)...)
		}
	}
	return
}

package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type PluginList struct {
	 List []PluginInfo
}

type PluginAPI struct {
	dependencyMap map[string]string
}

type PluginDependency struct {
	Name     string `json:"name"`
	Implied  bool   `json:"implied"`
	Optional bool   `json:"optional"`
	Title    string `json:"title"`
	Version  string `json:"version"`
}

type MaintainerInfo struct {
	ID		string  `json:"id"`
	Name 	string  `json:"name"`
	Email   string  `json:"email"`
}

type WikiInfo struct {
	Content 	string `json:"content"`
	URL         string `json:"url"`
}

type ScmInfo struct {
	Link			  	string `json:"link"`
	InLatestRelease   	string `json:"inLatestRelease"`
	SinceLatestRelease  string `json:"sinceLatestRelease"`
	PullRequests        string `json:"pullRequests"`
}

// NewPluginInfo hold the info of the new plugin
type NewPluginInfo struct {
	Wiki              []WikiInfo         `json:"wiki"`
	Sha1              string             `json:"sha1"`
	RequiredCore      string             `json:"requiredCore"`
	Maintainers       []MaintainerInfo   `json:"maintainers"`
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
	SecurityWarnings  string             `json:"securityWarnings"`
	Scm               []ScmInfo          `json:"scm"`

	Stats PluginInfoStats
}

// PluginInfo hold the info of a plugin
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

	Stats PluginInfoStats
}

type PluginInfoStats struct {
	CurrentInstalls                   int
	Installations                     []PluginInstallationInfo
	InstallationsPerVersion           []PluginInstallationInfo
	InstallationsPercentage           []PluginInstallationInfo
	InstallationsPercentagePerVersion []PluginInstallationInfo
	Trend                             int
}

type PluginInstallationInfo struct {
	Timestamp  int64
	Total      int
	Version    string
	Percentage float64
}

func (a *PluginAPI) ShowTrend(name string) {
	if plugin, err := a.getPlugin(name); err == nil {
		data := []float64{}
		installations := plugin.Stats.Installations
		offset, count := 0, 10
		if len(installations) > count {
			offset = len(installations) - count
		}
		for _, installation := range installations[offset:] {
			data = append(data, float64(installation.Total))
		}

		min, max := 0.0, 0.0
		for _, item := range data {
			if item < min {
				min = item
			} else if item > max {
				max = item
			}
		}

		unit := (max - min) / 100
		for _, num := range data {
			total := (int)(num / unit)
			if total == 0 {
				total = 1
			}
			arr := make([]int, total)
			for _ = range arr {
				fmt.Print("*")
			}
			fmt.Println("", num)
		}
	} else {
		log.Fatal(err)
	}
}

// DownloadPlugins will download those plugins from update center
func (d *PluginAPI) DownloadPlugins(names []string) {
	d.dependencyMap = make(map[string]string)
	fmt.Println("Start to collect plugin dependencies...")
	plugins := make([]PluginInfo, 0)
	for _, name := range names {
		plugins = append(plugins, d.collectDependencies(strings.ToLower(name))...)
	}

	fmt.Printf("Ready to download plugins, total: %d.\n", len(plugins))
	for i, plugin := range plugins {
		fmt.Printf("Start to download plugin %s, version: %s, number: %d\n",
			plugin.Name, plugin.Version, i)

		d.download(plugin.URL, plugin.Name)
	}
}

func (d *PluginAPI) download(url string, name string) {
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

func (d *PluginAPI) getPlugin(name string) (plugin *PluginInfo, err error) {
	var cli = http.Client{}
	cli.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	resp, err := cli.Get("https://plugins.jenkins.io/api/plugin/" + name)
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

func (d *PluginAPI) collectDependencies(pluginName string) (plugins []PluginInfo) {
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

// New Plugins will list all the new plugins that can be installed.
func (d *PluginAPI) NewPlugins() (pluginList *PluginList){
	var cli= http.Client{}
	cli.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	resp, _ := cli.Get("https://plugins.jenkins.io/api/plugins/new")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, pluginList)
	if err != nil {
		log.Println("error when unmarshal:", string(body))
	}
	return
}

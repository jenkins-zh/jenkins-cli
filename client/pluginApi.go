package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
		installations := plugin.Stats.Installations
		offset := 0
		if len(installations) > 5 {
			offset = len(installations) - 5
		}
		for _, installation := range installations[offset:] {
			fmt.Println(installation.Total)
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
		plugins = append(plugins, d.collectDependencies(name)...)
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

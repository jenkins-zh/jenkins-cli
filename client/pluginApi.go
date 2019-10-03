package client

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"net/http"
	"strings"
	"github.com/jenkins-zh/jenkins-cli/util"
)

// PluginAPI represetns a plugin API
type PluginAPI struct {
	dependencyMap map[string]string

	RoundTripper http.RoundTripper
}

// PluginDependency represents a plugin dependency
type PluginDependency struct {
	Name     string `json:"name"`
	Implied  bool   `json:"implied"`
	Optional bool   `json:"optional"`
	Title    string `json:"title"`
	Version  string `json:"version"`
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

// PluginInfoStats is the plugin info stats
type PluginInfoStats struct {
	CurrentInstalls                   int
	Installations                     []PluginInstallationInfo
	InstallationsPerVersion           []PluginInstallationInfo
	InstallationsPercentage           []PluginInstallationInfo
	InstallationsPercentagePerVersion []PluginInstallationInfo
	Trend                             int
}

// PluginInstallationInfo represents the plugin installation info
type PluginInstallationInfo struct {
	Timestamp  int64
	Total      int
	Version    string
	Percentage float64
}

// ShowTrend show the trend of plugins
func (d *PluginAPI) ShowTrend(name string) (trend string, err error) {
	var plugin *PluginInfo
	if plugin, err = d.getPlugin(name); err != nil {
		return
	}

	data := []float64{}
	installations := plugin.Stats.Installations
	offset, count := 0, 10
	if len(installations) > count {
		offset = len(installations) - count
	}
	for _, installation := range installations[offset:] {
		data = append(data, float64(installation.Total))
	}
	trend = util.PrintCollectTrend(data)
	return
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
	if d.RoundTripper == nil {
		cli.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	} else {
		cli.Transport = d.RoundTripper
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

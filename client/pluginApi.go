package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/util"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// PluginAPI represents a plugin API
type PluginAPI struct {
	dependencyMap map[string]string

	SkipDependency bool
	SkipOptional   bool
	UseMirror      bool
	ShowProgress   bool
	MirrorURL      string

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
	logger.Info("start to collect plugin dependencies...")
	plugins := make([]PluginInfo, 0)
	for _, name := range names {
		logger.Debug("start to collect dependency", zap.String("plugin", name))
		plugins = append(plugins, d.collectDependencies(strings.ToLower(name))...)
	}

	logger.Info("ready to download plugins", zap.Int("total", len(plugins)))
	var err error
	for i, plugin := range plugins {
		logger.Info("start to download plugin",
			zap.String("name", plugin.Name),
			zap.String("version", plugin.Version),
			zap.String("url", plugin.URL),
			zap.Int("number", i))

		if err = d.download(plugin.URL, plugin.Name); err != nil {
			logger.Error("download plugin error", zap.String("name", plugin.Name), zap.Error(err))
		}
	}
}

func (d *PluginAPI) getMirrorURL(url string) (mirror string) {
	mirror = url
	if d.UseMirror && d.MirrorURL != "" {
		logger.Debug("replace with mirror", zap.String("original", url))
		mirror = strings.Replace(url, "http://updates.jenkins-ci.org/download/", d.MirrorURL, -1)
	}
	return
}

func (d *PluginAPI) download(url string, name string) (err error) {
	url = d.getMirrorURL(url)
	logger.Info("prepare to download", zap.String("name", name), zap.String("url", url))

	downloader := util.HTTPDownloader{
		RoundTripper:   d.RoundTripper,
		TargetFilePath: fmt.Sprintf("%s.hpi", name),
		URL:            url,
		ShowProgress:   d.ShowProgress,
	}
	err = downloader.DownloadFile()
	return
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

	pluginAPI := fmt.Sprintf("https://plugins.jenkins.io/api/plugin/%s", name)
	logger.Debug("fetch data from plugin API", zap.String("url", pluginAPI))

	var resp *http.Response
	if resp, err = cli.Get(pluginAPI); err == nil {
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err == nil {
			plugin = &PluginInfo{}
			err = json.Unmarshal(body, plugin)
		}
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
	if d.SkipDependency {
		return
	}

	for _, dependency := range plugin.Dependencies {
		if d.SkipOptional && dependency.Optional {
			continue
		}
		if _, ok := d.dependencyMap[dependency.Name]; !ok {
			d.dependencyMap[dependency.Name] = dependency.Version

			plugins = append(plugins, d.collectDependencies(dependency.Name)...)
		}
	}
	return
}

package cmd

import (
	"fmt"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// ConfigOptions is the config cmd option
type ConfigOptions struct {
	ConfigFileLocation string
}

var configOptions ConfigOptions

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   i18n.T("Manage the config of jcli"),
	Long:    i18n.T("Manage the config of jcli"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		current := getCurrentJenkins()
		if current == nil {
			err = fmt.Errorf("no config file found or no current setting")
		} else {
			if current.Description != "" {
				cmd.Printf("Current Jenkins's name is %s, url is %s, description is %s\n", current.Name, current.URL, current.Description)
			} else {
				cmd.Printf("Current Jenkins's name is %s, url is %s\n", current.Name, current.URL)
			}
		}
		return
	},
	Example: `  jcli config generate
  jcli config list
  jcli config edit`,
}

// JenkinsServer holds the configuration of your Jenkins
type JenkinsServer struct {
	Name               string `yaml:"name"`
	URL                string `yaml:"url"`
	UserName           string `yaml:"username"`
	Token              string `yaml:"token"`
	Proxy              string `yaml:"proxy"`
	ProxyAuth          string `yaml:"proxyAuth"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
	Description        string `yaml:"description"`
}

// CommandHook is a hook
type CommandHook struct {
	Path    string `yaml:"path"`
	Command string `yaml:"cmd"`
}

// PluginSuite define a suite of plugins
type PluginSuite struct {
	Name        string   `yaml:"name"`
	Plugins     []string `yaml:"plugins"`
	Description string   `yaml:"description"`
}

// JenkinsMirror represents the mirror of Jenkins
type JenkinsMirror struct {
	Name string
	URL  string
}

// Config is a global config struct
type Config struct {
	Current        string          `yaml:"current"`
	Language       string          `yaml:"language"`
	JenkinsServers []JenkinsServer `yaml:"jenkins_servers"`
	PreHooks       []CommandHook   `yaml:"preHooks"`
	PostHooks      []CommandHook   `yaml:"postHooks"`
	PluginSuites   []PluginSuite   `yaml:"pluginSuites"`
	Mirrors        []JenkinsMirror `yaml:"mirrors"`
}

func setCurrentJenkins(name string) {
	found := false
	for _, jenkins := range getConfig().JenkinsServers {
		if jenkins.Name == name {
			found = true
			break
		}
	}

	if found {
		config.Current = name
		if err := saveConfig(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("Cannot found Jenkins by name %s", name)
	}
}

var config *Config

func getConfig() *Config {
	return config
}

func getJenkinsNames() []string {
	names := make([]string, 0)
	for _, j := range config.JenkinsServers {
		names = append(names, j.Name)
	}
	return names
}

func getCurrentJenkins() (jenkinsServer *JenkinsServer) {
	if config != nil {
		current := config.Current
		jenkinsServer = findJenkinsByName(current)
	}

	return
}

func findJenkinsByName(name string) (jenkinsServer *JenkinsServer) {
	if config == nil {
		return
	}

	for _, cfg := range config.JenkinsServers {
		if cfg.Name == name {
			jenkinsServer = &cfg
			break
		}
	}
	return
}

func findSuiteByName(name string) (suite *PluginSuite) {
	for _, cfg := range config.PluginSuites {
		if cfg.Name == name {
			suite = &cfg
			break
		}
	}
	return
}

func loadDefaultConfig() (err error) {
	var userHome string
	userHome, err = homedir.Dir()
	if err == nil {
		configPath := fmt.Sprintf("%s/.jenkins-cli.yaml", userHome)
		if _, err = os.Stat(configPath); err == nil {
			err = loadConfig(configPath)
		}
	}
	return
}

func loadConfig(path string) (err error) {
	configOptions.ConfigFileLocation = path

	var content []byte
	if content, err = ioutil.ReadFile(path); err == nil {
		err = yaml.Unmarshal([]byte(content), &config)
	}
	return
}

// getMirrors returns the mirror list, one official mirror should be returned if user don't give it
func getMirrors() (mirrors []JenkinsMirror) {
	if config != nil {
		mirrors = config.Mirrors
	}
	if len(mirrors) == 0 {
		mirrors = []JenkinsMirror{
			{
				Name: "default",
				URL:  "http://mirrors.jenkins.io/",
			},
		}
	}
	return
}

func getMirror(name string) string {
	mirrors := getMirrors()

	for _, mirror := range mirrors {
		if mirror.Name == name {
			logger.Debug("find mirror", zap.String("name", name), zap.String("url", mirror.URL))
			return mirror.URL
		}
	}
	return ""
}

func getDefaultMirror() string {
	return getMirror("default")
}

func saveConfig() (err error) {
	var data []byte
	config := getConfig()

	configPath := configOptions.ConfigFileLocation
	if rootOptions.ConfigFile != "" {
		configPath = rootOptions.ConfigFile
	}

	if data, err = yaml.Marshal(&config); err == nil {
		err = ioutil.WriteFile(configPath, data, 0644)
	}
	return
}

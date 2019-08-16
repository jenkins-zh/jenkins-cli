package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

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
	Short:   "Manage the config of jcli",
	Long:    `Manage the config of jcli`,
	Run: func(_ *cobra.Command, _ []string) {
		current := getCurrentJenkins()
		if current.Description != "" {
			fmt.Printf("Current Jenkins's name is %s, url is %s, description is %s\n", current.Name, current.URL, current.Description)
		} else {
			fmt.Printf("Current Jenkins's name is %s, url is %s\n", current.Name, current.URL)
		}
	},
	Example: `  jcli config generate
  jcli config list
  jcli config edit`,
}

// JenkinsServer holds the configuration of your Jenkins
type JenkinsServer struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	UserName    string `yaml:"username"`
	Token       string `yaml:"token"`
	Proxy       string `yaml:"proxy"`
	ProxyAuth   string `yaml:"proxyAuth"`
	Description string `yaml:"description"`
}

type Config struct {
	Current        string          `yaml:"current"`
	JenkinsServers []JenkinsServer `yaml:"jenkins_servers"`
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
	config := getConfig()
	names := make([]string, 0)
	for _, j := range config.JenkinsServers {
		names = append(names, j.Name)
	}
	return names
}

func getCurrentJenkins() (jenkinsServer *JenkinsServer) {
	config := getConfig()
	current := config.Current
	jenkinsServer = findJenkinsByName(current)

	return
}

func findJenkinsByName(name string) (jenkinsServer *JenkinsServer) {
	for _, cfg := range config.JenkinsServers {
		if cfg.Name == name {
			jenkinsServer = &cfg
			break
		}
	}
	return
}

func loadDefaultConfig() (err error) {
	userHome := userHomeDir()
	configPath := fmt.Sprintf("%s/.jenkins-cli.yaml", userHome)
	if _, err = os.Stat(configPath); err == nil {
		err = loadConfig(configPath)
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

func saveConfig() (err error) {
	var data []byte
	config := getConfig()

	if data, err = yaml.Marshal(&config); err == nil {
		err = ioutil.WriteFile(configOptions.ConfigFileLocation, data, 0644)
	}
	return
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

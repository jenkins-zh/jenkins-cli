package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/linuxsuren/jenkins-cli/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type ConfigOptions struct {
	Current  string
	Show     bool
	Generate bool
	List     bool

	ConfigFileLocation string
}

var configOptions ConfigOptions

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().StringVarP(&configOptions.Current, "current", "c", "", "Set the current Jenkins")
	configCmd.PersistentFlags().BoolVarP(&configOptions.Show, "show", "s", false, "Show the current Jenkins")
	configCmd.PersistentFlags().BoolVarP(&configOptions.Generate, "generate", "g", false, "Generate a sample config file for you")
	configCmd.PersistentFlags().BoolVarP(&configOptions.List, "list", "l", false, "Display all your Jenkins configs")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the config of jcli",
	Long:  `Manage the config of jcli`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("requires at least one argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		current := getCurrentJenkins()
		if configOptions.Show {
			fmt.Printf("Current Jenkins's name is %s, url is %s\n", current.Name, current.URL)
		}

		if configOptions.List {
			table := util.CreateTable(os.Stdout)
			table.AddRow("number", "name", "url")
			for i, jenkins := range getConfig().JenkinsServers {
				name := jenkins.Name
				if name == current.Name {
					name = fmt.Sprintf("*%s", name)
				}
				table.AddRow(fmt.Sprintf("%d", i), name, jenkins.URL)
			}
			table.Render()
		}

		if configOptions.Generate {
			if data, err := generateSampleConfig(); err == nil {
				fmt.Print(string(data))
			} else {
				log.Fatal(err)
			}
		}

		if configOptions.Current != "" {
			setCurrentJenkins(configOptions.Current)
		}
	},
	Example: "jcli config -l",
}

// JenkinsServer holds the configuration of your Jenkins
type JenkinsServer struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url"`
	UserName  string `yaml:"username"`
	Token     string `yaml:"token"`
	Proxy     string `yaml:"proxy"`
	ProxyAuth string `yaml:"proxyAuth"`
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

func generateSampleConfig() ([]byte, error) {
	sampleConfig := Config{
		Current: "yourServer",
		JenkinsServers: []JenkinsServer{
			{
				Name:     "yourServer",
				URL:      "http://localhost:8080/jenkins",
				UserName: "admin",
				Token:    "111e3a2f0231198855dceaff96f20540a9",
			},
		},
	}
	return yaml.Marshal(&sampleConfig)
}

var config Config

func getConfig() Config {
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

func addJenkins(jenkinsServer JenkinsServer) (err error) {
	jenkinsName := jenkinsServer.Name
	if jenkinsName == "" {
		err = fmt.Errorf("Name cannot be empty")
		return
	}

	if findJenkinsByName(jenkinsName) != nil {
		err = fmt.Errorf("Jenkins %s is existed", jenkinsName)
		return
	}

	config.JenkinsServers = append(config.JenkinsServers, jenkinsServer)
	err = saveConfig()
	return
}

func removeJenkins(name string) (err error) {
	current := getCurrentJenkins()
	if name == current.Name {
		err = fmt.Errorf("You cannot remove current Jenkins")
	}

	index := -1
	config := getConfig()
	for i, jenkins := range config.JenkinsServers {
		if name == jenkins.Name {
			index = i
			break
		}
	}

	if index == -1 {
		err = fmt.Errorf("Cannot found by name %s", name)
	} else {
		config.JenkinsServers[index] = config.JenkinsServers[len(config.JenkinsServers)-1]
		config.JenkinsServers[len(config.JenkinsServers)-1] = JenkinsServer{}
		config.JenkinsServers = config.JenkinsServers[:len(config.JenkinsServers)-1]

		err = saveConfig()
	}
	return
}

func loadDefaultConfig() {
	userHome := userHomeDir()
	if err := loadConfig(fmt.Sprintf("%s/.jenkins-cli.yaml", userHome)); err != nil {
		log.Fatalf("error: %v", err)
	}
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

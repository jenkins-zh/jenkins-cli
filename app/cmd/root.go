package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/spf13/cobra"
)

// RootOptions is a global option for whole cli
type RootOptions struct {
	ConfigFile string
	Jenkins    string
	Version    bool
	Debug      bool
}

var rootCmd = &cobra.Command{
	Use:   "jcli",
	Short: "jcli is a tool which could help you with your multiple Jenkins",
	Long: `jcli is Jenkins CLI which could help with your multiple Jenkins,
				  Manage your Jenkins and your pipelines
				  More information could found at https://jenkins-zh.cn`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Jenkins CLI (jcli) manage your Jenkins")

		current := getCurrentJenkinsFromOptionsOrDie()
		if current != nil {
			fmt.Println("Current Jenkins is:", current.Name)
		} else {
			fmt.Println("Cannot found the configuration")
		}

		if rootOptions.Version {
			fmt.Printf("Version: %s\n", app.GetVersion())
			fmt.Printf("Commit: %s\n", app.GetCommit())
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootOptions RootOptions

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&rootOptions.ConfigFile, "configFile", "", "", "An alternative config file")
	rootCmd.PersistentFlags().StringVarP(&rootOptions.Jenkins, "jenkins", "j", "", "Select a Jenkins server for this time")
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Version, "version", "v", false, "Print the version of Jenkins CLI")
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Debug, "debug", "", false, "Print the output into debug.html")
}

func initConfig() {
	if rootOptions.ConfigFile == "" {
		if err := loadDefaultConfig(); err != nil {
			configLoadErrorHandle(err)
		}
	} else {
		if err := loadConfig(rootOptions.ConfigFile); err != nil {
			configLoadErrorHandle(err)
		}
	}
}

func configLoadErrorHandle(err error) {
	if os.IsNotExist(err) {
		log.Printf("No config file found.")
		return
	}

	log.Fatalf("Config file is invalid: %v", err)
}

func getCurrentJenkinsFromOptions() (jenkinsServer *JenkinsServer) {
	jenkinsOpt := rootOptions.Jenkins
	if jenkinsOpt == "" {
		jenkinsServer = getCurrentJenkins()
	} else {
		jenkinsServer = findJenkinsByName(jenkinsOpt)
	}
	return
}

func getCurrentJenkinsFromOptionsOrDie() (jenkinsServer *JenkinsServer) {
	if jenkinsServer = getCurrentJenkinsFromOptions(); jenkinsServer == nil {
		log.Fatal("Cannot found Jenkins by", rootOptions.Jenkins) // TODO not accurate
	}
	return
}

func getCmdPath(cmd *cobra.Command) string {
	current := cmd.Use
	if cmd.HasParent() {
		parentName := getCmdPath(cmd.Parent())
		if parentName == "" {
			return current
		}

		return fmt.Sprintf("%s.%s", parentName, current)
	}
	// don't need the name of root cmd
	return ""
}

func executePreCmd(cmd *cobra.Command, _ []string, writer io.Writer) (err error) {
	config := getConfig()
	if config == nil {
		err = fmt.Errorf("Cannot find config file")
		return
	}

	path := getCmdPath(cmd)
	for _, preHook := range config.PreHooks {
		if path != preHook.Path {
			continue
		}

		if err = execute(preHook.Command, writer); err != nil {
			return
		}
	}
	return
}

func execute(command string, writer io.Writer) (err error) {
	array := strings.Split(command, " ")
	cmd := exec.Command(array[0], array[1:]...)
	cmd.Stdout = writer
	err = cmd.Run()
	return
}

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/spf13/cobra"
)

type RootOptions struct {
	Version bool
	Debug   bool
}

var rootCmd = &cobra.Command{
	Use:   "jcli",
	Short: "jcli is a tool which could help you with your multiple Jenkins",
	Long: `jcli is Jenkins CLI which could help with your multiple Jenkins,
				  Manage your Jenkins and your pipelines
				  More information could found at https://jenkins-zh.cn`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Jenkins CLI (jcli) manage your Jenkins")

		current := getCurrentJenkins()
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
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Version, "version", "v", false, "Print the version of Jenkins CLI")
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.Debug, "debug", "", false, "Print the output into debug.html")
}

func initConfig() {
	if err := loadDefaultConfig(); err != nil {
		if os.IsNotExist(err) {
			log.Printf("No config file found.")
			return
		}

		log.Fatalf("Config file is invalid: %v", err)
	}
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

func executePreCmd(cmd *cobra.Command, _ []string) {
	config := getConfig()
	if config == nil {
		log.Fatal("Cannot find config file")
		return
	}

	path := getCmdPath(cmd)
	for _, preHook := range config.PreHooks {
		if path != preHook.Path {
			continue
		}

		execute(preHook.Command)
	}
}

func execute(command string) {
	array := strings.Split(command, " ")
	cmd := exec.Command(array[0], array[1:]...)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

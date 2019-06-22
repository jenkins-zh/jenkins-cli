package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("hello")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type JenkinsServer struct {
	URL      string `yaml:"url"`
	UserName string `yaml:"username"`
	Token    string `yaml:"token"`
}

type Config struct {
	JenkinsServers []JenkinsServer `yaml:"jenkins_servers"`
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")
}

var config Config

func getConfig() Config {
	return config
}

func initConfig() {
	content, err := ioutil.ReadFile("/Users/mac/.jenkins-cli.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

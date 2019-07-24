package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	configCmd.AddCommand(configGenerateCmd)
}

var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a sample config file for you",
	Long:  `Generate a sample config file for you`,
	Run: func(cmd *cobra.Command, args []string) {
		if data, err := generateSampleConfig(); err == nil {
			fmt.Print(string(data))
			fmt.Println("# Goto 'http://localhost:8080/jenkins/me/configure', then you can generate your token.")
		} else {
			log.Fatal(err)
		}
	},
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

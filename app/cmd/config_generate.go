package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// ConfigGenerateOption is the config generate cmd option
type ConfigGenerateOption struct {
	Copy bool
}

var configGenerateOption ConfigGenerateOption

func init() {
	configCmd.AddCommand(configGenerateCmd)
	configGenerateCmd.Flags().BoolVarP(&configGenerateOption.Copy, "copy", "c", false, "Copy the output into clipboard")
}

var configGenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate a sample config file for you",
	Long:    `Generate a sample config file for you`,
	Run: func(_ *cobra.Command, _ []string) {
		if data, err := generateSampleConfig(); err == nil {
			configPath := configOptions.ConfigFileLocation

			if configPath == "" { // config file isn't exists
				userHome := userHomeDir()
				configPath = fmt.Sprintf("%s/.jenkins-cli.yaml", userHome)
			}

			_, err := os.Stat(configPath)
			if err != nil && os.IsNotExist(err) {
				confirm := false
				prompt := &survey.Confirm{
					Message: "Cannot found your config file, do you want to edit it?",
				}
				survey.AskOne(prompt, &confirm)
				if confirm {
					prompt := &survey.Editor{
						Message:       "Edit your config file",
						FileName:      "*.yaml",
						Default:       string(data),
						HideDefault:   true,
						AppendDefault: true,
					}

					var configContext string
					if err = survey.AskOne(prompt, &configContext); err != nil {
						log.Fatal(err)
					} else {
						if err = ioutil.WriteFile(configPath, []byte(configContext), 0644); err != nil {
							log.Fatal(err)
						}
					}
					return
				}
			}

			printCfg(data)

			if configGenerateOption.Copy {
				clipboard.WriteAll(string(data))
			}
		} else {
			log.Fatal(err)
		}
	},
}

func printCfg(data []byte) {
	fmt.Print(string(data))
	fmt.Println("# Language context is accept-language for HTTP header, It contains zh-CN/zh-TW/en/en-US/ja and so on")
	fmt.Println("# Goto 'http://localhost:8080/jenkins/me/configure', then you can generate your token.")
}

func getSampleConfig() (sampleConfig Config) {
	sampleConfig = Config{
		Current: "yourServer",
		JenkinsServers: []JenkinsServer{
			{
				Name:     "yourServer",
				URL:      "http://localhost:8080/jenkins",
				UserName: "admin",
				Token:    "111e3a2f0231198855dceaff96f20540a9",
			},
		},
		Mirrors: []JenkinsMirror{
			{
				Name: "default",
				URL:  "http://mirrors.jenkins.io/",
			},
			{
				Name: "tsinghua",
				URL:  "https://mirrors.tuna.tsinghua.edu.cn/jenkins/",
			},
			{
				Name: "huawei",
				URL:  "https://mirrors.huaweicloud.com/jenkins/",
			},
			{
				Name: "tencent",
				URL:  "https://mirrors.cloud.tencent.com/jenkins/",
			},
		},
	}
	return
}

func generateSampleConfig() ([]byte, error) {
	sampleConfig := getSampleConfig()
	return yaml.Marshal(&sampleConfig)
}

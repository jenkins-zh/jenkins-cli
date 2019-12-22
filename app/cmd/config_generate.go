package cmd

import (
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/atotto/clipboard"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// ConfigGenerateOption is the config generate cmd option
type ConfigGenerateOption struct {
	InteractiveOption
	Copy bool
}

var configGenerateOption ConfigGenerateOption

func init() {
	configCmd.AddCommand(configGenerateCmd)
	configGenerateCmd.Flags().BoolVarP(&configGenerateOption.Interactive, "interactive", "i", true,
		i18n.T("Interactive mode"))
	configGenerateCmd.Flags().BoolVarP(&configGenerateOption.Copy, "copy", "c", false,
		i18n.T("Copy the output into clipboard"))
}

var configGenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   i18n.T("Generate a sample config file for you"),
	Long:    i18n.T("Generate a sample config file for you"),
	Run: func(cmd *cobra.Command, _ []string) {
		data, err := generateSampleConfig()
		if err == nil {
			if configGenerateOption.Interactive {
				err = InteractiveWithConfig(cmd, data)
			} else {
				printCfg(cmd, data)
			}

			if configGenerateOption.Copy {
				err = clipboard.WriteAll(string(data))
			}
		}
		helper.CheckErr(cmd, err)
	},
}

// InteractiveWithConfig be friendly for a newer
func InteractiveWithConfig(cmd *cobra.Command, data []byte) (err error) {
	configPath := configOptions.ConfigFileLocation

	if configPath == "" { // config file isn't exists
		if configPath, err = GetConfigFromHome(); err != nil {
			return
		}
	}

	_, err = os.Stat(configPath)
	if err != nil && os.IsNotExist(err) {
		confirm := false
		prompt := &survey.Confirm{
			Message: "Cannot found your config file, do you want to edit it?",
		}
		err = survey.AskOne(prompt, &confirm)
		if err == nil && confirm {
			prompt := &survey.Editor{
				Message:       "Edit your config file",
				FileName:      "*.yaml",
				Default:       string(data),
				HideDefault:   true,
				AppendDefault: true,
			}

			var configContext string
			if err = survey.AskOne(prompt, &configContext); err == nil {
				err = ioutil.WriteFile(configPath, []byte(configContext), 0644)
			}
		}
	}
	return
}

func printCfg(cmd *cobra.Command, data []byte) {
	cmd.Print(string(data))
	cmd.Println("# Language context is accept-language for HTTP header, It contains zh-CN/zh-TW/en/en-US/ja and so on")
	cmd.Println("# Goto 'http://localhost:8080/jenkins/me/configure', then you can generate your token.")
}

func getSampleConfig() (sampleConfig Config) {
	sampleConfig = Config{
		Current: "yourServer",
		JenkinsServers: []JenkinsServer{
			{
				Name:               "yourServer",
				URL:                "http://localhost:8080/jenkins",
				UserName:           "admin",
				Token:              "111e3a2f0231198855dceaff96f20540a9",
				InsecureSkipVerify: true,
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

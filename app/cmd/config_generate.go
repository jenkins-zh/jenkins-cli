package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/atotto/clipboard"
	appCfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/spf13/cobra"
)

// ConfigGenerateOption is the config generate cmd option
type ConfigGenerateOption struct {
	common.InteractiveOption
	common.Option
	common.BatchOption

	Copy bool
}

var configGenerateOption ConfigGenerateOption

func init() {
	configCmd.AddCommand(configGenerateCmd)
	configGenerateCmd.Flags().BoolVarP(&configGenerateOption.Interactive, "interactive", "i", true,
		i18n.T("Interactive mode"))
	configGenerateCmd.Flags().BoolVarP(&configGenerateOption.Copy, "copy", "c", false,
		i18n.T("Copy the output into clipboard"))
	configGenerateOption.Option.Stdio = common.GetSystemStdio()
	configGenerateOption.BatchOption.Stdio = common.GetSystemStdio()
}

var configGenerateCmd = &cobra.Command{
	Use:               "generate",
	Aliases:           []string{"gen"},
	Short:             i18n.T("Generate a sample config file for you"),
	Long:              i18n.T("Generate a sample config file for you"),
	ValidArgsFunction: common.NoFileCompletion,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		var data []byte
		data, err = GenerateSampleConfig()
		if err == nil {
			if configGenerateOption.Interactive {
				err = configGenerateOption.InteractiveWithConfig(cmd, data)
			} else {
				printCfg(cmd, data)
			}

			if configGenerateOption.Copy {
				err = clipboard.WriteAll(string(data))
			}
		}
		return
	},
}

// InteractiveWithConfig be friendly for a newer
func (o *ConfigGenerateOption) InteractiveWithConfig(cmd *cobra.Command, data []byte) (err error) {
	configPath := configOptions.ConfigFileLocation
	if configPath == "" {
		configPath, err = getDefaultConfigPath()
	}

	if err == nil {
		_, err = os.Stat(configPath)
	}

	if err != nil && os.IsNotExist(err) {
		confirm := o.Confirm("Cannot found your config file, do you want to edit it?")
		if confirm {
			var content string
			content, err = o.Editor(string(data), "Edit your config file")
			if err == nil {
				logger.Debug("write generated config file", zap.String("path", configPath))
				err = ioutil.WriteFile(configPath, []byte(content), 0644)
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

func getSampleConfig() (sampleConfig appCfg.Config) {
	sampleConfig = appCfg.Config{
		Current: "yourServer",
		JenkinsServers: []appCfg.JenkinsServer{
			{
				Name:               "yourServer",
				URL:                "http://localhost:8080/jenkins",
				UserName:           "admin",
				Token:              "111e3a2f0231198855dceaff96f20540a9",
				InsecureSkipVerify: true,
			},
		},
		Mirrors: []appCfg.JenkinsMirror{
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

// GenerateSampleConfig returns a sample config
func GenerateSampleConfig() ([]byte, error) {
	sampleConfig := getSampleConfig()
	return yaml.Marshal(&sampleConfig)
}

// GetConfigFromHome returns the config file path from user home dir
func GetConfigFromHome() (configPath string, homeErr error) {
	userHome, homeErr := homedir.Dir()
	if homeErr == nil {
		configPath = fmt.Sprintf("%s/.jenkins-cli.yaml", userHome)
	}
	return
}

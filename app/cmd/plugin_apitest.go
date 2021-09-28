package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	jenkinsFormula "github.com/jenkins-zh/jenkins-formulas/pkg/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	pluginCmd.AddCommand(pluginAPITestCmd)
	flags := pluginAPITestCmd.Flags()
	flags.StringVarP(&pluginAPITestO.ip, "ip", "", "127.0.0.1",
		i18n.T("The ip address of the computer you want to use"))
	flags.StringVarP(&pluginAPITestO.port, "port", "", "8080",
		i18n.T("The port to connect to the jenkins you want to test"))
	flags.StringVar(&pluginAPITestO.yamlFile, "yaml", "",
		i18n.T("The yaml file you used to create the jenkins war with command 'jcli cwp' and this yaml is needed to decide which plugins is going to be tested.\n"+
			"And If you don't have the yaml file, you can use the 'jcli create yaml' to create one. You only need to provide either --yaml or --custom-yaml"))
	flags.StringVar(&pluginAPITestO.testYaml, "custom-yaml", "",
		i18n.T("The test yaml file is needed only you choose to conduct an API test for plugins in a custom mode.\n"+
			"And if you just want to conduct a simple API test, you don't need to specify the test yaml file. You only need to provide either --yaml or --custom-yaml"))
}

type pluginAPITestOption struct {
	ip                 string
	port               string
	jenkinsPluginTest  jenkinsFormula.CustomWarPackage
	pluginsWithProblem []string
	//the yaml file used to create the jenkins war
	yamlFile string
	//if choose to test API in a custom mode this testYaml file must be pointed out
	testYaml string
}

type apiTestOption struct {
	Plugins []plugin `yaml:"plugins"`
}

type plugin struct {
	ArtifactId string   `yaml:"artifactId"`
	API        []string `yaml:"api"`
}

var pluginAPITestO pluginAPITestOption
var apiTestO apiTestOption

var pluginAPITestCmd = &cobra.Command{
	Use:   "api test",
	Short: "Conduct an API test for plugins of jenkins started in a docker container with setupWizard=false",
	Long: "Conduct an API test for plugins of jenkins started ina docker container with setupWizard=false. The API test is provided in two modes：simple and custom. " +
		"Choose the simple mode, a yaml file created by 'jcli create yaml' is needed. Choose the custom mode, a yaml file contains plugins artifactID and api list is needed.",
	Example: "jcli plugin api test",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if pluginAPITestO.yamlFile != "" {
			if exist, _ := CheckFileExists(pluginAPITestO.yamlFile); !exist {
				prompt := fmt.Sprintf("The %s doesn't exist.", pluginAPITestO.yamlFile)
				panic(prompt)
			}
		} else if pluginAPITestO.testYaml != "" {
			if exist, _ := CheckFileExists(pluginAPITestO.testYaml); !exist {
				prompt := fmt.Sprintf("The %s doesn't exist.", pluginAPITestO.testYaml)
				panic(prompt)
			}
		}
		jClient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginFormulaOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))
		jClient.JenkinsCore.URL = fmt.Sprintf("http://%s:%s", pluginAPITestO.ip, pluginAPITestO.port)
		if pluginAPITestO.testYaml != "" {
			if file, err := ioutil.ReadFile(pluginAPITestO.testYaml); err == nil {
				err := yaml.Unmarshal(file, &apiTestO)
				if err != nil {
					return err
				}
				pluginsWithProblemMap := make(map[string]string)
				for _, plugin := range apiTestO.Plugins {
					apis := plugin.API
					for _, api := range apis {
						statusCode, _, err := jClient.JenkinsCore.Request(http.MethodGet, api, nil, nil)
						if err != nil {
							cmd.Println(err)
							return err
						}
						if statusCode != 200 {
							errorAPI, ok := pluginsWithProblemMap[plugin.ArtifactId]
							if ok {
								pluginsWithProblemMap[plugin.ArtifactId] = errorAPI + "   " + api
							} else {
								pluginsWithProblemMap[plugin.ArtifactId] = api
							}
						}
					}
				}
				if len(pluginsWithProblemMap) != 0 {
					cmd.Print("There's something wrong with the plugin(s):\n")
					for pluginName, url := range pluginsWithProblemMap {
						cmd.Println(fmt.Sprintf("%-18s: %s ", pluginName, url))
					}
				} else if len(pluginsWithProblemMap) == 0 {
					cmd.Println("Congratulations! All your plugins work fine.")
				}
			}

		} else if pluginAPITestO.yamlFile != "" {
			if file, err := ioutil.ReadFile(pluginAPITestO.yamlFile); err == nil {
				err := yaml.Unmarshal(file, &pluginAPITestO.jenkinsPluginTest)
				if err != nil {
					return err
				}
				pluginAPITestO.pluginsWithProblem = make([]string, 0)
				var pluginBuffer bytes.Buffer
				var i = 0
				for _, plugin := range pluginAPITestO.jenkinsPluginTest.Plugins {
					api := fmt.Sprintf("/pluginManager/plugin/%s/api/json", plugin.ArtifactId)
					statusCode, _, err := jClient.JenkinsCore.Request(http.MethodGet, api, nil, nil)
					if err != nil {
						cmd.Println(err)
						return err
					}
					if statusCode != 200 {
						pluginBuffer.WriteString(plugin.ArtifactId + "\n")
						i++
					}
				}
				pluginString := pluginBuffer.String()
				pluginAPITestO.pluginsWithProblem = strings.Split(pluginString, "\n")
				if len(pluginAPITestO.pluginsWithProblem) != 1 {
					cmd.Print("There's something wrong with the plugin(s):\n")
					for index, plugin := range pluginAPITestO.pluginsWithProblem {
						if index%5 == 0 {
							cmd.Println()
						}
						cmd.Print(fmt.Sprintf("%-25s", plugin))
					}
				} else if len(pluginAPITestO.pluginsWithProblem) == 1 {
					cmd.Println("Congratulations! All your plugins work fine.")
				}
			}
		}
		return err
	},
}

//CheckFileExists returns true if exits and returns false if not
func CheckFileExists(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

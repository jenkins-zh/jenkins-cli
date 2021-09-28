package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/pkg/docker"
	jenkinsFormula "github.com/jenkins-zh/jenkins-formulas/pkg/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var formulaYaml jenkinsFormula.CustomWarPackage
var all bool

func init() {
	pluginFormulaOption.DockerBuild = false
	pluginFormulaOption.OnlyRelease = true
	rootCmd.AddCommand(upgradeLabmouseCmd)
	upgradeLabmouseCmd.Flags().BoolVarP(&all, "all", "", false, i18n.T("Upgrade jenkins core and all plugins to update"))
	upgradeLabmouseCmd.Flags().StringVarP(&formulaYaml.Bundle.GroupId, "bundle-groupId", "", "io.jenkins.tools.jcli.yaml.demo", i18n.T("GourpId of Bundle in yaml"))
	upgradeLabmouseCmd.Flags().StringVarP(&formulaYaml.Bundle.ArtifactId, "bundle-artifactId", "", "jcli-yaml-demo", i18n.T("ArtifactId of Bundle in yaml"))
	upgradeLabmouseCmd.Flags().StringVarP(&formulaYaml.Bundle.Vendor, "bundle-vendor", "", "jenkins-cli", i18n.T("Vendor of Bundle in yaml"))
	upgradeLabmouseCmd.Flags().StringVarP(&formulaYaml.Bundle.Description, "bundle-description", "", "Upgraded jenkins core and plugins in a YAML specification", i18n.T("Description of Bundle in yaml"))
	healthCheckRegister.Register(getCmdPath(upgradeLabmouseCmd), &pluginFormulaOption)
	upgradeLabmouseCmd.Flags().StringVarP(&docker.DockerRunOption.IP, "ip", "", "127.0.0.1",
		i18n.T("The ip address of the computer you want to use"))
	upgradeLabmouseCmd.Flags().IntVarP(&docker.DockerRunOption.DockerPort, "docker-port", "", 2375,
		i18n.T("The port to connect to docker"))
	upgradeLabmouseCmd.Flags().IntVarP(&docker.DockerRunOption.JenkinsPort, "Jenkins-port", "", 8081,
		i18n.T("The port to connect to jenkins"))
	upgradeLabmouseCmd.Flags().StringVar(&pluginAPITestO.testYaml, "custom-yaml", "",
		i18n.T("The test yaml file is needed only you choose to conduct an API test for plugins in a custom mode.\n"+
			"And if you just want to conduct a simple API test, you don't need to specify the test yaml file. You only need to provide either --yaml or --custom-yaml"))
}

var upgradeLabmouseCmd = &cobra.Command{
	Use:   "upgrade labmouse",
	Short: i18n.T("This function is to test the viability of jenkins after upgrading jenkins core and(or) some plugins."),
	Long: i18n.T(`This function is to test the viability of jenkins after upgrading jenkins core and(or) some plugins. It will start a jenkins which contains the jenkins core and the plugins you want to upgrade, in a docker container.
	Then conduct a simple or custom API test and output which plugins fail the tests. Finally if the jenkins after upgrading works fine, you can choose to upgrade it and its plugins.`),
	Example: `upgrade labmouse --all
upgrade labmouse
upgrade labmouse --custom-yaml <yamlfile>`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = multipleChoice(cmd, args)
		if err != nil {
			return
		}
		err = cwpOptions.RunWithoutProcessExits(cmd, args)
		if err != nil {
			return
		}
		err = docker.DockerRunOption.CreateDockerfile(cmd, args)
		if err != nil {
			cmd.Println(err)
			return
		}
		docker.DockerRunOption.ImageName = "jclitest"
		docker.DockerRunOption.Tag = "latest"
		err = docker.DockerRunOption.CreateImageAndRunContainer(cmd, args)
		if err != nil {
			cmd.Println(err)
			return
		}
		pluginAPITestO.ip = docker.DockerRunOption.IP
		pluginAPITestO.port = strconv.Itoa(docker.DockerRunOption.JenkinsPort)
		ready, err := waitForJenkinsToBeReady(cmd)
		if err != nil && !ready {
			return fmt.Errorf("oops, jenkins didn't start successfully or needed more time to be ready")
		}
		err = pluginAPITestO.test(cmd, args)
		if err != nil {
			return
		}
		return
	},
}

// mutipleChoice prompt for users to choose which to upgrade
func multipleChoice(cmd *cobra.Command, args []string) (err error) {
	targetPlugins := make([]string, 0)
	formulaYaml.War.GroupId = "org.jenkins-ci.main"
	formulaYaml.War.ArtifactId = "jenkins-war"
	if all {
		if _, err := getLocalJenkinsAndPlugins(); err == nil {
			if pluginFormulaOption.OnlyRelease {
				formulaYaml.Plugins = removeSnapshotPluginsAndUpgradeOrNot(formulaYaml.Plugins, true)
			}
			if pluginFormulaOption.SortPlugins {
				formulaYaml.Plugins = SortPlugins(formulaYaml.Plugins)
			}
			if items, _, err := getVersionData(LtsURL); err == nil {
				formulaYaml.War.Source.Version = "\"" + items[0].Title[8:] + "\""
			}
			formulaYaml.BuildSettings.Docker = jenkinsFormula.BuildDockerSetting{
				Base:  fmt.Sprintf("jenkins/jenkins:%s", formulaYaml.War.Source.Version),
				Tag:   "jenkins/jenkins-formula:v0.0.1",
				Build: pluginFormulaOption.DockerBuild,
			}
		}
	} else if !all {
		var coreTemp bool
		if jenkinsVersion, err := getLocalJenkinsAndPlugins(); err == nil {
			promptCore := &survey.Confirm{
				Message: fmt.Sprintf("Please indicate whether do you want to upgrade or not"),
			}
			err = survey.AskOne(promptCore, &coreTemp)
			if err != nil {
				return err
			}
			if coreTemp {
				if items, _, err := getVersionData(LtsURL); err == nil {
					formulaYaml.War.Source.Version = "\"" + items[0].Title[8:] + "\""
				}
			} else if !coreTemp {
				formulaYaml.War.Source.Version = jenkinsVersion
			}
			prompt := &survey.MultiSelect{
				Message: fmt.Sprintf("Please select the plugins(%d) which you want to upgrade to the latest: ", len(formulaYaml.Plugins)),
				Options: ConvertPluginsToArray(formulaYaml.Plugins),
			}
			err = survey.AskOne(prompt, &targetPlugins)

			if err != nil {
				return err
			}
			tempMap := make(map[string]bool)
			for _, plugin := range targetPlugins {
				tempMap[plugin] = true
			}
			for index, plugin := range formulaYaml.Plugins {
				if _, exist := tempMap[plugin.ArtifactId]; exist {
					formulaYaml.Plugins[index].Source.Version, err = getNewVersionOfPlugin(plugin.ArtifactId)
				}
				if err != nil {
					return err
				}
			}
			formulaYaml.BuildSettings.Docker = jenkinsFormula.BuildDockerSetting{
				Base:  fmt.Sprintf("jenkins/jenkins:%s", formulaYaml.War.Source.Version),
				Tag:   "jenkins/jenkins-formula:v0.0.1",
				Build: pluginFormulaOption.DockerBuild,
			}
		}
	}
	renderYaml(formulaYaml)
	return nil
}

func getNewVersionOfPlugin(pluginName string) (version string, err error) {
	api := "https://plugins.jenkins.io/api/plugin/" + pluginName
	resp, err := http.Get(api)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	var newPluginOption NewPluginOption
	err = json.Unmarshal(bytes, &newPluginOption)
	if err != nil {
		return "", err
	}
	version = trimToID(newPluginOption.Version)

	return version, nil
}

func trimToID(content string) (version string) {
	startOfVersionNumber := strings.LastIndex(content, ":")
	version = content[startOfVersionNumber+1:]
	return version
}

func renderYaml(yamlTemp jenkinsFormula.CustomWarPackage) (err error) {
	data, err := yaml.Marshal(&yamlTemp)
	if err != nil {
		return err
	}
	dir, err := ioutil.TempDir("", "jenkins-cli")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dir+"/test.yaml", data, 0)
	cwpOptions.ConfigPath = dir + "/test.yaml"
	cwpOptions.TmpDir = dir + "/cwp"
	pluginAPITestO.yamlFile = dir + "/test.yaml"
	docker.DockerRunOption.WarPath = fmt.Sprintf("%s/cwp/output/target/jcli-yaml-demo-1.0-SNAPSHOT.war", dir)
	if err != nil {
		return err
	}
	return nil
}

func removeSnapshotPluginsAndUpgradeOrNot(plugins []jenkinsFormula.Plugin, allTmp bool) (result []jenkinsFormula.Plugin) {
	result = make([]jenkinsFormula.Plugin, 0)
	if allTmp {
		for i := range plugins {
			if strings.Contains(plugins[i].Source.Version, "SNAPSHOT") {
				continue
			}
			plugins[i].Source.Version, _ = getNewVersionOfPlugin(plugins[i].ArtifactId)
			result = append(result, plugins[i])
		}
	} else if !allTmp {
		for i := range plugins {
			if strings.Contains(plugins[i].Source.Version, "SNAPSHOT") {
				continue
			}
			result = append(result, plugins[i])
		}
	}
	return result
}
func getLocalJenkinsAndPlugins() (jenkinsVersion string, err error) {
	jClientPlugin := &client.PluginManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: pluginListOption.RoundTripper,
		},
	}
	GetCurrentJenkinsAndClient(&(jClientPlugin.JenkinsCore))
	if err = jClientPlugin.GetPluginsFormula(&(formulaYaml.Plugins)); err != nil {
		return "", err
	}

	jClientCore := &client.JenkinsStatusClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: centerOption.RoundTripper,
		},
	}
	GetCurrentJenkinsAndClient(&(jClientCore.JenkinsCore))
	status, error := jClientCore.Get()
	if error != nil {
		return "", error
	}
	jenkinsVersion = status.Version
	return jenkinsVersion, nil
}

// ConvertPluginsToArray convert jenkinsFormula.Plugin to slice for the sake of multiple select
func ConvertPluginsToArray(plugins []jenkinsFormula.Plugin) (pluginArray []string) {
	pluginArray = make([]string, 0)
	for _, plugin := range plugins {
		pluginArray = append(pluginArray, plugin.ArtifactId)
	}
	return pluginArray
}

func waitForJenkinsToBeReady(cmd *cobra.Command) (ready bool, err error) {
	var statusCode int
	if resp, _ := http.Get("http://" + pluginAPITestO.ip + ":" + pluginAPITestO.port); resp == nil {
		statusCode = 404
	} else {
		statusCode = resp.StatusCode
	}
	count := 0
	// only wait for 60*2 sec which equals to 2 mins for jenkins to be ready
	for count <= 60 && statusCode != 200 {
		time.Sleep(2 * time.Second)
		count++
		if resp, _ := http.Get("http://" + pluginAPITestO.ip + ":" + pluginAPITestO.port); resp == nil {
			statusCode = 404
		} else {
			statusCode = resp.StatusCode
		}
	}
	if count >= 61 {
		return false, err
	}
	return true, err
}

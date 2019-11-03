package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// DoctorOption is the doctor cmd option
type DoctorOption struct {
	OutputOption

	RoundTripper http.RoundTripper
}

var doctorOption DoctorOption

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check your Jenkins config list status and current Jenkins config plugins",
	Long:  `Check your Jenkins config list status and current Jenkins config plugins`,
	Run: func(_ *cobra.Command, _ []string) {
		jenkinsNames := getJenkinsNames()
		checkDuplicateName(jenkinsNames)
		jenkinsServers := getConfig().JenkinsServers
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: doctorOption.RoundTripper,
			},
		}
		checkJenkinsServersStatus(jenkinsServers, jclient)
		checkCurrentPlugins(jclient)
		fmt.Println("Checked is done.")
	},
}

func checkDuplicateName(jenkinsNames []string) {
	fmt.Println("Begining checking the name in the configuration file is duplicated：")
	var duplicateName = ""
	for i := range jenkinsNames {
		for j := range jenkinsNames {
			if i != j && jenkinsNames[i] == jenkinsNames[j] && !strings.Contains(duplicateName, jenkinsNames[i]) {
				duplicateName += jenkinsNames[i] + " "
			}
		}
	}
	if duplicateName == "" {
		fmt.Println("  Checked it sure. no duplicated config Name")
	} else {
		fmt.Printf("  Duplicate names: %s\n", duplicateName)
	}
}

func checkJenkinsServersStatus(jenkinsServers []JenkinsServer, jclient *client.PluginManager) {
	fmt.Println("Begining checking JenkinsServer status form the configuration files: ")
	for i := range jenkinsServers {
		jenkinsServer := jenkinsServers[i]
		jclient.URL = jenkinsServer.URL
		jclient.UserName = jenkinsServer.UserName
		jclient.Token = jenkinsServer.Token
		jclient.Proxy = jenkinsServer.Proxy
		jclient.ProxyAuth = jenkinsServer.ProxyAuth
		fmt.Printf("  checking the No.%d - %s status: ", i, jenkinsServer.Name)
		if _, err := jclient.GetPlugins(); err == nil {
			fmt.Println("***available***")
		} else {
			fmt.Println("***unavailable***", err)
		}
	}
}

func checkCurrentPlugins(jclient *client.PluginManager) {
	fmt.Println("Begining checking the current JenkinsServer's plugins status: ")
	getCurrentJenkinsAndClient(&jclient.JenkinsCore)
	if plugins, err := jclient.GetPlugins(2); err == nil {
		cyclePlugins(plugins)
	} else {
		fmt.Println("  No plugins have updated...")
	}
}

func cyclePlugins(plugins *client.InstalledPluginList) {
	for _, plugin := range plugins.Plugins {
		fmt.Printf("  Checking the plugin %s: \n", plugin.ShortName)
		dependencies := plugin.Dependencies
		if len(dependencies) != 0 {
			cycleDependencies(dependencies, plugins)
		} else {
			fmt.Println("    The Plugin no dependencies")
		}
	}
}

func cycleDependencies(dependencies []client.Dependence, plugins *client.InstalledPluginList) {
	for _, dependence := range dependencies {
		fmt.Printf("    Checking the dependence plugin %s: ", dependence.ShortName)
		hasInstalled := false
		needUpdate := false
		cycleMatchPlugins(plugins, dependence, hasInstalled, needUpdate)
	}
}

func cycleMatchPlugins(plugins *client.InstalledPluginList, dependence client.Dependence, hasInstalled bool, needUpdate bool) {
	for _, checkPlugin := range plugins.Plugins {
		checkPluginVersion := strings.Split(checkPlugin.Version, ".")
		dependenceVersion := strings.Split(dependence.Version, ".")
		if checkPlugin.ShortName == dependence.ShortName {
			hasInstalled = true
			matchPlugin(dependenceVersion, checkPluginVersion, needUpdate, dependence)
		}
		if needUpdate {
			break
		}
	}
	if !hasInstalled {
		fmt.Printf("    The dependence %s no install, please install it the version %s at least\n", dependence.ShortName, dependence.Version)
	}
}

func matchPlugin(dependenceVersion []string, checkPluginVersion []string, needUpdate bool, dependence client.Dependence) (isPass bool) {
	for i := range dependenceVersion {
		if strings.Contains(dependenceVersion[i], "-") && strings.Contains(checkPluginVersion[i], "-") {
			isPass = matchPlugin(strings.Split(dependenceVersion[i], "-"), strings.Split(checkPluginVersion[i], "-"), needUpdate, dependence)
			if isPass {
				break
			}
		} else if len(checkPluginVersion) >= i+1 && len(dependenceVersion) >= i+1 {
			checkPluginVersionInt, _ := strconv.Atoi(checkPluginVersion[i])
			dependenceVersionInt, _ := strconv.Atoi(dependenceVersion[i])
			if checkPluginVersionInt == dependenceVersionInt {
				if i+1 == len(dependenceVersion) {
					isPass = true
					fmt.Println("***true***")
					break
				} else {
					continue
				}
			} else if checkPluginVersionInt > dependenceVersionInt {
				isPass = true
				fmt.Println("***true***")
				break
			} else {
				isPass = true
				needUpdate = true
				fmt.Printf("The dependence %s need upgrade the version to %s\n", dependence.ShortName, dependence.Version)
				break
			}
		}
	}
	return
}

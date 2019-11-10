package cmd

import (
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
	Run: func(cmd *cobra.Command, _ []string) {
		jenkinsNames := getJenkinsNames()
		outString := ""
		if outputDuplicateNames, err := checkNameDuplicate(jenkinsNames); err == nil {
			outString += outputDuplicateNames
		}
		jenkinsServers := getConfig().JenkinsServers
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: doctorOption.RoundTripper,
			},
		}
		if outputJenkinsStatus, err := checkJenkinsServersStatus(jenkinsServers, jclient); err == nil {
			outString += outputJenkinsStatus
		}
		if outputCurrentPluginStatus, err := checkCurrentPlugins(jclient); err == nil {
			outString += outputCurrentPluginStatus
		}
		outString += "Checked is done.\n"
		cmd.Print(outString)
	},
}

func checkNameDuplicate(jenkinsNames []string) (outString string, err error) {
	outString = "Beginning to check the names are duplicated which in the configuration file：\n"
	var duplicateNames = ""
	for i := range jenkinsNames {
		for j := range jenkinsNames {
			if i != j && jenkinsNames[i] == jenkinsNames[j] && !strings.Contains(duplicateNames, jenkinsNames[i]) {
				duplicateNames += jenkinsNames[i] + " "
			}
		}
	}
	if duplicateNames != "" {
		outString += "  Duplicate names: " + duplicateNames + "\n"
	}
	return
}

func checkJenkinsServersStatus(jenkinsServers []JenkinsServer, jclient *client.PluginManager) (outString string, err error) {
	outString = "Begining to checking JenkinsServer status form the configuration files: \n"
	for i := range jenkinsServers {
		jenkinsServer := jenkinsServers[i]
		jclient.URL = jenkinsServer.URL
		jclient.UserName = jenkinsServer.UserName
		jclient.Token = jenkinsServer.Token
		jclient.Proxy = jenkinsServer.Proxy
		jclient.ProxyAuth = jenkinsServer.ProxyAuth
		outString += "  checking the No." + strconv.Itoa(i) + " - " + jenkinsServer.Name + " status: "
		if _, err := jclient.GetPlugins(); err == nil {
			outString += "***available***\n"
		} else {
			outString += "***unavailable*** " + err.Error() + "\n"
		}
	}
	return
}

func checkCurrentPlugins(jclient *client.PluginManager) (outString string, err error) {
	outString = "Begining to checking the current JenkinsServer's plugins status: \n"
	getCurrentJenkinsAndClient(&jclient.JenkinsCore)
	if plugins, err := jclient.GetPlugins(2); err == nil {
		if outputInstalledPluginsStatus, err := checkAllInstalledPlugins(plugins); err == nil {
			outString += outputInstalledPluginsStatus
		}
	} else {
		outString += "  No plugins have lost dependencies...\n"
	}
	return
}

// checkAllInstalledPlugins is check all installed plugins
func checkAllInstalledPlugins(plugins *client.InstalledPluginList) (outString string, err error) {
	outString = ""
	for _, plugin := range plugins.Plugins {
		outString += "  Checking the plugin " + plugin.ShortName + ": \n"
		dependencies := plugin.Dependencies
		if len(dependencies) != 0 {
			if outputDependenciesStatus, err := cycleDependencies(dependencies, plugins); err == nil {
				outString += outputDependenciesStatus
			}
		} else {
			outString += "    The Plugin no dependencies\n"
		}
	}
	return
}

func cycleDependencies(dependencies []client.PluginDependency, plugins *client.InstalledPluginList) (outString string, err error) {
	outString = ""
	for _, dependence := range dependencies {
		outString += "    Checking the dependence plugin " + dependence.ShortName + ": "
		hasInstalled := false
		needUpdate := false
		if outputMatchPlugins, err := cycleMatchPlugins(plugins, dependence, hasInstalled, needUpdate); err == nil {
			outString += outputMatchPlugins
		}
	}
	return
}

func cycleMatchPlugins(plugins *client.InstalledPluginList, dependence client.PluginDependency, hasInstalled bool, needUpdate bool) (outString string, err error) {
	outString = ""
	for _, checkPlugin := range plugins.Plugins {
		checkPluginVersion := strings.Split(checkPlugin.Version, ".")
		dependenceVersion := strings.Split(dependence.Version, ".")
		if checkPlugin.ShortName == dependence.ShortName {
			hasInstalled = true
			if outputMatchPluginStatus, _, err := matchPlugin(dependenceVersion, checkPluginVersion, dependence); err == nil {
				outString += outputMatchPluginStatus
			}

		}
		if needUpdate {
			break
		}
	}
	if !hasInstalled {
		outString += "\n    The dependence " + dependence.ShortName + " no install, please install it the version " + dependence.Version + " at least\n"
	}
	return
}

func matchPlugin(dependenceVersion []string, checkPluginVersion []string, dependence client.PluginDependency) (outString string, isPass bool, err error) {
	outString = ""
	for i := range dependenceVersion {
		if strings.Contains(dependenceVersion[i], "-") && strings.Contains(checkPluginVersion[i], "-") {
			dependenciesVersion := strings.Split(dependenceVersion[i], "-")
			checkPluginsVersion := strings.Split(checkPluginVersion[i], "-")
			if outputCycleMatchSplitValues, hasPass, err := cycleMatchSplitValues(dependenciesVersion, checkPluginsVersion, dependence); err == nil {
				outString += outputCycleMatchSplitValues
				isPass = hasPass
			}
			//_, isPass, _ = matchPlugin(strings.Split(dependenceVersion[i], "-"), strings.Split(checkPluginVersion[i], "-"), dependence)
		} else if len(checkPluginVersion) >= i+1 && len(dependenceVersion) >= i+1 {
			if outputJudgmentValue, hasPass, err := judgmentvalue(i, dependenceVersion, checkPluginVersion, dependence); err == nil {
				outString += outputJudgmentValue
				isPass = hasPass
			}
		}
		if isPass {
			break
		}
	}
	return
}

func judgmentvalue(i int, dependenceVersion []string, checkPluginVersion []string, dependence client.PluginDependency) (outString string, isPass bool, err error) {
	defaultValue := "***true***\n"
	checkPluginVersionInt, _ := strconv.Atoi(checkPluginVersion[i])
	dependenceVersionInt, _ := strconv.Atoi(dependenceVersion[i])
	if checkPluginVersionInt == dependenceVersionInt {
		if i+1 == len(dependenceVersion) {
			isPass = true
			outString += defaultValue
		}
	} else if checkPluginVersionInt > dependenceVersionInt {
		isPass = true
		outString += defaultValue
	} else {
		isPass = true
		outString += "\n      The dependence " + dependence.ShortName + " need upgrade the version to " + dependence.Version + "\n"
	}
	return
}

func cycleMatchSplitValues(dependenciesVersion []string, checkPluginsVersion []string, dependence client.PluginDependency) (outString string, isPass bool, err error) {
	for i := range checkPluginsVersion {
		checkPluginVersion := strings.Split(checkPluginsVersion[i], ".")
		dependenceVersion := strings.Split(dependenciesVersion[i], ".")
		_, isPass, _ = matchPlugin(strings.Split(dependenceVersion[i], "-"), strings.Split(checkPluginVersion[i], "-"), dependence)
		if isPass {
			outString += "***true***\n"
			break
		} else if len(checkPluginsVersion)-1 == i {
			outString += "\n      The dependence " + dependence.ShortName + " need upgrade the version to " + dependence.Version + "\n"
		}
	}
	return
}

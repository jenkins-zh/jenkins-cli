package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/util"
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
		err := checkDuplicateName(jenkinsNames)
		outString += err.Error()
		jenkinsServers := getConfig().JenkinsServers
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: doctorOption.RoundTripper,
			},
		}
		err = checkJenkinsServersStatus(jenkinsServers, jclient)
		outString += err.Error()
		err = checkCurrentPlugins(jclient)
		outString += err.Error()
		outString += "Checked is done."
		cmd.Print(outString)
	},
}

func checkDuplicateName(jenkinsNames []string) (err error){
	outString := "Begining checking the name in the configuration file is duplicated：\n"
	var duplicateName = ""
	for i := range jenkinsNames {
		for j := range jenkinsNames {
			if i != j && jenkinsNames[i] == jenkinsNames[j] && !strings.Contains(duplicateName, jenkinsNames[i]) {
				duplicateName += jenkinsNames[i] + " "
			}
		}
	}
	if duplicateName == "" {
		outString += "  Checked it sure. no duplicated config Name\n"
	} else {
		outString += "  Duplicate names: " + duplicateName + "\n"
	}
	err = errors.New(outString)
	return
}

func checkJenkinsServersStatus(jenkinsServers []JenkinsServer, jclient *client.PluginManager) (err error){
	outString := "Begining checking JenkinsServer status form the configuration files: \n"
	for i := range jenkinsServers {
		jenkinsServer := jenkinsServers[i]
		jclient.URL = jenkinsServer.URL
		jclient.UserName = jenkinsServer.UserName
		jclient.Token = jenkinsServer.Token
		jclient.Proxy = jenkinsServer.Proxy
		jclient.ProxyAuth = jenkinsServer.ProxyAuth
		outString += "  checking the No." + strconv.Itoa(i) + "- " + jenkinsServer.Name + " status: "
		if _, err := jclient.GetPlugins(); err == nil {
			outString += "***available***\n"
		} else {
			outString += "***unavailable***" + err.Error() + "\n"
		}
	}
	err = errors.New(outString)
	return
}

func checkCurrentPlugins(jclient *client.PluginManager) (err error){
	outString := "Begining checking the current JenkinsServer's plugins status: \n"
	getCurrentJenkinsAndClient(&jclient.JenkinsCore)
	if plugins, err := jclient.GetPlugins(2); err == nil {
		if err = cyclePlugins(plugins); err == nil {
			outString += err.Error()
		}
	} else {
	outString += "  No plugins have lost dependencies...\n"
	}
	err = errors.New(outString)
	return
}

func cyclePlugins(plugins *client.InstalledPluginList) (err error){
	outString := ""
	for _, plugin := range plugins.Plugins {
		outString += "  Checking the plugin " + plugin.ShortName + ": \n"
		dependencies := plugin.Dependencies
		if len(dependencies) != 0 {
			if err = cycleDependencies(dependencies, plugins); err == nil {
				outString += err.Error()
			}
		} else {
			outString += "    The Plugin no dependencies"
		}
	}
	err = errors.New(outString)
	return
}

func cycleDependencies(dependencies []client.Dependence, plugins *client.InstalledPluginList) (err error){
	outString := ""
	for _, dependence := range dependencies {
		outString += "    Checking the dependence plugin " + dependence.ShortName + ": "
		hasInstalled := false
		needUpdate := false
		if err = cycleMatchPlugins(plugins, dependence, hasInstalled, needUpdate); err == nil {
			outString += err.Error()
		}

	}
	err = errors.New(outString)
	return
}

func cycleMatchPlugins(plugins *client.InstalledPluginList, dependence client.Dependence, hasInstalled bool, needUpdate bool) (err error){
	outString := ""
	for _, checkPlugin := range plugins.Plugins {
		checkPluginVersion := strings.Split(checkPlugin.Version, ".")
		dependenceVersion := strings.Split(dependence.Version, ".")
		if checkPlugin.ShortName == dependence.ShortName {
			hasInstalled = true
			if _, err = matchPlugin(dependenceVersion, checkPluginVersion, needUpdate, dependence);err == nil {
				outString += err.Error()
			}

		}
		if needUpdate {
			break
		}
	}
	if !hasInstalled {
		outString += "    The dependence " + dependence.ShortName + " no install, please install it the version " + dependence.Version + " at least\n"
	}
	err = errors.New(outString)
	return
}

func matchPlugin(dependenceVersion []string, checkPluginVersion []string, needUpdate bool, dependence client.Dependence) (isPass bool, err error) {
	outString := ""
	for i := range dependenceVersion {
		if strings.Contains(dependenceVersion[i], "-") && strings.Contains(checkPluginVersion[i], "-") {
			isPass, _ = matchPlugin(strings.Split(dependenceVersion[i], "-"), strings.Split(checkPluginVersion[i], "-"), needUpdate, dependence)
			if isPass {
				break
			}
		} else if len(checkPluginVersion) >= i+1 && len(dependenceVersion) >= i+1 {
			checkPluginVersionInt, _ := strconv.Atoi(checkPluginVersion[i])
			dependenceVersionInt, _ := strconv.Atoi(dependenceVersion[i])
			if checkPluginVersionInt == dependenceVersionInt {
				if i+1 == len(dependenceVersion) {
					isPass = true
					outString += "***true***\n"
					break
				} else {
					continue
				}
			} else if checkPluginVersionInt > dependenceVersionInt {
				isPass = true
				outString += "***true***\n"
				break
			} else {
				isPass = true
				needUpdate = true
				outString += "The dependence " + dependence.ShortName + " need upgrade the version to " + dependence.Version + "\n"
				break
			}
		}
	}
	err = errors.New(outString)
	return
}

// Output renders data into a table
func (o *DoctorOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		buf := new(bytes.Buffer)


		jobCategories := obj.([]client.JobCategory)
		table := util.CreateTable(buf)
		table.AddRow("number", "name", "type")
		for _, jobCategory := range jobCategories {
			for i, item := range jobCategory.Items {
				table.AddRow(fmt.Sprintf("%d", i), item.DisplayName,
					jobCategory.Name)
			}
		}
		table.Render()
		err = nil
		data = buf.Bytes()
	}
	return
}
package cmd

import (
	"fmt"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// DoctorOption is the doctor cmd option
type DoctorOption struct {
	OutputOption
}

var doctorOption DoctorOption

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.PersistentFlags().StringVarP(&jobOption.Format, "output", "o", "json", "Format the output")
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Print the doctor of your Jenkins",
	Long:  `Print the doctor of your Jenkins`,
	Run: func(_ *cobra.Command, _ []string) {
		jenkinsNames := getJenkinsNames()
		checkDuplicateName(jenkinsNames)
		jenkinsServers := getConfig().JenkinsServers
		checkJenkinsServersStatus(jenkinsServers)
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
		fmt.Println("Checked it sure. no duplicated config Name")
	} else {
		fmt.Printf("Duplicate names: %s\n", duplicateName)
	}
}

func checkJenkinsServersStatus(jenkinsServers []JenkinsServer) {
	fmt.Println("Begining checking jenkinsServer status form the configuration files: ")
	for i := range jenkinsServers {
		jenkinsServer := jenkinsServers[i]
		jclient := &client.PluginManager{}
		jclient.URL = jenkinsServer.URL
		jclient.UserName = jenkinsServer.UserName
		jclient.Token = jenkinsServer.Token
		jclient.Proxy = jenkinsServer.Proxy
		jclient.ProxyAuth = jenkinsServer.ProxyAuth
		fmt.Printf("checking the number: %d, name: %s's status now: ", i, jenkinsServer.Name)
		if _, err := jclient.GetPlugins(); err == nil {
			fmt.Println("***ok***")
		} else {
			fmt.Println(err)
		}
	}
}

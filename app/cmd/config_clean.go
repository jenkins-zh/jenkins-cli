package cmd

import (
	cfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"time"
)

// ConfigCleanOption option for config list command
type ConfigCleanOption struct {
	Timeout      int
	CleanTimeout bool

	Logger helper.Printer
}

var configCleanOption ConfigCleanOption

func init() {
	configCmd.AddCommand(configCleanCmd)
	configCleanCmd.Flags().IntVarP(&configCleanOption.Timeout, "timeout", "t", 5,
		i18n.T("Timeout in second value when checking with the Jenkins URL"))
	configCleanCmd.Flags().BoolVarP(&configCleanOption.CleanTimeout, "clean-timeout", "", false,
		i18n.T("Clean the config items when timeout with API request"))
}

// CheckResult is the result of checking
type CheckResult struct {
	Name       string
	StatusCode int
	Timeout    bool
}

var configCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: i18n.T("Clean up some unavailable config items"),
	Long:  i18n.T("Clean up some unavailable config items"),
	RunE:  configCleanOption.Run,
}

// Run is the main entry point of command config-clean
func (o *ConfigCleanOption) Run(cmd *cobra.Command, args []string) (err error) {
	if config := getConfig(); config == nil {
		cmd.Println("cannot found config file")
	}
	o.Logger = cmd

	itemCount := len(config.JenkinsServers)
	checkResult := make(chan CheckResult, itemCount)

	for _, jenkins := range config.JenkinsServers {
		go func(target cfg.JenkinsServer) {
			checkResult <- o.Check(target)
		}(jenkins)
	}

	checkResultList := make([]CheckResult, itemCount)
	for i := range config.JenkinsServers {
		checkResultList[i] = <-checkResult
	}

	// do the clean work
	err = o.CleanByCondition(checkResultList)
	cmd.Println()
	return
}

// Check check the target JenkinsServer config
// make a request to a Jenkins API
func (o *ConfigCleanOption) Check(jenkins cfg.JenkinsServer) (result CheckResult) {
	result.Name = jenkins.Name

	jClient := &client.PluginManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: pluginListOption.RoundTripper,
			Timeout:      time.Duration(o.Timeout),
		},
	}
	jClient.URL = jenkins.URL
	jClient.UserName = jenkins.UserName
	jClient.Token = jenkins.Token
	jClient.Proxy = jenkins.Proxy
	jClient.ProxyAuth = jenkins.ProxyAuth
	jClient.InsecureSkipVerify = jenkins.InsecureSkipVerify

	var statusCode int
	var err error
	if statusCode, _, err = jClient.Request(http.MethodGet, "/api/json", nil, nil); err != nil {
		if err, ok := err.(net.Error); ok {
			result.Timeout = err.Timeout()
		} else {
			o.Logger.Println("check request failed, error is", err)
		}
	}
	result.StatusCode = statusCode
	return
}

// CleanByCondition do the clean work by conditions
func (o *ConfigCleanOption) CleanByCondition(resultList []CheckResult) (err error) {
	if len(resultList) == 0 {
		return
	}

	for _, result := range resultList {
		if o.CleanTimeout && result.Timeout {
			if err = removeJenkins(result.Name); err == nil {
				o.Logger.Printf("removed invalid item %s due to timeout reason\n", result.Name)
			}
		} else {
			o.Logger.Printf("status code of %s is %d\n", result.Name, result.StatusCode)
		}
	}
	return
}

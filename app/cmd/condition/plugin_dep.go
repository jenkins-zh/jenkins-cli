package condition

import (
	"fmt"
	"github.com/hashicorp/go-version"
	appCfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/jenkins-zh/jenkins-cli/client"
	"net/http"
)

// PluginDepCheck is the checker of plugin deps
type PluginDepCheck struct {
	client                    *client.PluginManager
	pluginName, targetVersion string
}

// NewChecker returns a plugin dep checker
func NewChecker(jenkins *appCfg.JenkinsServer, roundTripper http.RoundTripper, pluginName, targetVersion string) (
	checker *PluginDepCheck) {
	checker = &PluginDepCheck{
		pluginName:    pluginName,
		targetVersion: targetVersion,
	}

	jClient := &client.PluginManager{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: roundTripper,
		},
	}
	jClient.URL = jenkins.URL
	jClient.UserName = jenkins.UserName
	jClient.Token = jenkins.Token
	jClient.Proxy = jenkins.Proxy
	jClient.ProxyAuth = jenkins.ProxyAuth
	jClient.InsecureSkipVerify = jenkins.InsecureSkipVerify
	checker.client = jClient
	return
}

// FindPlugin find a plugin by name
func (p *PluginDepCheck) FindPlugin(name string) (plugin *client.InstalledPlugin, err error) {
	if plugin, err = p.client.FindInstalledPlugin(name); err == nil && plugin == nil {
		err = fmt.Errorf(fmt.Sprintf("lack of plugin %s", name))
	}
	return
}

// Check check if the target plugin with a specific version does exists
func (p *PluginDepCheck) Check() (err error) {
	var plugin *client.InstalledPlugin
	if plugin, err = p.FindPlugin(p.pluginName); err == nil {
		var (
			current      *version.Version
			target       *version.Version
			versionMatch bool
		)

		if current, err = version.NewVersion(plugin.Version); err == nil {
			if target, err = version.NewVersion(p.targetVersion); err == nil {
				versionMatch = current.GreaterThanOrEqual(target)
			}
		}

		if err == nil && !versionMatch {
			err = fmt.Errorf("%s version is %s, should be %s", p.pluginName, plugin.Version, p.targetVersion)
		}
	}
	return
}

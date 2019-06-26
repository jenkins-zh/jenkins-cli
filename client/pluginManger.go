package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type PluginManager struct {
	JenkinsCore
}

// PluginList represent a list of plugins
type PluginList struct {
	Plugins []Plugin
}

// Plugin represent the plugin from Jenkins
type Plugin struct {
	Active             bool
	ShortName          string
	LongName           string
	Version            string
	URL                string
	HasUpdate          bool
	Enable             bool
	Downgradable       bool
	Pinned             bool
	RequiredCoreVesion string
	MinimumJavaVersion string
	SupportDynamicLoad string
	Deleted            bool
	Bundled            bool
	BackVersion        string
}

// CheckUpdate fetch the lastest plugins from update center site
func (p *PluginManager) CheckUpdate(handle func(*http.Response)) {
	api := fmt.Sprintf("%s/pluginManager/checkUpdatesServer", p.URL)
	req, err := http.NewRequest("POST", api, nil)
	if err == nil {
		p.AuthHandle(req)
	} else {
		log.Fatal(err)
	}

	if err = p.CrumbHandle(req); err != nil {
		log.Fatal(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if response, err := client.Do(req); err == nil {
		p.handleCheck(handle)(response)
	} else {
		log.Fatal(err)
	}
}

func (p *PluginManager) GetPlugins() (pluginList *PluginList, err error) {
	api := fmt.Sprintf("%s/pluginManager/api/json?pretty=true&depth=1", p.URL)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("GET", api, nil)
	if err == nil {
		p.AuthHandle(req)
	} else {
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			if err == nil {
				pluginList = &PluginList{}
				err = json.Unmarshal(data, pluginList)
			}
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

// InstallPlugin install a plugin by name
func (p *PluginManager) InstallPlugin(names []string) (err error) {
	for i, name := range names {
		names[i] = fmt.Sprintf("plugin.%s", name)
	}
	api := fmt.Sprintf("%s/pluginManager/install?%s", p.URL, strings.Join(names, "=&"))
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("POST", api, nil)
	if err == nil {
		p.AuthHandle(req)
	} else {
		return
	}

	if err = p.CrumbHandle(req); err != nil {
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			fmt.Println("install succeed.")
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

// UninstallPlugin uninstall a plugin by name
func (p *PluginManager) UninstallPlugin(name string) (err error) {
	api := fmt.Sprintf("%s/pluginManager/plugin/%s/uninstall", p.URL, name)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("POST", api, nil)
	if err == nil {
		p.AuthHandle(req)
	} else {
		return
	}

	if err = p.CrumbHandle(req); err != nil {
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			fmt.Println("uninstall succeed.")
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (p *PluginManager) handleCheck(handle func(*http.Response)) func(*http.Response) {
	if handle == nil {
		handle = func(*http.Response) {}
	}
	return handle
}

package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

// PrepareForOneInstalledPluginWithPluginName only for test
func PrepareForOneInstalledPluginWithPluginName(roundTripper *mhttp.MockRoundTripper, rootURL, pluginName string) (
	request *http.Request, response *http.Response) {
	request, response = PrepareForEmptyInstalledPluginList(roundTripper, rootURL, 1)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
			"plugins": [{
				"shortName": "%s",
				"version": "1.0",
				"hasUpdate": true,
				"enable": true,
				"active": true
			}]
		}`, pluginName)))
	return
}

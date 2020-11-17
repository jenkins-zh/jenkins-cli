package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// CenterLoginOption option for upgrade Jenkins
type CenterLoginOption struct {
	common.Option
	RoundTripper http.RoundTripper
}

var centerLoginOption CenterLoginOption

func init() {
	centerCmd.AddCommand(centerLoginCmd)
	healthCheckRegister.Register(getCmdPath(centerLoginCmd), &centerLoginOption)
}

// Check do the health check of center login cmd
func (o *CenterLoginOption) Check() (err error) {
	opt := PluginOptions{
		Option: common.Option{RoundTripper: o.RoundTripper},
	}
	const pluginName = "pipeline-restful-api"
	const targetVersion = "0.11"
	var plugin *client.InstalledPlugin
	if plugin, err = opt.FindPlugin(pluginName); err == nil {
		var (
			current      *version.Version
			target       *version.Version
			versionMatch bool
		)

		if current, err = version.NewVersion(plugin.Version); err == nil {
			if target, err = version.NewVersion(targetVersion); err == nil {
				versionMatch = current.GreaterThanOrEqual(target)
			}
		}

		if err == nil && !versionMatch {
			err = fmt.Errorf("%s version is %s, should be %s", pluginName, plugin.Version, targetVersion)
		}
	}
	return
}

var centerLoginCmd = &cobra.Command{
	Use:   "login",
	Short: i18n.T("Login Jenkins and fetch the token"),
	Long:  i18n.T("Login Jenkins and fetch the token"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		port := 17890
		srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}
		httpServerDone := &sync.WaitGroup{}

		jenkins := getCurrentJenkins()
		fmt.Fprintf(cmd.OutOrStdout(), "Try to login %s, %s\n", jenkins.Name, jenkins.URL)

		httpServerDone.Add(1)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var body []byte
			if body, err = ioutil.ReadAll(r.Body); err == nil {
				token := &client.Token{}
				if err = json.Unmarshal(body, token); err == nil {
					for i, cfg := range config.JenkinsServers {
						if cfg.Name == jenkins.Name {
							config.JenkinsServers[i].Token = token.Data.TokenValue
							err = saveConfig()

							fmt.Fprintln(cmd.OutOrStdout(), "Good to go!")
							break
						}
					}
				}
			}
			srv.Close()
			httpServerDone.Done()
		})
		go func() {
			srv.ListenAndServe()
		}()

		callback := fmt.Sprintf(jenkins.URL+"/jcliPluginManager/test?callback=http://localhost:%d", port)

		util.Open(callback, "", nil)
		httpServerDone.Wait()
		return
	},
}

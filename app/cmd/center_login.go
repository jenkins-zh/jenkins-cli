package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/condition"
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

	if jenkins := getCurrentJenkinsFromOptions(); jenkins != nil {
		healthCheckRegister.Register(getCmdPath(centerLoginCmd), condition.NewChecker(jenkins, centerLoginOption.RoundTripper,
			"pipeline-restful-api", "0.11"))
	}
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
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Try to login %s, %s\n", jenkins.Name, jenkins.URL)

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

							_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Good to go!")
							break
						}
					}
				}
			}
			if srvErr := srv.Close(); srvErr != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Cannot close the local server: %v. %v", srv, srvErr)
			}
			httpServerDone.Done()
		})
		go func() {
			if srvErr := srv.ListenAndServe(); srvErr != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Got error when starting a local server: %v. %v", srv, srvErr)
			}
		}()

		callback := fmt.Sprintf(jenkins.URL+"/jcliPluginManager/test?callback=http://localhost:%d", port)

		_ = util.Open(callback, "", nil)
		httpServerDone.Wait()
		return
	},
}

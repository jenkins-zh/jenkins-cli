package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// CommonOption contains the common options
type CommonOption struct {
	SystemCallExec util.SystemCallExec
	RoundTripper   http.RoundTripper
}

// OutputOption represent the format of output
type OutputOption struct {
	Format string

	WithoutHeaders bool
}

// FormatOutput is the interface of format output
type FormatOutput interface {
	Output(obj interface{}, format string) (data []byte, err error)
}

const (
	// JSONOutputFormat is the format of json
	JSONOutputFormat string = "json"
	// YAMLOutputFormat is the format of yaml
	YAMLOutputFormat string = "yaml"
	// TableOutputFormat is the format of table
	TableOutputFormat string = "table"
)

// Output print the object into byte array
func (o *OutputOption) Output(obj interface{}) (data []byte, err error) {
	switch o.Format {
	case JSONOutputFormat:
		return json.MarshalIndent(obj, "", "  ")
	case YAMLOutputFormat:
		return yaml.Marshal(obj)
	}

	return nil, fmt.Errorf("not support format %s", o.Format)
}

// SetFlag set flag of output format
func (o *OutputOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Format, "output", "o", TableOutputFormat, "Format the output, supported formats: table, json, yaml")
	cmd.Flags().BoolVarP(&o.WithoutHeaders, "no-headers", "", false,
		`When using the default output format, don't print headers (default print headers)`)
}

// BatchOption represent the options for a batch operation
type BatchOption struct {
	Batch bool
}

// Confirm prompte user if they really want to do this
func (b *BatchOption) Confirm(message string) bool {
	if !b.Batch {
		confirm := false
		prompt := &survey.Confirm{
			Message: message,
		}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			return false
		}
	}

	return true
}

// SetFlag the flag for batch option
func (b *BatchOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&b.Batch, "batch", "b", false, "Batch mode, no need confirm")
}

// WatchOption for the resources which can be watched
type WatchOption struct {
	Watch    bool
	Interval int
	Count    int
}

// SetFlag for WatchOption
func (o *WatchOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&o.Interval, "interval", "i", 1, "Interval of watch")
	cmd.Flags().IntVarP(&o.Count, "count", "", 9999, "Count of watch")
}

// InteractiveOption allow user to choose whether the mode is interactive
type InteractiveOption struct {
	Interactive bool
}

// SetFlag set the option flag to this cmd
func (b *InteractiveOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&b.Interactive, "interactive", "i", false,
		i18n.T("Interactive mode"))
}

// HookOption is the option whether skip command hook
type HookOption struct {
	SkipPreHook  bool
	SkipPostHook bool
}

func getCurrentJenkinsAndClientOrDie(jclient *client.JenkinsCore) (jenkins *JenkinsServer) {
	jenkins = getCurrentJenkinsFromOptionsOrDie()
	jclient.URL = jenkins.URL
	jclient.UserName = jenkins.UserName
	jclient.Token = jenkins.Token
	jclient.Proxy = jenkins.Proxy
	jclient.ProxyAuth = jenkins.ProxyAuth
	return
}

func getCurrentJenkinsAndClient(jClient *client.JenkinsCore) (jenkins *JenkinsServer) {
	if jenkins = getCurrentJenkinsFromOptions(); jenkins != nil {
		jClient.URL = jenkins.URL
		jClient.UserName = jenkins.UserName
		jClient.Token = jenkins.Token
		jClient.Proxy = jenkins.Proxy
		jClient.ProxyAuth = jenkins.ProxyAuth
	}
	return
}

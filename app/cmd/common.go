package cmd

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

const (
	since = "since"
)

// CommonOption contains the common options
type CommonOption struct {
	ExecContext     util.ExecContext
	SystemCallExec  util.SystemCallExec
	LookPathContext util.LookPathContext
	RoundTripper    http.RoundTripper
}

// OutputOption represent the format of output
type OutputOption struct {
	Format string

	Columns        string
	WithoutHeaders bool
	Filter         []string

	Writer        io.Writer
	CellRenderMap map[string]RenderCell
}

// RenderCell render a specific cell in a table
type RenderCell = func(string) string

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
// Deprecated see also OutputV2
func (o *OutputOption) Output(obj interface{}) (data []byte, err error) {
	switch o.Format {
	case JSONOutputFormat:
		return json.MarshalIndent(obj, "", "  ")
	case YAMLOutputFormat:
		return yaml.Marshal(obj)
	}

	return nil, fmt.Errorf("not support format %s", o.Format)
}

// OutputV2 print the data line by line
func (o *OutputOption) OutputV2(obj interface{}) (err error) {
	if o.Writer == nil {
		err = fmt.Errorf("no writer found")
		return
	}

	if len(o.Columns) == 0 {
		err = fmt.Errorf("no columns found")
		return
	}

	logger.Debug("start to output", zap.Any("filter", o.Filter))
	obj = o.ListFilter(obj)

	var data []byte
	switch o.Format {
	case JSONOutputFormat:
		data, err = json.MarshalIndent(obj, "", "  ")
	case YAMLOutputFormat:
		data, err = yaml.Marshal(obj)
	case TableOutputFormat, "":
		table := util.CreateTableWithHeader(o.Writer, o.WithoutHeaders)
		table.AddHeader(strings.Split(o.Columns, ",")...)
		items := reflect.ValueOf(obj)
		for i := 0; i < items.Len(); i++ {
			table.AddRow(o.GetLine(items.Index(i))...)
		}
		table.Render()
	default:
		err = fmt.Errorf("not support format %s", o.Format)
	}

	if err == nil && len(data) > 0 {
		_, err = o.Writer.Write(data)
	}
	return
}

// ListFilter filter the data list by fields
func (o *OutputOption) ListFilter(obj interface{}) interface{} {
	if len(o.Filter) == 0 {
		return obj
	}

	elemType := reflect.TypeOf(obj).Elem()
	elemSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 10)
	items := reflect.ValueOf(obj)
	for i := 0; i < items.Len(); i++ {
		item := items.Index(i)
		if o.Match(item) {
			elemSlice = reflect.Append(elemSlice, item)
		}
	}
	return elemSlice.Interface()
}

// Match filter an item
func (o *OutputOption) Match(item reflect.Value) bool {
	if len(o.Filter) == 0 {
		return true
	}

	for _, f := range o.Filter {
		arr := strings.Split(f, "=")
		if len(arr) < 2 {
			continue
		}

		key := arr[0]
		val := arr[1]

		if !strings.Contains(util.ReflectFieldValueAsString(item, key), val) {
			continue
		} else {
			return true
		}
	}
	return false
}

// GetLine returns the line of a table
func (o *OutputOption) GetLine(obj reflect.Value) []string {
	columns := strings.Split(o.Columns, ",")
	values := make([]string, 0)

	if o.CellRenderMap == nil {
		o.CellRenderMap = make(map[string]RenderCell, 0)
	}

	for _, col := range columns {
		cell := util.ReflectFieldValueAsString(obj, col)
		if renderCell, ok := o.CellRenderMap[col]; ok && renderCell != nil {
			cell = renderCell(cell)
		}

		values = append(values, cell)
	}
	return values
}

// SetFlag set flag of output format
// Deprecated, see also SetFlagWithHeaders
func (o *OutputOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Format, "output", "o", TableOutputFormat,
		i18n.T("Format the output, supported formats: table, json, yaml"))
	cmd.Flags().BoolVarP(&o.WithoutHeaders, "no-headers", "", false,
		i18n.T(`When using the default output format, don't print headers (default print headers)`))
	cmd.Flags().StringArrayVarP(&o.Filter, "filter", "", []string{},
		i18n.T("Filter for the list by fields"))
}

// SetFlagWithHeaders set the flags of output
func (o *OutputOption) SetFlagWithHeaders(cmd *cobra.Command, headers string) {
	o.SetFlag(cmd)
	cmd.Flags().StringVarP(&o.Columns, "columns", "", headers,
		i18n.T("The columns of table"))
}

// BatchOption represent the options for a batch operation
type BatchOption struct {
	Batch bool
}

// Confirm promote user if they really want to do this
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
		jClient.InsecureSkipVerify = jenkins.InsecureSkipVerify
	}
	return
}

// GetAliasesDel returns the aliases for delete command
func GetAliasesDel() []string {
	return []string{"remove", "del"}
}

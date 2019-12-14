package cmd

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"net/http"

	"github.com/hashicorp/go-version"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

// JobSearchOption is the options of job search command
type JobSearchOption struct {
	OutputOption
	Start int
	Limit int
	Name  string
	Type  string

	RoundTripper http.RoundTripper
}

var jobSearchOption JobSearchOption

func init() {
	jobCmd.AddCommand(jobSearchCmd)
	jobSearchCmd.Flags().IntVarP(&jobSearchOption.Start, "start", "", 0,
		i18n.T("The list of items offset"))
	jobSearchCmd.Flags().IntVarP(&jobSearchOption.Limit, "limit", "", 50,
		i18n.T("The list of items limit"))
	jobSearchCmd.Flags().StringVarP(&jobSearchOption.Name, "name", "", "",
		i18n.T("The name of plugin for search"))
	jobSearchCmd.Flags().StringVarP(&jobSearchOption.Type, "type", "", "",
		i18n.T("The type of plugin for search"))
	jobSearchCmd.Flags().StringVarP(&jobSearchOption.Format, "output", "o", "table",
		i18n.T(`Formats of the output which contain name, url`))

	healthCheckRegister.Register(getCmdPath(jobSearchCmd), &jobSearchOption)
}

var jobSearchCmd = &cobra.Command{
	Use:     "search",
	Short:   i18n.T("Print the job of your Jenkins"),
	Long:    i18n.T(`Print the job of your Jenkins`),
	Example: "jcli job search [keyword] --name keyword --type Folder",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			jobSearchOption.Name = args[0]
		}
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jClient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: jobSearchOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		var items []client.JenkinsItem
		if items, err = jClient.Search(jobSearchOption.Name, jobSearchOption.Type,
			jobSearchOption.Start, jobSearchOption.Limit); err == nil {
			var data []byte
			if data, err = jobSearchOption.Output(items); err == nil {
				cmd.Print(string(data))
			}
		}
		return
	},
}

// Output render data into byte array
func (o *JobSearchOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		items := obj.([]client.JenkinsItem)
		buf := new(bytes.Buffer)
		table := util.CreateTable(buf)
		table.AddRow("number", "name", "displayname", "type", "url")
		for i, item := range items {
			table.AddRow(fmt.Sprintf("%d", i), item.Name, item.DisplayName, item.Type, item.URL)
		}
		table.Render()
		data = buf.Bytes()
		err = nil
	}
	return
}

// Check do the conditions check
func (o *JobSearchOption) Check() (err error) {
	opt := PluginOptions{
		CommonOption: CommonOption{RoundTripper: o.RoundTripper},
	}
	const pluginName = "pipeline-restful-api"
	const targetVersion = "0.3"
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

package cmd

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/jenkins-zh/jenkins-cli/util"
	"net/http"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
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

	Filter []string

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
	jobSearchCmd.Flags().StringVarP(&jobSearchOption.Columns, "columns", "", "Name,DisplayName,Type,URL",
		i18n.T("The columns of table"))
	jobSearchCmd.Flags().StringArrayVarP(&jobSearchOption.Filter, "filter", "", []string{},
		i18n.T("Filter for the list"))
	jobSearchOption.SetFlag(jobSearchCmd)

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
			items = jobSearchOption.ItemsFilter(items)
			jobSearchOption.Writer = cmd.OutOrStdout()
			err = jobSearchOption.OutputV2(items)
		}
		return
	},
}

func (o *JobSearchOption) ItemsFilter(items []client.JenkinsItem) (result []client.JenkinsItem) {
	if len(o.Filter) == 0 {
		result = items
		return
	}

	result = make([]client.JenkinsItem, 0)
	for _, item := range items {
		for _, f := range o.Filter {
			arr := strings.Split(f, "=")
			if len(arr) < 2 {
				continue
			}

			key := arr[0]
			val := arr[1]

			if !strings.Contains(util.GetFieldValueAsString(item, key), val) {
				continue
			}

			result = append(result, item)
		}
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

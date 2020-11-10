package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	jenkinsFormula "github.com/jenkins-zh/jenkins-formulas/pkg/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"net/http"
	"sort"
	"strings"
)

// PluginFormulaOption option for plugin formula command
type PluginFormulaOption struct {
	common.OutputOption

	// OnlyRelease indicated that we only output the release version of plugins
	OnlyRelease bool
	// DockerBuild indicated if build docker image
	DockerBuild  bool
	SortPlugins  bool
	RoundTripper http.RoundTripper
}

var pluginFormulaOption PluginFormulaOption

func init() {
	pluginCmd.AddCommand(pluginFormulaCmd)
	flags := pluginFormulaCmd.Flags()
	flags.BoolVarP(&pluginFormulaOption.OnlyRelease, "only-release", "", true,
		`Indicated that we only output the release version of plugins`)
	flags.BoolVarP(&pluginFormulaOption.DockerBuild, "docker-build", "", false,
		`Indicated if build docker image`)
	flags.BoolVarP(&pluginFormulaOption.SortPlugins, "sort-plugins", "", true,
		`Indicated if sort the plugins by name`)

	healthCheckRegister.Register(getCmdPath(pluginFormulaCmd), &pluginFormulaOption)
}

// Check do the health check of plugin formula cmd
func (o *PluginFormulaOption) Check() (err error) {
	opt := PluginOptions{
		CommonOption: common.CommonOption{RoundTripper: o.RoundTripper},
	}
	_, err = opt.FindPlugin("pipeline-restful-api")
	return
}

var pluginFormulaCmd = &cobra.Command{
	Use:   "formula",
	Short: i18n.T("Print a formula which contains all plugins come from current Jenkins server"),
	Long: i18n.T(`Print a formula which contains all plugins come from current Jenkins server
Want to know more about what's a Jenkins formula? Please visit https://github.com/jenkins-zh/jenkins-formulas'`),
	Example: `Once you generate the formula file by: jcli plugin formula > test.yaml
than you can package the Jenkins distribution by: jcli cwp --config-path test.yaml`,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jClient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginFormulaOption.RoundTripper,
			},
		}
		jCoreClient := &client.JenkinsStatusClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginFormulaOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))
		getCurrentJenkinsAndClient(&(jCoreClient.JenkinsCore))

		var status *client.JenkinsStatus
		if status, err = jCoreClient.Get(); err != nil {
			err = fmt.Errorf("cannot get the version of current Jenkins, error is %v", err)
			return
		}

		// make the formula
		formula := jenkinsFormula.CustomWarPackage{
			Bundle: jenkinsFormula.Bundle{
				GroupId:     "io.github.jenkins-zh",
				ArtifactId:  "jenkins-zh",
				Description: "Jenkins formula generated by jcli",
				Vendor:      "Chinese Jenkins Community",
			},
			BuildSettings: jenkinsFormula.BuildSettings{
				Docker: jenkinsFormula.BuildDockerSetting{
					Build: pluginFormulaOption.DockerBuild,
					Base:  fmt.Sprintf("jenkins/jenkins:%s", status.Version),
					Tag:   "jenkins/jenkins-formula:v0.0.1",
				},
			},
			War: jenkinsFormula.CustomWar{
				GroupId:    "org.jenkins-ci.main",
				ArtifactId: "jenkins-war",
				Source: jenkinsFormula.Source{
					Version: status.Version,
				},
			},
		}
		if err = jClient.GetPluginsFormula(&formula.Plugins); err == nil {
			if pluginFormulaOption.OnlyRelease {
				formula.Plugins = removeSnapshotPlugins(formula.Plugins)
			}

			if pluginFormulaOption.SortPlugins {
				formula.Plugins = SortPlugins(formula.Plugins)
			}

			var data []byte
			if data, err = yaml.Marshal(formula); err == nil {
				_, _ = cmd.OutOrStdout().Write(data)
			}
		}
		return
	},
	Annotations: map[string]string{
		common.Since: common.VersionSince0031,
	},
}

// SortPlugins sort the plugins by asc
func SortPlugins(plugins []jenkinsFormula.Plugin) []jenkinsFormula.Plugin {
	sort.SliceStable(plugins, func(i, j int) bool {
		if strings.Compare(plugins[j].GroupId, plugins[i].GroupId) > 0 {
			return true
		}
		if strings.Compare(plugins[j].ArtifactId, plugins[i].ArtifactId) > 0 {
			return true
		}
		return false
	})
	return plugins
}

func removeSnapshotPlugins(plugins []jenkinsFormula.Plugin) (result []jenkinsFormula.Plugin) {
	result = make([]jenkinsFormula.Plugin, 0)

	for i := range plugins {
		if strings.Contains(plugins[i].Source.Version, "SNAPSHOT") {
			continue
		}

		result = append(result, plugins[i])
	}
	return
}

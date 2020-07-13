package cmd

import (
	"github.com/google/go-github/v29/github"
	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

// VersionOption is the version option
type VersionOption struct {
	Changelog  bool
	ShowLatest bool

	GitHubClient     *github.Client
	JenkinsClient    common.JenkinsClient
	JenkinsConfigMgr common.JenkinsConfigMgr
}

// NewVersionCmd create a command for version
func NewVersionCmd(client common.JenkinsClient, jenkinsConfigMgr common.JenkinsConfigMgr) (cmd *cobra.Command) {
	opt := &VersionOption{}

	cmd = &cobra.Command{
		Use:    "version",
		Short:  "Print the version of Jenkins CLI",
		Long:   `Print the version of Jenkins CLI`,
		PreRun: opt.PreRun,
		RunE:   opt.RunE,
		Annotations: map[string]string{
			common.Since: "v0.0.26",
		},
	}

	flags := cmd.Flags()
	opt.addFlags(flags)

	cmd.AddCommand(NewSelfUpgradeCmd(client, jenkinsConfigMgr))
	return
}

func (o *VersionOption) addFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&o.Changelog, "changelog", "", false,
		i18n.T("Output the changelog of current version"))
	flags.BoolVarP(&o.ShowLatest, "show-latest", "", false,
		i18n.T("Output the latest version"))
}

// PreRun is the pre-check of current command
func (o *VersionOption) PreRun(cmd *cobra.Command, _ []string) {
	if o.GitHubClient == nil {
		o.GitHubClient = github.NewClient(nil)
	}
}

// RunE is the main point of current command
func (o *VersionOption) RunE(cmd *cobra.Command, _ []string) (err error) {
	cmd.Println(i18n.T("Jenkins CLI (jcli) manage your Jenkins"))

	version := app.GetVersion()
	cmd.Printf("Version: %s\n", version)
	cmd.Printf("Last Commit: %s\n", app.GetCommit())
	cmd.Printf("Build Date: %s\n", app.GetDate())

	if strings.HasPrefix(version, "dev-") {
		version = strings.ReplaceAll(version, "dev-", "")
	}

	ghClient := &client.GitHubReleaseClient{
		Client: o.GitHubClient,
	}
	var asset *client.ReleaseAsset
	if o.Changelog {
		if asset, err = ghClient.GetJCLIAsset(version); err == nil && asset != nil {
			cmd.Println()
			cmd.Println(asset.Body)
		}
	} else if o.ShowLatest {
		if asset, err = ghClient.GetLatestJCLIAsset(); err == nil && asset != nil {
			cmd.Println()
			cmd.Println(asset.TagName)
			cmd.Println(asset.Body)
		}
	}
	return
}

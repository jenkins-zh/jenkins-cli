package version

import (
	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

// NewVersionCmd create a command for version
func NewVersionCmd(client common.JenkinsClient, jenkinsConfigMgr common.JenkinsConfigMgr) (cmd *cobra.Command) {
	opt := &VersionPrintOption{
		JenkinsConfigMgr: jenkinsConfigMgr,
	}

	cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of Jenkins CLI",
		Long:  `Print the version of Jenkins CLI`,
		RunE:  opt.RunE,
		Annotations: map[string]string{
			common.Since: "v0.0.26",
		},
	}

	flags := cmd.Flags()
	opt.addFlags(flags)

	cmd.AddCommand(NewSelfUpgradeCmd(client, jenkinsConfigMgr))
	return
}

func (o *VersionPrintOption) addFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&o.Changelog, "changelog", "", false,
		i18n.T("Output the changelog of current version"))
	flags.BoolVarP(&o.ShowLatest, "show-latest", "", false,
		i18n.T("Output the latest version"))
}

// RunE is the main point of current command
func (o *VersionPrintOption) RunE(cmd *cobra.Command, _ []string) (err error) {
	cmd.Println(i18n.T("Jenkins CLI (jcli) manage your Jenkins"))

	version := app.GetVersion()
	cmd.Printf("Version: %s\n", version)
	cmd.Printf("Last Commit: %s\n", app.GetCommit())
	cmd.Printf("Build Date: %s\n", app.GetDate())

	if strings.HasPrefix(version, "dev-") {
		version = strings.ReplaceAll(version, "dev-", "")
	}

	ghClient := &client.GitHubReleaseClient{
		Client: o.JenkinsConfigMgr.GetGitHubClient(),
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

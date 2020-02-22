package cmd

import (
	"fmt"
	"github.com/google/go-github/v29/github"
	"github.com/jenkins-zh/jenkins-cli/app"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"strings"
)

// VersionOption is the version option
type VersionOption struct {
	Changelog  bool
	ShowLatest bool

	GitHubClient *github.Client
}

var versionOption VersionOption

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&versionOption.Changelog, "changelog", "", false,
		i18n.T("Output the changelog of current version"))
	versionCmd.Flags().BoolVarP(&versionOption.ShowLatest, "show-latest", "", false,
		i18n.T("Output the latest version"))
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the user of your Jenkins",
	Long:  `Print the user of your Jenkins`,
	PreRun: func(cmd *cobra.Command, _ []string) {
		if versionOption.GitHubClient == nil {
			versionOption.GitHubClient = github.NewClient(nil)
		}
	},
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		cmd.Println(i18n.T("Jenkins CLI (jcli) manage your Jenkins"))

		version := app.GetVersion()
		cmd.Printf("Version: %s\n", version)
		cmd.Printf("Commit: %s\n", app.GetCommit())

		if rootOptions.Jenkins != "" {
			current := getCurrentJenkinsFromOptions()
			if current != nil {
				cmd.Println("Current Jenkins is:", current.Name)
			} else {
				err = fmt.Errorf("cannot found the configuration: %s", rootOptions.Jenkins)
				return
			}
		}

		if strings.HasPrefix(version, "dev-") {
			version = strings.ReplaceAll(version, "dev-", "")
		}

		ghClient := &client.GitHubReleaseClient{
			Client: versionOption.GitHubClient,
		}
		var asset *client.ReleaseAsset
		if versionOption.Changelog {
			if asset, err = ghClient.GetJCLIAsset(version); err == nil && asset != nil {
				cmd.Println()
				cmd.Println(asset.Body)
			}
		} else if versionOption.ShowLatest {
			if asset, err = ghClient.GetLatestJCLIAsset(); err == nil && asset != nil {
				cmd.Println()
				cmd.Println(asset.TagName)
				cmd.Println(asset.Body)
			}
		}
		return
	},
	Annotations: map[string]string{
		since: "v0.0.26",
	},
}

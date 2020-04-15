package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// PluginUploadOption will hold the options of plugin cmd
type PluginUploadOption struct {
	Remote         string
	RemoteUser     string
	RemotePassword string
	RemoteJenkins  string
	ShowProgress   bool

	RoundTripper http.RoundTripper

	HookOption

	pluginFilePath string
}

var pluginUploadOption PluginUploadOption

func init() {
	pluginCmd.AddCommand(pluginUploadCmd)
	pluginUploadCmd.Flags().BoolVarP(&pluginUploadOption.ShowProgress, "show-progress", "", true,
		i18n.T("Whether show the upload progress"))
	pluginUploadCmd.Flags().StringVarP(&pluginUploadOption.Remote, "remote", "r", "",
		i18n.T("Remote plugin URL"))
	pluginUploadCmd.Flags().StringVarP(&pluginUploadOption.RemoteUser, "remote-user", "", "",
		i18n.T("User of remote plugin URL"))
	pluginUploadCmd.Flags().StringVarP(&pluginUploadOption.RemotePassword, "remote-password", "", "",
		i18n.T("Password of remote plugin URL"))
	pluginUploadCmd.Flags().StringVarP(&pluginUploadOption.RemoteJenkins, "remote-jenkins", "", "",
		i18n.T("Remote Jenkins which will find from config list"))

	pluginUploadCmd.Flags().BoolVarP(&pluginUploadOption.SkipPreHook, "skip-prehook", "", false,
		i18n.T("Whether skip the previous command hook"))
	pluginUploadCmd.Flags().BoolVarP(&pluginUploadOption.SkipPostHook, "skip-posthook", "", false,
		i18n.T("Whether skip the post command hook"))
}

var pluginUploadCmd = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"up"},
	Short:   i18n.T("Upload a plugin  to your Jenkins"),
	Long:    i18n.T(`Upload a plugin from local filesystem or remote URL to your Jenkins`),
	Example: `  jcli plugin upload --remote https://server/sample.hpi
jcli plugin upload sample.hpi
jcli plugin upload sample.hpi --show-progress=false`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if pluginUploadOption.Remote != "" {
			file, err := ioutil.TempFile(".", "jcli-plugin")
			if err != nil {
				log.Fatal(err)
			}

			defer os.Remove(file.Name())

			if pluginUploadOption.RemoteJenkins != "" {
				if jenkins := findJenkinsByName(pluginUploadOption.RemoteJenkins); jenkins != nil {
					pluginUploadOption.RemoteUser = jenkins.UserName
					pluginUploadOption.RemotePassword = jenkins.Token
				}
			}

			pluginUploadOption.pluginFilePath = fmt.Sprintf("%s.hpi", file.Name())
			downloader := util.HTTPDownloader{
				TargetFilePath: pluginUploadOption.pluginFilePath,
				URL:            pluginUploadOption.Remote,
				UserName:       pluginUploadOption.RemoteUser,
				Password:       pluginUploadOption.RemotePassword,
				ShowProgress:   true,
				Debug:          rootOptions.Debug,
			}

			if err := downloader.DownloadFile(); err != nil {
				log.Fatal(err)
			}
		} else if len(args) == 0 {
			if !pluginUploadOption.SkipPreHook {
				executePreCmd(cmd, args, os.Stdout)
			}

			path, _ := os.Getwd()
			dirName := filepath.Base(path)
			dirName = strings.Replace(dirName, "-plugin", "", -1)
			path += fmt.Sprintf("/target/%s.hpi", dirName)

			pluginUploadOption.pluginFilePath = path
		} else {
			pluginUploadOption.pluginFilePath = args[0]
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if pluginUploadOption.SkipPostHook {
			return
		}

		executePostCmd(cmd, args, cmd.OutOrStdout())
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		targetFiles := make([]string, 0)
		if files, err := filepath.Glob("*.hpi"); err == nil {
			targetFiles = append(targetFiles, files...)
		}
		if files, err := filepath.Glob("target/*.hpi"); err == nil {
			targetFiles = append(targetFiles, files...)
		}
		return targetFiles, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, _ []string) {
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginUploadOption.RoundTripper,
				Output:       cmd.OutOrStdout(),
			},
			ShowProgress: pluginUploadOption.ShowProgress,
		}
		getCurrentJenkinsAndClientOrDie(&(jclient.JenkinsCore))
		jclient.Debug = rootOptions.Debug

		if pluginUploadOption.Remote != "" {
			defer os.Remove(pluginUploadOption.pluginFilePath)
		}

		err := jclient.Upload(pluginUploadOption.pluginFilePath)
		helper.CheckErr(cmd, err)
	},
}

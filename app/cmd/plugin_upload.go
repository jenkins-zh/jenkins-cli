package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	FileName       string
	// Timeout is the timeout when upload the plugin
	Timeout int64

	RoundTripper http.RoundTripper

	common.HookOption

	pluginFilePath string
}

var pluginUploadOption PluginUploadOption

func init() {
	pluginCmd.AddCommand(pluginUploadCmd)
	flags := pluginUploadCmd.Flags()

	flags.BoolVarP(&pluginUploadOption.ShowProgress, "show-progress", "", true,
		i18n.T("Whether show the upload progress"))
	flags.StringVarP(&pluginUploadOption.FileName, "file", "f", "",
		i18n.T("The plugin file path which should end with .hpi"))
	flags.StringVarP(&pluginUploadOption.Remote, "remote", "r", "",
		i18n.T("Remote plugin URL"))
	flags.StringVarP(&pluginUploadOption.RemoteUser, "remote-user", "", "",
		i18n.T("User of remote plugin URL"))
	flags.StringVarP(&pluginUploadOption.RemotePassword, "remote-password", "", "",
		i18n.T("Password of remote plugin URL"))
	flags.StringVarP(&pluginUploadOption.RemoteJenkins, "remote-jenkins", "", "",
		i18n.T("Remote Jenkins which will find from config list"))

	flags.BoolVarP(&pluginUploadOption.SkipPreHook, "skip-prehook", "", false,
		i18n.T("Whether skip the previous command hook"))
	flags.BoolVarP(&pluginUploadOption.SkipPostHook, "skip-posthook", "", false,
		i18n.T("Whether skip the post command hook"))

	flags.Int64VarP(&pluginUploadOption.Timeout, "timeout", "", 120,
		"Timeout in second when upload the plugin")

	if err := pluginUploadCmd.RegisterFlagCompletionFunc("file", pluginUploadOption.HPICompletion); err != nil {
		pluginCmd.PrintErrln(err)
	}
}

func (o *PluginUploadOption) HPICompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
}

var pluginUploadCmd = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"up"},
	Short:   i18n.T("Upload a plugin  to your Jenkins"),
	Long:    i18n.T(`Upload a plugin from local filesystem or remote URL to your Jenkins`),
	Example: `  jcli plugin upload --remote https://server/sample.hpi
jcli plugin upload sample.hpi
jcli plugin upload sample.hpi --show-progress=false`,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if pluginUploadOption.Remote != "" {
			var file *os.File
			if file, err = ioutil.TempFile(".", "jcli-plugin"); err != nil {
				return
			}

			defer func() {
				_= os.Remove(file.Name())
			}()

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

			err = downloader.DownloadFile()
		} else if len(args) == 0 {
			if !pluginUploadOption.SkipPreHook {
				if err = executePreCmd(cmd, args, os.Stdout); err != nil {
					return
				}
			}

			path, _ := os.Getwd()
			dirName := filepath.Base(path)
			dirName = strings.Replace(dirName, "-plugin", "", -1)
			path += fmt.Sprintf("/target/%s.hpi", dirName)

			pluginUploadOption.pluginFilePath = path
		} else {
			pluginUploadOption.pluginFilePath = args[0]
		}
		return
	},
	PostRunE: func(cmd *cobra.Command, args []string) (err error) {
		if pluginUploadOption.SkipPostHook {
			return
		}

		err = executePostCmd(cmd, args, cmd.OutOrStdout())
		return
	},
	ValidArgsFunction: pluginUploadOption.HPICompletion,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jclient := &client.PluginManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: pluginUploadOption.RoundTripper,
				Output:       cmd.OutOrStdout(),
				Timeout:      time.Duration(pluginUploadOption.Timeout) * time.Second,
			},
			ShowProgress: pluginUploadOption.ShowProgress,
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if pluginUploadOption.Remote != "" {
			defer func() {
				_ = os.Remove(pluginUploadOption.pluginFilePath)
			}()
		}

		err = jclient.Upload(pluginUploadOption.pluginFilePath)
		return
	},
}

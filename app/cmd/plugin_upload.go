package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
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

	pluginFilePath string
}

var pluginUploadOption PluginUploadOption

func init() {
	pluginCmd.AddCommand(pluginUploadCmd)
	pluginUploadCmd.Flags().StringVarP(&pluginUploadOption.Remote, "remote", "r", "", "Remote plugin URL")
	pluginUploadCmd.Flags().StringVarP(&pluginUploadOption.RemoteUser, "remote-user", "", "", "User of remote plugin URL")
	pluginUploadCmd.Flags().StringVarP(&pluginUploadOption.RemotePassword, "remote-password", "", "", "Password of remote plugin URL")
	pluginUploadCmd.Flags().StringVarP(&pluginUploadOption.RemoteJenkins, "remote-jenkins", "", "", "Remote Jenkins which will find from config list")
}

var pluginUploadCmd = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"up"},
	Short:   "Upload a plugin  to your Jenkins",
	Long:    `Upload a plugin from local filesystem or remote URL to your Jenkins`,
	Example: `  jcli plugin upload --remote https://server/sample.hpi
  jcli plugin upload sample.hpi`,
	PreRun: func(_ *cobra.Command, args []string) {
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
			path, _ := os.Getwd()
			dirName := filepath.Base(path)
			dirName = strings.Replace(dirName, "-plugin", "", -1)
			path += fmt.Sprintf("/target/%s.hpi", dirName)

			pluginUploadOption.pluginFilePath = path
		} else {
			pluginUploadOption.pluginFilePath = args[0]
		}
	},
	Run: func(_ *cobra.Command, _ []string) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.PluginManager{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth
		jclient.Debug = rootOptions.Debug

		if pluginUploadOption.Remote != "" {
			defer os.Remove(pluginUploadOption.pluginFilePath)
		}

		jclient.Upload(pluginUploadOption.pluginFilePath)
	},
}

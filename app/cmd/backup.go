package cmd

import (
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

type BackupOption struct {
	RoundTripper http.RoundTripper
}

var backupOption BackupOption

func init() {
	rootCmd.AddCommand(backupCmd)
	healthCheckRegister.Register(getCmdPath(backupCmd), &backupOption)
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: i18n.T("Backup your global and job specific configurations."),
	Long: i18n.T(`Backup your global and job specific configurations. 
	This command uses Thin Backup Plugin, so please be sure that you have the plugin installed in your jenkins.`),
	RunE: backupOption.Backup,
}

//Check will find out whether Thin Backup Plugin installed or not
func (o *BackupOption) Check() (err error) {
	opt := PluginOptions{
		Option: common.Option{RoundTripper: o.RoundTripper},
	}
	_, err = opt.FindPlugin("thinBackup")
	return
}

//Backup will trigger thinBackup plugin to make a backup
func (o *BackupOption) Backup(cmd *cobra.Command, _ []string) (err error) {
	jClient := &client.CoreClient{
		JenkinsCore: client.JenkinsCore{
			RoundTripper: backupOption.RoundTripper,
		},
	}
	GetCurrentJenkinsAndClient(&(jClient.JenkinsCore))
	cmd.Println("Please wait while making a backup.")
	_, err = jClient.RequestWithoutData(http.MethodGet, "/thinBackup/backupManual", nil, nil, 200)
	if err != nil {
		cmd.Println("Backup failed. Please see logging message for detailed reasons.")
	} else {
		cmd.Println("Backup succeeds.")
	}
	return
}

package cmd

import (
	"net/http"
	"os"
	"time"

	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
)

//BackupOption is an option for backup
type BackupOption struct {
	RoundTripper http.RoundTripper
	BackupDir    string
	WaitTime     int
}

var backupOption BackupOption

func init() {
	rootCmd.AddCommand(backupCmd)
	healthCheckRegister.Register(getCmdPath(backupCmd), &backupOption)
	backupCmd.Flags().StringVarP(&backupOption.BackupDir, "backup-dir", "d", "", i18n.T("the backup directory in thinBackup setting"))
	backupCmd.Flags().IntVarP(&backupOption.WaitTime, "wait-time", "t", 300, i18n.T("the maximum time you would like to wait for jenkins to be idle and make a backup"))
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
	err = ThinBackupAPI(jClient)
	if err != nil {
		cmd.Println("Backup failed. Please see logging message(/yourJenkinsHome/backup.log) for detailed reasons.")
	}
	i := 0
	now := time.Now()
	for o.WaitTime-i*15 > 0 {
		time.Sleep(15 * time.Second)
		stat, err := os.Stat(o.BackupDir)
		if err != nil {
			//check if the directory be modified and it is a hint that a backup may succeed
			if stat.ModTime().After(now) {
				modifyTime := stat.ModTime().Format("03:04:05")
				hour := modifyTime[0:2]
				min := modifyTime[3:5]
				//check if thinBackup plugin createed the bakcup directory
				_, err = os.Stat("*" + stat.ModTime().Format("2021-09-21") + "_" + hour + "-" + min)
				if err == nil {
					cmd.Println("Back up successfully!")
					return nil
				}
			}
		}

	}
	cmd.Println("Backup failed or thinBackup is still waiting for jenkins to be idle to make a backup. Please see logging message(/yourJenkinsHome/backup.log) for detailed reasons.")
	return
}

//ThinBackupAPI requests backupManual api
func ThinBackupAPI(client *client.CoreClient) (err error) {
	_, err = client.RequestWithoutData(http.MethodGet, "/thinBackup/backupManual", nil, nil, 200)
	return err
}

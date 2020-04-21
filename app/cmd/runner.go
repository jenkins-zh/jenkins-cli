package cmd

import (

	"fmt"
	"log"
	"strings"
	"path/filepath"
	"go.uber.org/zap"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"github.com/mitchellh/go-homedir"
	"os"


)

// RunnerOption is the wrapper of jenkinsfile runner cli
type RunnerOption struct {
	BatchOption
	CommonOption

	Safe bool
	Version string
	WarPath string 
}

var runnerOption RunnerOption

func init() {
	rootCmd.AddCommand(runnerCmd)
	runnerOption.SetFlag(runnerCmd)
	runnerCmd.Flags().StringVarP(&runnerOption.WarPath, "path", "p", "",
	i18n.T("The jenkins.war path"))
	runnerCmd.Flags().BoolVarP(&runnerOption.Safe, "safe", "s", true,
	i18n.T("Puts Jenkins into the quiet mode, wait for existing builds to be completed, and then restart Jenkins"))
	runnerCmd.Flags().StringVarP(&runnerOption.Version, "version", "v", "2.190.3",
		i18n.T("The of version of jenkins.war"))
	runnerOption.BatchOption.Stdio = GetSystemStdio()
	runnerOption.CommonOption.Stdio = GetSystemStdio()
}

var runnerCmd = &cobra.Command{
	Use:   "runner",
	Short: i18n.T("The wrapper of jenkinsfile runner"),
	Long: i18n.T(`The wrapper of jenkinsfile runner
Get more about jenkinsfile runner from https://github.com/jenkinsci/jenkinsfile-runner`),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
	 
	//Start by downloading the mirror for the jenkinsfileRunner
	jenkinsfileRunnerVersion := "jenkinsfile-runner-1.0-beta-11.jar"
	jenkinsfileRunnerURL := fmt.Sprintf("https://repo.jenkins-ci.org/list/releases/io/jenkins/jenkinsfile-runner/jenkinsfile-runner/1.0-beta-11/%s", jenkinsfileRunnerVersion) 
	logger.Info("Prepare to start Downloading jenkinfileRunner", zap.String("URL", jenkinsfileRunnerURL))
	downloader:= util.HTTPDownloader{
		URL: jenkinsfileRunnerURL,
		ShowProgress: true,
		TargetFilePath: jenkinsfileRunnerVersion,
	}
	if err := downloader.DownloadFile(); err != nil {
		//Fatal error has occured while downloading the file.
		log.Fatal(err)
	}

		if runnerOption.WarPath == "" {
			// If it does not exist download the jenkins.war
			var userHome string
			if userHome, err = homedir.Dir(); err != nil {
				return err
			}
			jenkinsWar := fmt.Sprintf("%s/.jenkins-cli/cache/%s/jenkins.war", userHome, runnerOption.Version)
			logger.Info("prepare to download jenkins.war as pre-requisite for jfr", zap.String("localPath", jenkinsWar))

			if _, fileErr := os.Stat(jenkinsWar); fileErr != nil {
				download := &CenterDownloadOption{
					Mirror:       "default",
					Output:       jenkinsWar,
					ShowProgress: true,
					Version:      runnerOption.Version,
				}
				if err = download.DownloadJenkins(); err != nil {
					return err
				}
			}
		}

	//Check if jenkins war path exists
	if filepath.Ext(strings.TrimSpace(runnerOption.WarPath)) != ".war" {
		return fmt.Errorf("incorrect file path : %s", runnerOption.WarPath)
	}
	

	return nil
	},
}
package cmd

import (

	"fmt"
	"go.uber.org/zap"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
	"github.com/mitchellh/go-homedir"

)

// RunnerOption is the wrapper of jenkinsfile runner cli
type RunnerOption struct {
	BatchOption
	CommonOption

	Safe bool
}

var runnerOption RunnerOption

func init() {
	rootCmd.AddCommand(runnerCmd)
	runnerOption.SetFlag(runnerCmd)
	runnerCmd.Flags().BoolVarP(&runnerOption.Safe, "safe", "s", true,
	i18n.T("Puts Jenkins into the quiet mode, wait for existing builds to be completed, and then restart Jenkins"))
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
	downloader.DownloadFile()

	//Check if jenkins war exists
	var userHome string
		if userHome, err = homedir.Dir(); err != nil {
			return
		}
	searchDir := fmt.Sprintf("%s/.jenkins-cli/cache/%s/jenkins.war", userHome, centerDownloadOption.Version)
	return
	},
}
package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunnerOption is the wrapper of jenkinsfile runner cli
type RunnerOption struct {
	RoundTripper http.RoundTripper

	Safe            bool
	WarVersion      string
	WarPath         string
	PluginPath      string
	JenkinsfilePath string
	JfrVersion      string
}

var runnerOption RunnerOption

func init() {
	rootCmd.AddCommand(runnerCmd)
	runnerOption.SetFlag(runnerCmd)
	runnerCmd.Flags().StringVarP(&runnerOption.WarPath, "path", "w", "",
		i18n.T("The jenkins.war path"))
	runnerCmd.Flags().BoolVarP(&runnerOption.Safe, "safe", "s", true,
		i18n.T("Puts Jenkins into the quiet mode, wait for existing builds to be completed, and then restart Jenkins"))
	runnerCmd.Flags().StringVarP(&runnerOption.WarVersion, "war-version", "v", "2.190.3",
		i18n.T("The of version of jenkins.war"))
	runnerCmd.Flags().StringVarP(&runnerOption.PluginPath, "plugin-path", "p", "",
		i18n.T("The path to plugins.txt"))
	runnerCmd.Flags().StringVarP(&runnerOption.JenkinsfilePath, "jenkinsfile-path", "j", "",
		i18n.T("The path to jenkinsfile"))
	runnerCmd.Flags().StringVarP(&runnerOption.JfrVersion, "jfr-version", "f", "1.0-beta-11",
		i18n.T("The path to jenkinsfile"))
}

var runnerCmd = &cobra.Command{
	Use:   "runner",
	Short: i18n.T("The wrapper of jenkinsfile runner"),
	Long: i18n.T(`The wrapper of jenkinsfile runner
Get more about jenkinsfile runner from https://github.com/jenkinsci/jenkinsfile-runner`),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {

		var userHome string
		if userHome, err = homedir.Dir(); err != nil {
			return err
		}
		//Start by downloading the mirror for the jenkinsfileRunner
		jenkinsfileRunnerVersion := fmt.Sprintf("jenkinsfile-runner-%s.jar", runnerOption.JfrVersion)
		jenkinsfileRunnerTargetPath := fmt.Sprintf("%s/.jenkins-cli/cache/%s/%s", userHome, runnerOption.WarVersion, jenkinsfileRunnerVersion)
		jenkinsfileRunnerURL := fmt.Sprintf("https://repo.jenkins-ci.org/list/releases/io/jenkins/jenkinsfile-runner/jenkinsfile-runner/1.0-beta-11/%s", jenkinsfileRunnerVersion)
		logger.Info("Prepare to start Downloading jenkinfileRunner", zap.String("URL", jenkinsfileRunnerURL))
		downloader := util.HTTPDownloader{
			URL:            jenkinsfileRunnerURL,
			ShowProgress:   true,
			TargetFilePath: jenkinsfileRunnerTargetPath,
		}
		if err := downloader.DownloadFile(); err != nil {
			//Fatal error has occured while downloading the file.
			log.Fatal(err)
		}

		if runnerOption.WarPath == "" {
			// If it does not exist download the jenkins.war

			jenkinsWar := fmt.Sprintf("%s/.jenkins-cli/cache/%s/jenkins.war", userHome, runnerOption.WarVersion)
			logger.Info("prepare to download jenkins.war as pre-requisite for jfr", zap.String("localPath", jenkinsWar))

			if _, fileErr := os.Stat(jenkinsWar); fileErr != nil {
				download := &CenterDownloadOption{
					Mirror:       "default",
					Output:       jenkinsWar,
					ShowProgress: true,
					Version:      runnerOption.WarVersion,
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

		//Check if plugin.txt path is provided
		if runnerOption.PluginPath == "" {
			return fmt.Errorf("plugin.txt path is not provided kinldy proivde path")
		}

		//Check if the plugin file has a valid extension
		if filepath.Ext(strings.TrimSpace(runnerOption.WarPath)) != ".txt" {
			return fmt.Errorf("incorrect file type it should be a txt file : %s", runnerOption.PluginPath)
		}

		//Check if jenkinsfile path is empty
		if runnerOption.PluginPath == "" {
			return fmt.Errorf("kindly provide valid path to jenkinsfile")
		}

		if filepath.Base(runnerOption.JenkinsfilePath) != "Jenkinsfile" {
			return fmt.Errorf("invalid file type. kindly provide a jenkinsfile")
		}

		//TO-DO
		/*
		a) Build JFR using mvn clean package
		b) Run jfr using arguments
		*/

		return nil
	},
}

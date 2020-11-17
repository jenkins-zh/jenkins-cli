package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RunnerOption is the wrapper of jenkinsfile runner cli
type RunnerOption struct {
	common.BatchOption
	common.Option
	RoundTripper http.RoundTripper

	Safe            bool
	WarVersion      string
	WarPath         string
	PluginPath      string
	JenkinsfilePath string
	JfrVersion      string
	LTS             bool
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
	runnerCmd.Flags().StringVarP(&runnerOption.JenkinsfilePath, "jenkinsfile-path", "z", "",
		i18n.T("The path to jenkinsfile"))
	runnerCmd.Flags().StringVarP(&runnerOption.JfrVersion, "jfr-version", "f", "1.0-beta-11",
		i18n.T("The path to jenkinsfile"))
	runnerCmd.Flags().BoolVarP(&runnerOption.LTS, "lts", "", true,
		i18n.T("If you want to download Jenkins as LTS"))
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
		jenkinsfileRunnerURL := fmt.Sprintf("https://repo.jenkins-ci.org/list/releases/io/jenkins/jenkinsfile-runner/jenkinsfile-runner/%s/%s", runnerOption.JfrVersion, jenkinsfileRunnerVersion)
		logger.Info("Prepare to start Downloading jenkinfileRunner", zap.String("URL", jenkinsfileRunnerURL))

		if _, fErr := os.Stat(jenkinsfileRunnerTargetPath); fErr != nil && os.IsNotExist(fErr) {
			downloader := util.HTTPDownloader{
				URL:            jenkinsfileRunnerURL,
				ShowProgress:   true,
				TargetFilePath: jenkinsfileRunnerTargetPath,
			}
			if err = downloader.DownloadFile(); err != nil {
				return
			}
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
					LTS:          runnerOption.LTS,
					Version:      runnerOption.WarVersion,
				}
				if err = download.DownloadJenkins(); err != nil {
					return err
				}
			}
		}

		//Check if jenkins war path exists
		if runnerOption.WarPath != "" && filepath.Ext(strings.TrimSpace(runnerOption.WarPath)) != ".war" {
			return fmt.Errorf("incorrect file path : %s", runnerOption.WarPath)
		}

		//Check if plugin.txt path is provided
		if runnerOption.PluginPath == "" {
			return fmt.Errorf("plugin.txt path is not provided kinldy proivde path")
		}

		//Check if the plugin file has a valid extension
		if filepath.Ext(strings.TrimSpace(runnerOption.PluginPath)) != ".txt" {
			return fmt.Errorf("incorrect file type it should be a txt file : %s", runnerOption.PluginPath)
		}

		//Check if jenkinsfile path is empty
		if runnerOption.PluginPath == "" {
			return fmt.Errorf("kindly provide valid path to jenkinsfile")
		}

		if filepath.Base(runnerOption.JenkinsfilePath) != "Jenkinsfile" {
			return fmt.Errorf("invalid file type. kindly provide a jenkinsfile")
		}

		var binary string
		binary, err = util.LookPath("java", centerStartOption.LookPathContext)
		if err == nil {
			jenkinsWarArgs := []string{"java", "-jar", jenkinsfileRunnerTargetPath, "-w", runnerOption.WarPath, "-p", runnerOption.PluginPath, "-f", runnerOption.JenkinsfilePath}
			env := os.Environ()
			err = util.Exec(binary, jenkinsWarArgs, env, centerStartOption.SystemCallExec)
		}

		return
	},
}

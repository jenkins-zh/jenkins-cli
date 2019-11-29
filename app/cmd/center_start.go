package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"

	"github.com/spf13/cobra"
)

// CenterStartOption option for upgrade Jenkins
type CenterStartOption struct {
	RoundTripper http.RoundTripper

	Port        int
	Context     string
	SetupWizard bool

	Admin string

	Download bool
	Version  string
}

var centerStartOption CenterStartOption

func init() {
	centerCmd.AddCommand(centerStartCmd)
	centerStartCmd.Flags().IntVarP(&centerStartOption.Port, "port", "", 8080,
		i18n.T("Port of Jenkins"))
	centerStartCmd.Flags().StringVarP(&centerStartOption.Context, "context", "", "/",
		i18n.T("The address of update center site mirror"))
	centerStartCmd.Flags().BoolVarP(&centerStartOption.SetupWizard, "setup-wizard", "", true,
		i18n.T("If you want to enable update center server"))
	centerStartCmd.Flags().BoolVarP(&centerStartOption.Download, "download", "", true,
		i18n.T("If you want to enable update center server"))
	centerStartCmd.Flags().StringVarP(&centerStartOption.Version, "version", "", "lts",
		i18n.T("The address of update center site mirror"))
}

var centerStartCmd = &cobra.Command{
	Use:   "start",
	Short: i18n.T("Start Jenkins server from a cache directory"),
	Long:  i18n.T("Start Jenkins server from a cache directory"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jenkinsWar := fmt.Sprintf("/Users/mac/.jenkins-cli/cache/%s/jenkins.war", centerStartOption.Version)

		if _, fileErr := os.Stat(jenkinsWar); fileErr != nil {
			download := &CenterDownloadOption{
				Mirror:       "default",
				LTS:          true,
				Output:       jenkinsWar,
				ShowProgress: true,
				Version:      centerStartOption.Version,
			}

			if err = os.MkdirAll(strings.Replace(jenkinsWar, "jenkins.war", "", -1), os.FileMode(0755)); err != nil {
				return
			}

			if err = download.DownloadJenkins(); err != nil {
				return
			}
		}

		var binary string
		binary, err = exec.LookPath("java")
		if err == nil {
			env := os.Environ()
			env = append(env, fmt.Sprintf("JENKINS_HOME=%s/%s/web", "/Users/mac/.jenkins-cli/cache", centerStartOption.Version))
			env = append(env, fmt.Sprintf("jenkins.install.runSetupWizard=%v", centerStartOption.SetupWizard))

			jenkinsWarArgs := []string{"java"}
			jenkinsWarArgs = append(jenkinsWarArgs, "-jar", jenkinsWar)
			jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpPort=%d", centerStartOption.Port))

			fmt.Println(jenkinsWarArgs)
			err = syscall.Exec(binary, jenkinsWarArgs, env)
		}
		return
	},
}

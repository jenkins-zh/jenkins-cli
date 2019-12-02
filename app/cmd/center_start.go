package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"github.com/mitchellh/go-homedir"

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

	HTTPSEnable      bool
	HTTPSPort        int
	HTTPSCertificate string
	HTTPSPrivateKey  string

	Environments []string

	Download bool
	Version  string

	DryRun bool
}

var centerStartOption CenterStartOption

func init() {
	centerCmd.AddCommand(centerStartCmd)
	centerStartCmd.Flags().IntVarP(&centerStartOption.Port, "port", "", 8080,
		i18n.T("Port of Jenkins"))
	centerStartCmd.Flags().StringVarP(&centerStartOption.Context, "context", "", "/",
		i18n.T("Web context of Jenkins server"))
	centerStartCmd.Flags().StringArrayVarP(&centerStartOption.Environments, "env", "", nil,
		i18n.T("Environments for the Jenkins which as key-value format"))
	centerStartCmd.Flags().BoolVarP(&centerStartOption.SetupWizard, "setup-wizard", "", true,
		i18n.T("If you want to show the setup wizard at first start"))
	centerStartCmd.Flags().BoolVarP(&centerStartOption.Download, "download", "", true,
		i18n.T("If you want to download jenkins.war when it does not exist"))
	centerStartCmd.Flags().StringVarP(&centerStartOption.Version, "version", "", "2.190.3",
		i18n.T("The of version of jenkins.war"))

	centerStartCmd.Flags().BoolVarP(&centerStartOption.HTTPSEnable, "https-enable", "", false,
		i18n.T("If you want to enable https"))
	centerStartCmd.Flags().IntVarP(&centerStartOption.HTTPSPort, "https-port", "", 8083,
		i18n.T("The port of https protocol"))
	centerStartCmd.Flags().StringVarP(&centerStartOption.HTTPSCertificate, "https-cert", "", "",
		i18n.T("Certificate file path for https"))
	centerStartCmd.Flags().StringVarP(&centerStartOption.HTTPSPrivateKey, "https-private", "", "",
		i18n.T("Private key file path for https"))

	centerStartCmd.Flags().BoolVarP(&centerStartOption.DryRun, "dry-run", "", false,
		i18n.T("Don't run jenkins.war really"))
}

var centerStartCmd = &cobra.Command{
	Use:   "start",
	Short: i18n.T("Start Jenkins server from a cache directory"),
	Long:  i18n.T("Start Jenkins server from a cache directory"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		var userHome string
		if userHome, err = homedir.Dir(); err != nil {
			return
		}

		jenkinsWar := fmt.Sprintf("%s/.jenkins-cli/cache/%s/jenkins.war", userHome, centerStartOption.Version)

		if !centerStartOption.DryRun {
			if _, fileErr := os.Stat(jenkinsWar); fileErr != nil {
				download := &CenterDownloadOption{
					Mirror:       "default",
					LTS:          true,
					Output:       jenkinsWar,
					ShowProgress: true,
					Version:      centerStartOption.Version,
				}

				if err = download.DownloadJenkins(); err != nil {
					return
				}
			}
		}

		var binary string

		if !centerStartOption.DryRun {
			binary, err = exec.LookPath("java")
		}

		if err == nil {
			env := os.Environ()
			env = append(env, fmt.Sprintf("JENKINS_HOME=%s/.jenkins-cli/cache/%s/web", userHome, centerStartOption.Version))
			env = append(env, fmt.Sprintf("jenkins.install.runSetupWizard=%v", centerStartOption.SetupWizard))

			if centerStartOption.Environments != nil {
				for _, item := range centerStartOption.Environments {
					env = append(env, item)
				}
			}

			jenkinsWarArgs := []string{"java"}
			jenkinsWarArgs = append(jenkinsWarArgs, "-jar", jenkinsWar)
			jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpPort=%d", centerStartOption.Port))
			jenkinsWarArgs = append(jenkinsWarArgs, "--argumentsRealm.passwd.admin=admin --argumentsRealm.roles.admin=admin")
			jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--prefix=%s", centerStartOption.Context))

			if centerStartOption.HTTPSEnable {
				jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpsPort=%d", centerStartOption.HTTPSPort))
				jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpsCertificate=%s", centerStartOption.HTTPSCertificate),
					fmt.Sprintf("--httpsPrivateKey=%s", centerStartOption.HTTPSPrivateKey))
			}

			if !centerStartOption.DryRun {
				err = syscall.Exec(binary, jenkinsWarArgs, env)
			}
		}
		return
	},
}

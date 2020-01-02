package cmd

import (
	"fmt"
	"os"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

// CenterStartOption option for upgrade Jenkins
type CenterStartOption struct {
	CommonOption

	Port                      int
	Context                   string
	SetupWizard               bool
	AdminCanGenerateNewTokens bool

	// comes from folder plugin
	ConcurrentIndexing int

	Admin string

	HTTPSEnable      bool
	HTTPSPort        int
	HTTPSCertificate string
	HTTPSPrivateKey  string

	Environments []string
	System       []string

	RandomWebDir bool
	Download     bool
	Version      string

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
	centerStartCmd.Flags().StringArrayVarP(&centerStartOption.System, "sys", "", nil,
		i18n.T("System property key-value"))
	centerStartCmd.Flags().BoolVarP(&centerStartOption.SetupWizard, "setup-wizard", "", true,
		i18n.T("If you want to show the setup wizard at first start"))
	centerStartCmd.Flags().BoolVarP(&centerStartOption.AdminCanGenerateNewTokens, "admin-can-generate-new-tokens", "", false,
		i18n.T("If enabled, the users with administer permissions can generate new tokens for other users"))

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

	centerStartCmd.Flags().IntVarP(&centerStartOption.ConcurrentIndexing, "concurrent-indexing", "", -1,
		i18n.T("Concurrent indexing limit, take this value only it is bigger than -1"))

	centerStartCmd.Flags().BoolVarP(&centerStartOption.RandomWebDir, "random-web-dir", "", false,
		i18n.T("If start jenkins.war in a random web dir"))
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
		binary, err = util.LookPath("java", centerStartOption.LookPathContext)
		if err == nil {
			env := os.Environ()

			if centerStartOption.RandomWebDir {
				env = append(env, fmt.Sprintf("JENKINS_HOME=%s/.jenkins-cli/cache/%s/web", os.TempDir(), centerStartOption.Version))
			} else {
				env = append(env, fmt.Sprintf("JENKINS_HOME=%s/.jenkins-cli/cache/%s/web", userHome, centerStartOption.Version))
			}

			if centerStartOption.Environments != nil {
				for _, item := range centerStartOption.Environments {
					env = append(env, item)
				}
			}

			jenkinsWarArgs := []string{"java"}
			jenkinsWarArgs = centerStartOption.setSystemProperty(jenkinsWarArgs)
			jenkinsWarArgs = append(jenkinsWarArgs, "-jar", jenkinsWar)
			jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpPort=%d", centerStartOption.Port))
			jenkinsWarArgs = append(jenkinsWarArgs, "--argumentsRealm.passwd.admin=admin")
			jenkinsWarArgs = append(jenkinsWarArgs, "--argumentsRealm.roles.admin=admin")
			jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--prefix=%s", centerStartOption.Context))

			if centerStartOption.HTTPSEnable {
				jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpsPort=%d", centerStartOption.HTTPSPort))
				jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpsCertificate=%s", centerStartOption.HTTPSCertificate),
					fmt.Sprintf("--httpsPrivateKey=%s", centerStartOption.HTTPSPrivateKey))
			}

			err = util.Exec(binary, jenkinsWarArgs, env, centerStartOption.SystemCallExec)
		}
		return
	},
}

func (c *CenterStartOption) setSystemProperty(jenkinsWarArgs []string) []string {
	if c.System == nil {
		c.System = []string{}
	}

	c.System = append(c.System, fmt.Sprintf("jenkins.install.runSetupWizard=%v", c.SetupWizard))
	c.System = append(c.System, fmt.Sprintf("jenkins.security.ApiTokenProperty.adminCanGenerateNewTokens=%v", c.AdminCanGenerateNewTokens))
	if c.ConcurrentIndexing > -1 {
		c.System = append(c.System, fmt.Sprintf("com.cloudbees.hudson.plugins.folder.computed.ThrottleComputationQueueTaskDispatcher.LIMIT=%d", c.ConcurrentIndexing))
	}

	for _, item := range centerStartOption.System {
		jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("-D%s", item))
	}
	return jenkinsWarArgs
}

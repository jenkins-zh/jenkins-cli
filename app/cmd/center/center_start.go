package center

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"go.uber.org/zap"
	"os"
	"path/filepath"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

func NewCenterStartcmd(client common.JenkinsClient) (cmd *cobra.Command) {
	centerStartOption := CenterStartOption{}

	cmd = &cobra.Command{
		Use:   "start",
		Short: i18n.T("Start Jenkins server from a cache directory"),
		Long:  i18n.T("Start Jenkins server from a cache directory"),
		RunE:  centerStartOption.RunE,
	}

	cmd.Flags().IntVarP(&centerStartOption.Port, "port", "", 8080,
		i18n.T("Port of Jenkins"))
	cmd.Flags().StringVarP(&centerStartOption.Context, "context", "", "/",
		i18n.T("Web context of Jenkins server"))
	cmd.Flags().StringArrayVarP(&centerStartOption.Environments, "env", "", nil,
		i18n.T("Environments for the Jenkins which as key-value format"))
	cmd.Flags().StringArrayVarP(&centerStartOption.System, "sys", "", nil,
		i18n.T("System property key-value"))
	cmd.Flags().BoolVarP(&centerStartOption.SetupWizard, "setup-wizard", "", true,
		i18n.T("If you want to show the setup wizard at first start"))
	cmd.Flags().BoolVarP(&centerStartOption.AdminCanGenerateNewTokens, "admin-can-generate-new-tokens", "", false,
		i18n.T("If enabled, the users with administer permissions can generate new tokens for other users"))
	cmd.Flags().BoolVarP(&centerStartOption.ShowProperties, "show-properties", "", false,
		"Show all the Jenkins properties and exit")

	cmd.Flags().BoolVarP(&centerStartOption.Download, "download", "", true,
		i18n.T("If you want to download jenkins.war when it does not exist"))
	cmd.Flags().StringVarP(&centerStartOption.Version, "version", "", "2.190.3",
		i18n.T("The of version of jenkins.war"))
	cmd.Flags().BoolVarP(&centerStartOption.LTS, "lts", "", true,
		i18n.T("If you want to download Jenkins as LTS"))
	cmd.Flags().StringVarP(&centerStartOption.Formula, "formula", "", "",
		i18n.T("The formula of jenkins.war, only support zh currently"))

	cmd.Flags().BoolVarP(&centerStartOption.HTTPSEnable, "https-enable", "", false,
		i18n.T("If you want to enable https"))
	cmd.Flags().IntVarP(&centerStartOption.HTTPSPort, "https-port", "", 8083,
		i18n.T("The port of https protocol"))
	cmd.Flags().StringVarP(&centerStartOption.HTTPSCertificate, "https-cert", "", "",
		i18n.T("Certificate file path for https"))
	cmd.Flags().StringVarP(&centerStartOption.HTTPSPrivateKey, "https-private", "", "",
		i18n.T("Private key file path for https"))

	cmd.Flags().IntVarP(&centerStartOption.ConcurrentIndexing, "concurrent-indexing", "", -1,
		i18n.T("Concurrent indexing limit, take this value only it is bigger than -1"))

	cmd.Flags().BoolVarP(&centerStartOption.RandomWebDir, "random-web-dir", "", false,
		i18n.T("If start jenkins.war in a random web dir"))
	cmd.Flags().BoolVarP(&centerStartOption.DryRun, "dry-run", "", false,
		i18n.T("Don't run jenkins.war really"))

	err := cmd.RegisterFlagCompletionFunc("version", func(cmd *cobra.Command, args []string, toComplete string) (strings []string, directive cobra.ShellCompDirective) {
		var userHome string
		var err error
		if userHome, err = homedir.Dir(); err != nil {
			return
		}

		var machedPathes []string
		jenkinsWar := fmt.Sprintf("%s/.jenkins-cli/cache/*/jenkins.war", userHome)
		if machedPathes, err = filepath.Glob(jenkinsWar); err != nil {
			return
		}

		versionArray := make([]string, len(machedPathes))
		for _, path := range machedPathes {
			versionArray = append(versionArray, filepath.Base(filepath.Dir(path)))
		}

		return versionArray, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		cmd.PrintErrf("register flag version failed %#v\n", err)
	}
	return
}

func (o *CenterStartOption) RunE(cmd *cobra.Command, _ []string) (err error) {
	if o.ShowProperties {
		for _, item := range o.getAllJenkinsSystemProperties() {
			cmd.Println(item.Key)
		}
		return
	}

	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	jenkinsWar := fmt.Sprintf("%s/.jenkins-cli/cache/%s/jenkins.war", userHome, o.Version)

	o.Logger.Info("prepare to download jenkins.war", zap.String("localPath", jenkinsWar))

	if !o.DryRun {
		if _, fileErr := os.Stat(jenkinsWar); fileErr != nil {
			download := &CenterDownloadOption{
				Mirror:       "default",
				Formula:      o.Formula,
				LTS:          o.LTS,
				Output:       jenkinsWar,
				ShowProgress: true,
				Version:      o.Version,
			}

			if err = download.DownloadJenkins(); err != nil {
				return
			}
		}
	}

	var binary string
	binary, err = util.LookPath("java", o.LookPathContext)
	if err == nil {
		env := os.Environ()

		if o.RandomWebDir {
			randomWebDir := fmt.Sprintf("JENKINS_HOME=%s/.jenkins-cli/cache/%s/web", os.TempDir(), o.Version)
			defer func(logger helper.Printer, randomWebDir string) {
				if err := os.RemoveAll(randomWebDir); err != nil {
					logger.PrintErr(fmt.Sprintf("remove random web dir [%s] of Jenkins failed, %#v", randomWebDir, err))
				}
			}(cmd, randomWebDir)

			env = append(env, randomWebDir)
		} else {
			env = append(env, fmt.Sprintf("JENKINS_HOME=%s/.jenkins-cli/cache/%s/web", userHome, o.Version))
		}

		if o.Environments != nil {
			for _, item := range o.Environments {
				env = append(env, item)
			}
		}

		jenkinsWarArgs := []string{"java"}
		jenkinsWarArgs = o.setSystemProperty(jenkinsWarArgs)
		jenkinsWarArgs = append(jenkinsWarArgs, "-jar", jenkinsWar)
		jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpPort=%d", o.Port))
		jenkinsWarArgs = append(jenkinsWarArgs, "--argumentsRealm.passwd.admin=admin")
		jenkinsWarArgs = append(jenkinsWarArgs, "--argumentsRealm.roles.admin=admin")
		jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--prefix=%s", o.Context))

		if o.HTTPSEnable {
			jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpsPort=%d", o.HTTPSPort))
			jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("--httpsCertificate=%s", o.HTTPSCertificate),
				fmt.Sprintf("--httpsPrivateKey=%s", o.HTTPSPrivateKey))
		}

		err = util.Exec(binary, jenkinsWarArgs, env, o.SystemCallExec)
	}
	return
}

type JenkinsSystemProperty struct {
	Key          string
	DefaultValue string
	Description  string
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

	for _, item := range c.System {
		jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("-D%s", item))
	}
	return jenkinsWarArgs
}

func (c *CenterStartOption) getAllJenkinsSystemProperties() (properties []JenkinsSystemProperty) {
	properties = []JenkinsSystemProperty{{
		Key: "jenkins.install.runSetupWizard",
	}, {
		Key: "jenkins.security.ApiTokenProperty.adminCanGenerateNewTokens",
	}, {
		Key: "com.cloudbees.hudson.plugins.folder.computed.ThrottleComputationQueueTaskDispatcher.LIMIT",
	}, {
		Key: "hudson.model.DownloadService.noSignatureCheck",
	}, {
		Key: "hudson.model.DirectoryBrowserSupport.CSP",
	}, {
		Key: "hudson.security.csrf.DefaultCrumbIssuer.EXCLUDE_SESSION_ID",
	}, {
		Key: "kubernetes.websocket.ping.interval",
	}, {
		Key: "org.jenkinsci.plugins.gitclient.Git.timeOut",
	}}
	return
}

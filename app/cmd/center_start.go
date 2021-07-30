package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/helper"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

// CenterStartOption option for upgrade Jenkins
type CenterStartOption struct {
	common.Option

	Port                      int
	AgentPort                 int
	Context                   string
	SetupWizard               bool
	AdminCanGenerateNewTokens bool
	CleanHome                 bool
	CrumbExcludeSessionID     bool

	// comes from folder plugin
	ConcurrentIndexing int

	Admin string

	HTTPSEnable      bool
	HTTPSPort        int
	HTTPSCertificate string
	HTTPSPrivateKey  string

	Environments []string
	System       []string

	Download     bool
	Mirror       string
	Thread       int
	Version      string
	LTS          bool
	Formula      string
	RandomWebDir bool

	Mode          string
	Image         string
	ForcePull     bool
	ContainerUser string
	DryRun        bool
}

var centerStartOption CenterStartOption

func init() {
	jenkinsVersion := util.GetEnvOrDefault("JCLI_JENKINS_VERSION", "2.249.1")

	centerCmd.AddCommand(centerStartCmd)
	flags := centerStartCmd.Flags()
	flags.IntVarP(&centerStartOption.Port, "port", "", 8080,
		i18n.T("Port of Jenkins"))
	flags.IntVarP(&centerStartOption.AgentPort, "agent-port", "", 50000,
		i18n.T("Port of Jenkins agent"))
	flags.StringVarP(&centerStartOption.Context, "context", "", "/",
		i18n.T("Web context of Jenkins server"))
	flags.StringArrayVarP(&centerStartOption.Environments, "env", "", nil,
		i18n.T("Environments for the Jenkins which as key-value format"))
	flags.StringArrayVarP(&centerStartOption.System, "sys", "", nil,
		i18n.T("System property key-value"))
	flags.BoolVarP(&centerStartOption.SetupWizard, "setup-wizard", "", true,
		i18n.T("If you want to show the setup wizard at first start"))
	flags.BoolVarP(&centerStartOption.CrumbExcludeSessionID, "crumb-exclude-sessionId", "", false,
		i18n.T(`Add system properties with 'hudson.security.csrf.DefaultCrumbIssuer.EXCLUDE_SESSION_ID=true'`))
	flags.BoolVarP(&centerStartOption.AdminCanGenerateNewTokens, "admin-can-generate-new-tokens", "", false,
		i18n.T("If enabled, the users with administer permissions can generate new tokens for other users"))
	flags.BoolVarP(&centerStartOption.CleanHome, "clean-home", "", false,
		i18n.T("If you want to clean the JENKINS_HOME before start it"))

	flags.BoolVarP(&centerStartOption.Download, "download", "", true,
		i18n.T("If you want to download jenkins.war when it does not exist"))
	flags.StringVarP(&centerStartOption.Mirror, "mirror", "", "default",
		i18n.T("The mirror site of Jenkins"))
	flags.IntVarP(&centerStartOption.Thread, "thread", "t", 0, "Using multi-thread to download jenkins.war")
	flags.StringVarP(&centerStartOption.Version, "version", "", jenkinsVersion,
		i18n.T("The of version of jenkins.war. You can give it another default value by setting env JCLI_JENKINS_VERSION"))
	flags.BoolVarP(&centerStartOption.LTS, "lts", "", true,
		i18n.T("If you want to download Jenkins as LTS"))
	flags.StringVarP(&centerStartOption.Formula, "formula", "", "",
		i18n.T("The formula of jenkins.war, only support zh currently"))

	flags.BoolVarP(&centerStartOption.HTTPSEnable, "https-enable", "", false,
		i18n.T("If you want to enable https"))
	flags.IntVarP(&centerStartOption.HTTPSPort, "https-port", "", 8083,
		i18n.T("The port of https protocol"))
	flags.StringVarP(&centerStartOption.HTTPSCertificate, "https-cert", "", "",
		i18n.T("Certificate file path for https"))
	flags.StringVarP(&centerStartOption.HTTPSPrivateKey, "https-private", "", "",
		i18n.T("Private key file path for https"))

	flags.IntVarP(&centerStartOption.ConcurrentIndexing, "concurrent-indexing", "", -1,
		i18n.T("Concurrent indexing limit, take this value only it is bigger than -1"))

	flags.StringVarP(&centerStartOption.Mode, "mode", "m", "java",
		i18n.T("Which mode do you want to run. Supported mode contains: java, docker"))
	flags.StringVarP(&centerStartOption.Image, "image", "", "jenkins/jenkins",
		i18n.T("Which docker image do you want to run. It works only the mode is docker. Please use --version if you want to specific the version of docker image."))
	flags.BoolVarP(&centerStartOption.ForcePull, "force-pull", "", false,
		"Indicate if you want to force pull image. Sometimes your local image is not the latest.")
	flags.StringVarP(&centerStartOption.ContainerUser, "c-user", "", "",
		i18n.T("Container Username or UID (format: <name|uid>[:<group|gid>])"))

	flags.BoolVarP(&centerStartOption.RandomWebDir, "random-web-dir", "", false,
		i18n.T(fmt.Sprintf("If start jenkins.war in a random web dir which under %s/.jenkins-cli/cache", os.TempDir())))
	flags.BoolVarP(&centerStartOption.DryRun, "dry-run", "", false,
		i18n.T("Don't run jenkins.war really"))

	err := centerStartCmd.RegisterFlagCompletionFunc("version", func(cmd *cobra.Command, args []string, toComplete string) (strings []string, directive cobra.ShellCompDirective) {
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
		centerCmd.PrintErrf("register flag version failed %#v\n", err)
	}
}

var centerStartCmd = &cobra.Command{
	Use:               "start",
	Short:             i18n.T("Start Jenkins server from a cache directory"),
	Long:              i18n.T("Start Jenkins server from a cache directory"),
	RunE:              centerStartOption.run,
	ValidArgsFunction: common.NoFileCompletion,
}

func (c *CenterStartOption) run(cmd *cobra.Command, _ []string) (err error) {
	switch c.Mode {
	case "java":
		err = c.createJavaArgs(cmd)
	case "docker":
		if c.ForcePull {
			if err = c.pullImage(fmt.Sprintf("%s:%s", c.Image, c.Version)); err != nil {
				return
			}
		}
		err = c.createDockerArgs(cmd)
	}
	return
}

func (c *CenterStartOption) pullImage(image string) (err error) {
	var binary string
	binary, err = util.LookPath("docker", c.LookPathContext)
	if err == nil {
		env := os.Environ()

		dockerArgs := []string{"docker", "pull", image}
		err = util.Exec(binary, dockerArgs, env, centerStartOption.SystemCallExec)
	}
	return
}

func (c *CenterStartOption) createDockerArgs(cmd *cobra.Command) (err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	var binary string
	binary, err = util.LookPath("docker", c.LookPathContext)
	if err == nil {
		env := os.Environ()

		jenkinsHome := fmt.Sprintf("%s/.jenkins-cli/cache/%s/web", userHome, strings.Split(c.Version, "@")[0])

		dockerArgs := []string{"docker", "run"}
		dockerArgs = append(dockerArgs, "-v", fmt.Sprintf("%s:/var/jenkins_home", jenkinsHome))
		dockerArgs = append(dockerArgs, "-p", fmt.Sprintf("%d:8080", c.Port))
		dockerArgs = append(dockerArgs, "-p", fmt.Sprintf("%d:50000", c.AgentPort))

		if c.ContainerUser != "" {
			dockerArgs = append(dockerArgs, "-u", c.ContainerUser)
		}

		javaOpts := ""
		args := make([]string, 0)
		args = c.setSystemProperty(args)
		for _, arg := range args {
			javaOpts += " " + arg
		}
		if strings.TrimSpace(javaOpts) != "" {
			dockerArgs = append(dockerArgs, "-e", fmt.Sprintf("JAVA_OPTS=%s", strings.TrimSpace(javaOpts)))
		}

		dockerArgs = append(dockerArgs, fmt.Sprintf("%s:%s", c.Image, c.Version))
		dockerArgs = append(dockerArgs, fmt.Sprintf("--prefix=%s", centerStartOption.Context))

		if c.CleanHome {
			logger.Debug("start to clean JENKINS_HOME before start it", zap.String("JENKINS_HOME", jenkinsHome))
			if err = os.RemoveAll(jenkinsHome); err != nil {
				err = fmt.Errorf("cannot remove JENKINS_HOME %s, error: %v", jenkinsHome, err)
				return
			}
		}
		err = util.Exec(binary, dockerArgs, env, centerStartOption.SystemCallExec)
	}
	return
}

func (c *CenterStartOption) createJavaArgs(cmd *cobra.Command) (err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}
	jenkinsWar := fmt.Sprintf("%s/.jenkins-cli/cache/%s/jenkins.war", userHome, centerStartOption.Version)

	logger.Info("prepare to download jenkins.war", zap.String("localPath", jenkinsWar))

	if !centerStartOption.DryRun {
		if _, fileErr := os.Stat(jenkinsWar); fileErr != nil {
			download := &CenterDownloadOption{
				Mirror:       c.Mirror,
				Formula:      centerStartOption.Formula,
				LTS:          centerStartOption.LTS,
				Output:       jenkinsWar,
				ShowProgress: true,
				Version:      centerStartOption.Version,
				Thread:       centerStartOption.Thread,
			}

			if err = download.DownloadJenkins(); err != nil {
				return
			}
		}
	}

	var binary string
	var jenkinsHome string
	binary, err = util.LookPath("java", centerStartOption.LookPathContext)
	if err == nil {
		env := os.Environ()

		if centerStartOption.RandomWebDir {
			rand.Seed(1)
			jenkinsHome = fmt.Sprintf("%s/.jenkins-cli/cache/%d/%s/web", os.TempDir(), rand.Int(), centerStartOption.Version)
			randomWebDir := fmt.Sprintf("JENKINS_HOME=%s", jenkinsHome)
			defer func(logger helper.Printer, randomWebDir string) {
				if err := os.RemoveAll(randomWebDir); err != nil {
					logger.PrintErr(fmt.Sprintf("remove random web dir [%s] of Jenkins failed, %#v", randomWebDir, err))
				}
			}(cmd, randomWebDir)

			env = append(env, randomWebDir)
		} else {
			jenkinsHome = fmt.Sprintf("%s/.jenkins-cli/cache/%s/web", userHome, centerStartOption.Version)
			env = append(env, fmt.Sprintf("JENKINS_HOME=%s", jenkinsHome))
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

		if c.CleanHome {
			logger.Debug("start to clean JENKINS_HOME before start it", zap.String("JENKINS_HOME", jenkinsHome))
			if err = os.RemoveAll(jenkinsHome); err != nil {
				err = fmt.Errorf("cannot remove JENKINS_HOME %s, error: %v", jenkinsHome, err)
				return
			}
		}
		err = util.Exec(binary, jenkinsWarArgs, env, centerStartOption.SystemCallExec)
	}
	return
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

	if c.CrumbExcludeSessionID {
		c.System = append(c.System, "hudson.security.csrf.DefaultCrumbIssuer.EXCLUDE_SESSION_ID=true")
	}

	for _, item := range centerStartOption.System {
		jenkinsWarArgs = append(jenkinsWarArgs, fmt.Sprintf("-D%s", item))
	}
	return jenkinsWarArgs
}

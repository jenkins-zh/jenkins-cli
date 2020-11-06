package cmd

import (
	"encoding/xml"
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"html/template"
	"io/ioutil"
	"os"
	"path"
)

func init() {
	rootCmd.AddCommand(cwpCmd)

	cwpCmd.Flags().BoolVarP(&cwpOptions.BatchMode, "batch-mode", "", false,
		i18n.T("Enables the batch mode for the build"))
	cwpCmd.Flags().StringVarP(&cwpOptions.BomPath, "bom-path", "", "",
		i18n.T("Path to the BOM file. If defined, it will override settings in Config YAML"))
	cwpCmd.Flags().StringVarP(&cwpOptions.Environment, "environment", "", "",
		i18n.T("Environment to be used"))
	cwpCmd.Flags().BoolVarP(&cwpOptions.InstallArtifacts, "install-artifacts", "", false,
		i18n.T("If set, the final artifacts will be automatically installed to the local repository (current version - only WAR)"))
	cwpCmd.Flags().StringVarP(&cwpOptions.ConfigPath, "config-path", "", "",
		i18n.T("Path to the configuration YAML. See the tool's README for format"))
	cwpCmd.Flags().BoolVarP(&cwpOptions.Demo, "demo", "", false,
		i18n.T("Enables demo mode with predefined config file"))
	cwpCmd.Flags().StringVarP(&cwpOptions.MvnSettingsFile, "mvn-settings-file", "", "",
		i18n.T("Path to a custom Maven settings file to be used within the build"))
	cwpCmd.Flags().StringVarP(&cwpOptions.TmpDir, "tmp-dir", "", "",
		i18n.T("Temporary directory for generated files and the output WAR."))
	cwpCmd.Flags().StringVarP(&cwpOptions.Version, "version", "", "1.0-SNAPSHOT",
		i18n.T("Version of WAR to be set."))

	cwpCmd.Flags().BoolVarP(&cwpOptions.ShowProgress, "show-progress", "", true,
		i18n.T("Show the progress of downloading files"))
	cwpCmd.Flags().StringVarP(&cwpOptions.MetadataURL, "metadata-url", "",
		"https://repo.jenkins-ci.org/list/releases/io/jenkins/tools/custom-war-packager/custom-war-packager-cli/maven-metadata.xml",
		i18n.T("The metadata URL"))
	cwpCmd.Flags().StringToStringVarP(&cwpOptions.ValueSet, "value-set", "", nil,
		`The value set of config template`)

	localCache := path.Join(os.TempDir(), "/", ".jenkins-cli")
	if userHome, err := homedir.Dir(); err == nil {
		localCache = path.Join(userHome, "/", ".jenkins-cli")
	}
	cwpCmd.Flags().StringVarP(&cwpOptions.LocalCache, "local-cache", "", localCache,
		i18n.T("The local cache directory"))
}

// CWPOptions is the option of custom-war-packager
// see also https://github.com/jenkinsci/custom-war-packager
type CWPOptions struct {
	common.CommonOption

	ConfigPath      string
	Version         string
	TmpDir          string
	Environment     string
	BomPath         string
	MvnSettingsFile string

	BatchMode        bool
	Demo             bool
	InstallArtifacts bool

	ShowProgress bool
	MetadataURL  string
	LocalCache   string

	ValueSet map[string]string
}

var cwpOptions CWPOptions

var cwpCmd = &cobra.Command{
	Use:   "cwp",
	Short: i18n.T("Custom Jenkins WAR packager for Jenkins"),
	Long: i18n.T(`Custom Jenkins WAR packager for Jenkins
This's a wrapper of https://github.com/jenkinsci/custom-war-packager`),
	RunE:    cwpOptions.Run,
	Example: `jcli cwp --config-path test.yaml`,
	Annotations: map[string]string{
		common.Since: "v0.0.27",
	},
}

// Run is the main logic of cwp cmd
func (o *CWPOptions) Run(cmd *cobra.Command, args []string) (err error) {
	localCWP := o.getLocalCWP()
	_, err = os.Stat(localCWP)
	if os.IsNotExist(err) {
		if err = o.Download(); err != nil {
			return
		}
	} else if err != nil {
		return
	}

	var binary string
	binary, err = util.LookPath("java", o.LookPathContext)
	if err == nil {
		env := os.Environ()

		cwpArgs := []string{"java"}
		cwpArgs = append(cwpArgs, "-jar", localCWP)

		if o.Demo {
			cwpArgs = append(cwpArgs, "-demo")
		}

		if o.BatchMode {
			cwpArgs = append(cwpArgs, "--batch-mode")
		}

		if o.InstallArtifacts {
			cwpArgs = append(cwpArgs, "--installArtifacts")
		}

		if o.ConfigPath != "" {
			configPath := o.ConfigPath
			if configPath, err = RenderTemplate(o.ConfigPath, o.ValueSet); err != nil {
				return
			}
			defer os.RemoveAll(configPath)

			cwpArgs = append(cwpArgs, "-configPath", configPath)
		}

		if o.TmpDir != "" {
			cwpArgs = append(cwpArgs, "-tmpDir", o.TmpDir)
		}

		if o.Version != "" {
			cwpArgs = append(cwpArgs, "-version", o.Version)
		}
		err = util.Exec(binary, cwpArgs, env, o.SystemCallExec)
	}
	return
}

// RenderTemplate render a go template to a temporary file
func RenderTemplate(filepath string, values map[string]string) (result string, err error) {
	var t *template.Template
	tmp := template.New(path.Base(filepath)).Funcs(template.FuncMap{
		"default": func(arg interface{}, value interface{}) interface{} {
			if value == nil {
				return arg
			}
			return value
		},
	})

	if t, err = tmp.ParseFiles(filepath); err == nil {
		logger.Debug("parse template done", zap.String("path", filepath))
		f, _ := ioutil.TempFile("/tmp", ".yaml")
		err = t.Execute(f, values)

		result = f.Name()
	} else {
		logger.Error("error when parsing template file", zap.Error(err))
	}
	return
}

// Download get the latest cwp from server into local
func (o *CWPOptions) Download() (err error) {
	var latest string
	if latest, err = o.GetLatest(); err == nil {
		cwpURL := o.GetCWPURL(latest)

		err = o.downloadFile(cwpURL, o.getLocalCWP())
	}
	return
}

// GetCWPURL returns the download URL of a specific version cwp
func (o *CWPOptions) GetCWPURL(version string) string {
	return fmt.Sprintf("https://repo.jenkins-ci.org/list/releases/io/jenkins/tools/custom-war-packager/custom-war-packager-cli/%s/custom-war-packager-cli-%s-jar-with-dependencies.jar",
		version, version)
}

func (o *CWPOptions) getLocalCWP() string {
	return path.Join(o.LocalCache, "cwp-cli.jar")
}

// GetLatest returns the latest of cwp
func (o *CWPOptions) GetLatest() (version string, err error) {
	metadataURL := o.MetadataURL
	output := "metadata.xml"

	if err = o.downloadFile(metadataURL, output); err == nil {
		var data []byte

		mavenMeta := MavenMetadata{}
		if data, err = ioutil.ReadFile(output); err == nil {
			err = xml.Unmarshal(data, &mavenMeta)
		}

		if err == nil {
			version = mavenMeta.Versioning.Latest
		}
	}
	return
}

func (o *CWPOptions) downloadFile(url, output string) (err error) {
	downloader := util.HTTPDownloader{
		RoundTripper:   o.RoundTripper,
		TargetFilePath: output,
		URL:            url,
		ShowProgress:   o.ShowProgress,
	}
	err = downloader.DownloadFile()
	return
}

// MavenMetadata is the maven metadata xml root
type MavenMetadata struct {
	XMLName    xml.Name        `xml:"metadata"`
	Versioning MavenVersioning `xml:"versioning"`
}

// MavenVersioning is the versioning of maven
type MavenVersioning struct {
	XMLName xml.Name `xml:"versioning"`
	Latest  string   `xml:"latest"`
	Release string   `xml:"release"`
}

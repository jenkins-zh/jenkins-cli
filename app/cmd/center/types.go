package center

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"net/http"
)

// CenterOption is the center cmd option
type CenterOption struct {
	common.WatchOption
	JenkinsClient    common.JenkinsClient
	JenkinsConfigMgr common.JenkinsConfigMgr

	RoundTripper http.RoundTripper
	CenterStatus string
}

// CenterStartOption option for upgrade Jenkins
type CenterStartOption struct {
	common.CommonOption

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

	Environments   []string
	System         []string
	ShowProperties bool

	Download     bool
	Version      string
	LTS          bool
	Formula      string
	RandomWebDir bool

	DryRun bool
}

// CenterUpgradeOption option for upgrade Jenkins
type CenterUpgradeOption struct {
	RoundTripper  http.RoundTripper
	JenkinsClient common.JenkinsClient
}

// CenterWatchOption as the options of watch command
type CenterWatchOption struct {
	common.WatchOption
	*CenterOption
	JenkinsClient       common.JenkinsClient
	UtilNeedRestart     bool
	UtilInstallComplete bool

	RoundTripper  http.RoundTripper
	CeneterStatus string
}

// CenterMirrorOption option for upgrade Jenkins
type CenterMirrorOption struct {
	RoundTripper  http.RoundTripper
	JenkinsClient common.JenkinsClient

	Enable    bool
	MirrorURL string
}

// CenterIdentityOption option for upgrade Jenkins
type CenterIdentityOption struct {
	common.CommonOption
	JenkinsClient common.JenkinsClient
}

// CenterDownloadOption as the options of download command
type CenterDownloadOption struct {
	common.CommonOption
	JenkinsClient    common.JenkinsClient
	JenkinsConfigMgr common.JenkinsConfigMgr
	LTS              bool
	Mirror           string
	Version          string

	Output       string
	ShowProgress bool

	Formula string

	RoundTripper http.RoundTripper
}

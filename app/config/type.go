package config

const (
	// ANNOTATION_CONFIG_LOAD annotation for config loading set
	ANNOTATION_CONFIG_LOAD string = "config.load"
)

// JenkinsServer holds the configuration of your Jenkins
type JenkinsServer struct {
	Name               string `yaml:"name"`
	URL                string `yaml:"url"`
	UserName           string `yaml:"username"`
	Token              string `yaml:"token"`
	Proxy              string `yaml:"proxy"`
	ProxyAuth          string `yaml:"proxyAuth"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
	Description        string `yaml:"description"`
}

// CommandHook is a hook
type CommandHook struct {
	Path    string `yaml:"path"`
	Command string `yaml:"cmd"`
}

// PluginSuite define a suite of plugins
type PluginSuite struct {
	Name        string   `yaml:"name"`
	Plugins     []string `yaml:"plugins"`
	Description string   `yaml:"description"`
}

// JenkinsMirror represents the mirror of Jenkins
type JenkinsMirror struct {
	Name string
	URL  string
}

// Config is a global config struct
type Config struct {
	Current        string          `yaml:"current"`
	Language       string          `yaml:"language"`
	JenkinsServers []JenkinsServer `yaml:"jenkins_servers"`
	PreHooks       []CommandHook   `yaml:"preHooks"`
	PostHooks      []CommandHook   `yaml:"postHooks"`
	PluginSuites   []PluginSuite   `yaml:"pluginSuites"`
	Mirrors        []JenkinsMirror `yaml:"mirrors"`
}

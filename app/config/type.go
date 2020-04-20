package config

const (
	// ANNOTATION_CONFIG_LOAD annotation for config loading set
	ANNOTATION_CONFIG_LOAD string = "config.load"
)

// JenkinsServer holds the configuration of your Jenkins
type JenkinsServer struct {
	Name               string            `yaml:"name"`
	URL                string            `yaml:"url"`
	UserName           string            `yaml:"username"`
	Token              string            `yaml:"token"`
	Proxy              string            `yaml:"proxy,omitempty"`
	ProxyAuth          string            `yaml:"proxyAuth,omitempty"`
	InsecureSkipVerify bool              `yaml:"insecureSkipVerify,omitempty"`
	Description        string            `yaml:"description,omitempty"`
	Data               map[string]string `yaml:"data,omitempty"`
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
	Language       string          `yaml:"language,omitempty"`
	JenkinsServers []JenkinsServer `yaml:"jenkins_servers"`
	PreHooks       []CommandHook   `yaml:"preHooks,omitempty"`
	PostHooks      []CommandHook   `yaml:"postHooks,omitempty"`
	PluginSuites   []PluginSuite   `yaml:"pluginSuites,omitempty"`
	Mirrors        []JenkinsMirror `yaml:"mirrors"`
}

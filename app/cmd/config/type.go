package config

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"io"
	"net/http"
)

type (
	configPluginListCmd struct {
		*common.Option
		common.OutputOption
	}
	jcliPluginFetchCmd struct {
		*common.Option
		PluginRepo string
		Reset      bool

		Username   string
		Password   string
		SSHKeyFile string

		output io.Writer
	}
	jcliPluginInstallCmd struct {
		*common.Option
		RoundTripper http.RoundTripper
		ShowProgress bool

		output io.Writer
	}
	jcliPluginUninstallCmd struct {
		*common.Option
	}
	jcliPluginUpdateCmd struct {
		*common.Option
	}
	plugin struct {
		Use          string `yaml:"use"`
		Short        string
		Long         string
		Main         string
		Version      string
		DownloadLink string `yaml:"downloadLink"`
		Installed    bool
	}
	pluginError struct {
		error
		code int
	}
)

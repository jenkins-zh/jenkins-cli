package config_plugin

import (
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"io"
	"net/http"
)

type (
	configPluginListCmd struct {
		*common.CommonOption
		common.OutputOption
	}
	jcliPluginFetchCmd struct {
		*common.CommonOption
		PluginRepo string
		Reset      bool

		Username   string
		Password   string
		SSHKeyFile string

		output io.Writer
	}
	jcliPluginInstallCmd struct {
		*common.CommonOption
		RoundTripper http.RoundTripper
		ShowProgress bool

		output io.Writer
	}
	jcliPluginUninstallCmd struct {
		*common.CommonOption
	}
	jcliPluginUpdateCmd struct {
		*common.CommonOption
	}
	plugin struct {
		Use          string
		Short        string
		Long         string
		Main         string
		Version      string
		DownloadLink string `yaml:"downloadLink"`
	}
	pluginError struct {
		error
		code int
	}
)

package client

import (
	"github.com/jenkins-zh/jenkins-cli/util"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
	"go.uber.org/zap"
	"net/http"
)

var logger *zap.Logger

// SetLogger set a global logger
func SetLogger(zapLogger *zap.Logger) {
	logger = zapLogger
}

func init() {
	if logger == nil {
		var err error
		if logger, err = util.InitLogger("warn"); err != nil {
			panic(err)
		}
	}
}

// CoreClient hold the client of Jenkins core
type CoreClient struct {
	JenkinsCore
}

// Reload will send the reload request
func (q *CoreClient) Reload() (err error) {
	_, err = q.RequestWithoutData(http.MethodPost, "/reload", map[string]string{
		httpdownloader.ContentType: httpdownloader.ApplicationForm},
		nil, 503)
	return
}

// Restart will send the restart request
func (q *CoreClient) Restart() (err error) {
	_, err = q.RequestWithoutData(http.MethodPost, "/safeRestart", nil, nil, 503)
	return
}

// RestartDirectly restart Jenkins directly
func (q *CoreClient) RestartDirectly() (err error) {
	_, err = q.RequestWithoutData(http.MethodPost, "/restart", nil, nil, 503)
	return
}

// Shutdown puts Jenkins into the quiet mode, wait for existing builds to be completed, and then shut down Jenkins
func (q *CoreClient) Shutdown(safe bool) (err error) {
	if safe {
		_, err = q.RequestWithoutData(http.MethodPost, "/safeExit", nil, nil, 200)
	} else {
		_, err = q.RequestWithoutData(http.MethodPost, "/exit", nil, nil, 200)
	}
	return
}

// PrepareShutdown Put Jenkins in a Quiet mode, in preparation for a restart. In that mode Jenkins don’t start any build
func (q *CoreClient) PrepareShutdown(cancel bool) (err error) {
	if cancel {
		_, err = q.RequestWithoutData(http.MethodPost, "/cancelQuietDown", nil, nil, 200)
	} else {
		_, err = q.RequestWithoutData(http.MethodPost, "/quietDown", nil, nil, 200)
	}
	return
}

// JenkinsIdentity belongs to a Jenkins
type JenkinsIdentity struct {
	Fingerprint   string
	PublicKey     string
	SystemMessage string
}

// GetIdentity returns the identity of a Jenkins
func (q *CoreClient) GetIdentity() (identity JenkinsIdentity, err error) {
	err = q.RequestWithData(http.MethodGet, "/instance", nil, nil, 200, &identity)
	return
}

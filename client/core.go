package client

import (
	"github.com/jenkins-zh/jenkins-cli/util"
	"go.uber.org/zap"
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

// Restart will send the restart request
func (q *CoreClient) Restart() (err error) {
	_, err = q.RequestWithoutData("POST", "/safeRestart", nil, nil, 503)
	return
}

// RestartDirectly restart Jenkins directly
func (q *CoreClient) RestartDirectly() (err error) {
	_, err = q.RequestWithoutData("POST", "/restart", nil, nil, 503)
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
	err = q.RequestWithData("GET", "/instance", nil, nil, 200, &identity)
	return
}

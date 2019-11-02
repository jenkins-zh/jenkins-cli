package client

import "go.uber.org/zap"

var logger *zap.Logger

// SetLogger set a global logger
func SetLogger(zapLogger *zap.Logger) {
	logger = zapLogger
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

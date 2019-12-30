package util

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

// InitLogger returns a logger
func InitLogger(level string) (logger *zap.Logger, err error) {
	rawJSON := []byte(fmt.Sprintf(`{
	  "level": "%s",
	  "encoding": "json",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`, level))
	var cfg zap.Config
	if err = json.Unmarshal(rawJSON, &cfg); err == nil {
		logger, err = cfg.Build()
	}
	return
}

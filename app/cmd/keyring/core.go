package keyring

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/zalando/go-keyring"
)

const (
	// PlaceHolder is the replacer of original credential
	PlaceHolder = "******"
	// KeyTokenPrefix is the prefix of keyring service
	KeyTokenPrefix = "jcli-config-token"
)

// SaveTokenToKeyring store the token to keyring
func SaveTokenToKeyring(config *config.Config) {
	if config == nil {
		return
	}
	for i, item := range config.JenkinsServers {
		token := item.Token
		if token == PlaceHolder {
			continue
		}

		if err := keyring.Set(fmt.Sprintf("%s-%s", KeyTokenPrefix, item.Name), item.UserName, token); err == nil {
			(&item).Token = PlaceHolder
			config.JenkinsServers[i] = item
		}
	}
}

// LoadTokenFromKeyring load token from keyring
func LoadTokenFromKeyring(config *config.Config) {
	for i, item := range config.JenkinsServers {
		if item.Token != PlaceHolder {
			continue
		}
		if token, err := keyring.Get(fmt.Sprintf("%s-%s", KeyTokenPrefix, item.Name), item.UserName); err == nil {
			(&item).Token = token
			config.JenkinsServers[i] = item
		}
	}
}

// DelToken removes the token from keyring
func DelToken(jenkins config.JenkinsServer) (err error) {
	err = keyring.Delete(fmt.Sprintf("%s-%s", KeyTokenPrefix, jenkins.Name), jenkins.UserName)
	return
}

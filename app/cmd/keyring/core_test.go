package keyring_test

import (
	"fmt"
	innerKeyring "github.com/jenkins-zh/jenkins-cli/app/cmd/keyring"
	cfg "github.com/jenkins-zh/jenkins-cli/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
	"testing"
)

func TestSaveTokenToKeyring(t *testing.T) {
	keyring.MockInit()

	// nil case should be handled nicely
	var config *cfg.Config
	innerKeyring.SaveTokenToKeyring(config)

	// empty struct should be handled nicely
	config = &cfg.Config{}
	innerKeyring.SaveTokenToKeyring(config)

	// give it a empty JenkinsServer
	config.JenkinsServers = append(config.JenkinsServers,
		cfg.JenkinsServer{})
	innerKeyring.SaveTokenToKeyring(config)

	// give it a real JenkinsService with token
	config.JenkinsServers = append(config.JenkinsServers,
		cfg.JenkinsServer{
			Token: "I'm a fake token",
		})
	innerKeyring.SaveTokenToKeyring(config)
	assert.Equal(t, innerKeyring.PlaceHolder, config.JenkinsServers[0].Token)
}

func TestLoadTokenFromKeyring(t *testing.T) {
	keyring.MockInit()

	const (
		service  = "fake-service"
		username = "fake-username"
		token    = "fake-token"
	)
	err := keyring.Set(fmt.Sprintf("%s-%s", innerKeyring.KeyTokenPrefix, service), username, token)
	assert.Nil(t, err, "got error when set keyring")

	config := &cfg.Config{
		JenkinsServers: []cfg.JenkinsServer{{
			Name:     service,
			UserName: username,
			Token:    token,
		}},
	}
	innerKeyring.LoadTokenFromKeyring(config)
	assert.Equal(t, token, config.JenkinsServers[0].Token)
}

func TestDelToken(t *testing.T) {
	keyring.MockInit()

	jenkins := cfg.JenkinsServer{}
	// delete a non-exists keyring item
	err := innerKeyring.DelToken(jenkins)
	assert.NotNil(t, err, "got error when delete token from keyring")

	// prepare a keyring item
	const (
		service  = "fake-service"
		username = "fake-username"
		token    = "fake-token"
	)
	err = keyring.Set(fmt.Sprintf("%s-%s", innerKeyring.KeyTokenPrefix, service), username, token)
	assert.Nil(t, err, "got error when set keyring")
	// delete an existing keyring
	jenkins.UserName = username
	jenkins.Token = token
	jenkins.Name = service
	err = innerKeyring.DelToken(jenkins)
	assert.Nil(t, err, "got error when delete keyring")
}

package client

import (
	"encoding/json"
	"fmt"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
)

// CredentialsManager hold the info of credentials client
type CredentialsManager struct {
	JenkinsCore
}

// GetList returns the credential list
func (c *CredentialsManager) GetList(store string) (credentialList CredentialList, err error) {
	api := fmt.Sprintf("/credentials/store/%s/domain/_/api/json?pretty=true&depth=1", store)
	err = c.RequestWithData(http.MethodGet, api, nil, nil, 200, &credentialList)
	return
}

// Delete removes a credential by id from a store
func (c *CredentialsManager) Delete(store, id string) (err error) {
	api := fmt.Sprintf("/credentials/store/%s/domain/_/credential/%s/doDelete", store, id)
	_, err = c.RequestWithoutData(http.MethodPost, api, nil, nil, 200)
	return
}

// Create create a credential in Jenkins
func (c *CredentialsManager) Create(store, credential string) (err error) {
	api := fmt.Sprintf("/credentials/store/%s/domain/_/createCredentials", store)

	logger.Debug("create credential", zap.String("api", api), zap.String("payload", credential))

	formData := url.Values{}
	formData.Add("json", fmt.Sprintf(`{"credentials": %s}`, credential))
	payload := strings.NewReader(formData.Encode())

	_, err = c.RequestWithoutData(http.MethodPost, api,
		map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 200)
	return
}

// CreateUsernamePassword create username and password credential in Jenkins
func (c *CredentialsManager) CreateUsernamePassword(store string, cred UsernamePasswordCredential) (err error) {
	var payload []byte
	cred.Class = "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
	if payload, err = json.Marshal(cred); err == nil {
		err = c.Create(store, string(payload))
	}
	return
}

// CreateSecret create token credential in Jenkins
func (c *CredentialsManager) CreateSecret(store string, cred StringCredentials) (err error) {
	var payload []byte
	cred.Class = "org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl"
	if payload, err = json.Marshal(cred); err == nil {
		err = c.Create(store, string(payload))
	}
	return
}

// Credential of Jenkins
type Credential struct {
	Description string `json:"description"`
	DisplayName string
	Fingerprint interface{}
	FullName    string
	ID          string `json:"id"`
	TypeName    string
	Class       string `json:"$class"`
	Scope       string `json:"scope"`
}

// UsernamePasswordCredential hold the username and password
type UsernamePasswordCredential struct {
	Credential `json:",inline"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

// StringCredentials hold a token
type StringCredentials struct {
	Credential `json:",inline"`
	Secret     string `json:"secret"`
}

// CredentialList contains many credentials
type CredentialList struct {
	Description     string
	DisplayName     string
	FullDisplayName string
	FullName        string
	Global          bool
	URLName         string
	Credentials     []Credential
}

package client

import "fmt"

// CredentialsManager hold the info of credentials client
type CredentialsManager struct {
	JenkinsCore
}

// GetList returns the credential list
func (c *CredentialsManager) GetList(store string) (credentialList CredentialList, err error) {
	api := fmt.Sprintf("/credentials/store/%s/domain/_/api/json?pretty=true&depth=1", store)
	err = c.RequestWithData("GET", api, nil, nil, 200, &credentialList)
	return
}

// Delete removes a credential by id from a store
func (c *CredentialsManager) Delete(store, id string) (err error) {
	api := fmt.Sprintf("/credentials/store/%s/domain/_/credential/%s/doDelete", store, id)
	_, err = c.RequestWithoutData("POST", api, nil, nil, 200)
	return
}

// Credential of Jenkins
type Credential struct {
	Description string
	DisplayName string
	Fingerprint string
	FullName    string
	ID          string
	TypeName    string
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

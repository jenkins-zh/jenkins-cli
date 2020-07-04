package client

import "net/http"

// CASCManager is the client of configuration as code
type CASCManager struct {
	JenkinsCore
}

// Export exports the config of configuration-as-code
func (c *CASCManager) Export() (config string, err error) {
	var (
		data       []byte
		statusCode int
	)

	if statusCode, data, err = c.Request(http.MethodPost, "/configuration-as-code/export",
		nil, nil); err == nil &&
		statusCode != 200 {
		err = c.ErrorHandle(statusCode, data)
	}
	config = string(data)
	return
}

// Schema get the schema of configuration-as-code
func (c *CASCManager) Schema() (schema string, err error) {
	var (
		data       []byte
		statusCode int
	)

	if statusCode, data, err = c.Request(http.MethodPost, "/configuration-as-code/schema",
		nil, nil); err == nil &&
		statusCode != 200 {
		err = c.ErrorHandle(statusCode, data)
	}
	schema = string(data)
	return
}

// Reload reload the config of configuration-as-code
func (c *CASCManager) Reload() (err error) {
	_, err = c.RequestWithoutData(http.MethodPost, "/configuration-as-code/reload",
		nil, nil, 200)
	return
}

// Apply apply the config of configuration-as-code
func (c *CASCManager) Apply() (err error) {
	_, err = c.RequestWithoutData(http.MethodPost, "/configuration-as-code/apply",
		nil, nil, 200)
	return
}

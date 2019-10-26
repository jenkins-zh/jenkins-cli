package client

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"net/http"
)

//PrepareRestart only for test
func PrepareRestart(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, statusCode int)  {
	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/safeRestart", rootURL), nil)
	response := PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	response.StatusCode = statusCode
	return
}

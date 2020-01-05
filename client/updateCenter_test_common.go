package client

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"github.com/jenkins-zh/jenkins-cli/util"
	"net/http"
	"net/url"
	"strings"
)

// PrepareForSetMirrorCertificate only for test
func PrepareForSetMirrorCertificate(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, enable bool) {
	api := "/update-center-mirror/use"
	if !enable {
		api = "/update-center-mirror/remove"
	}

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", rootURL, api), nil)
	request.Header.Add(util.ContentType, util.ApplicationForm)
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForChangeUpdateCenterSite only for test
func PrepareForChangeUpdateCenterSite(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, name, updateCenterURL string) {
	formData := url.Values{}
	formData.Add("site", updateCenterURL)
	payload := strings.NewReader(formData.Encode())

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/pluginManager/siteConfigure", rootURL), payload)
	request.Header.Add(util.ContentType, util.ApplicationForm)
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

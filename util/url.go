package util

import (
	"net/url"
	"path"
	"strings"
)

// URLJoin is a util function to join host URL and API URL
func URLJoin(host, api string) (targetURL *url.URL, err error) {
	if targetURL, err = url.Parse(host); err == nil {
		pathURL, _ := url.Parse(path.Join(targetURL.Path, api))
		targetURL = targetURL.ResolveReference(pathURL)
	}
	return
}

//URLJoinAsString  is a util function to join host URL and API URL
func URLJoinAsString(host, api string) (targetURLStr string, err error) {
	var targetURL *url.URL
	if targetURL, err = URLJoin(host, api); err == nil {
		targetURLStr = targetURL.String()
	}
	return
}
//ArraySplitAndDeleteEmpty split string and delete empty element
func ArraySplitAndDeleteEmpty(s, sep string) []string {
	var r []string
	stringList := strings.Split(s, sep)
	for _, str := range stringList {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

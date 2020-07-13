package cmd

import "net/http"

type SelfUpgradeOption struct {
	ShowProgress bool
	RoundTripper http.RoundTripper
}

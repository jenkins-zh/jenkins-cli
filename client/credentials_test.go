package client_test

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("job test", func() {
	var (
		ctrl               *gomock.Controller
		credentialsManager CredentialsManager
		roundTripper       *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		credentialsManager = CredentialsManager{}
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		credentialsManager.RoundTripper = roundTripper
		credentialsManager.URL = "http://localhost"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("", func() {

	})
})

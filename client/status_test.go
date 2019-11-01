package client

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("status test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		statusClient JenkinsStatusClient

		username string
		password string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		statusClient = JenkinsStatusClient{}
		statusClient.RoundTripper = roundTripper
		statusClient.URL = "http://localhost"

		username = "admin"
		password = "token"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Get", func() {
		It("should success", func() {
			statusClient.UserName = username
			statusClient.Token = password

			PrepareGetStatus(roundTripper, statusClient.URL, username, password)

			status, err := statusClient.Get()
			Expect(err).To(BeNil())
			Expect(status).NotTo(BeNil())
			Expect(status.NodeName).To(Equal("master"))
			Expect(status.Version).To(Equal("version"))
		})
	})
})

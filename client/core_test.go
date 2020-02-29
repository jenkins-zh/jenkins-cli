package client

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("core test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		coreClient   CoreClient

		username string
		password string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		coreClient = CoreClient{}
		coreClient.RoundTripper = roundTripper
		coreClient.URL = "http://localhost"

		username = "admin"
		password = "token"

		coreClient.UserName = username
		coreClient.Token = password
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Get", func() {
		It("should success", func() {
			PrepareRestart(roundTripper, coreClient.URL, username, password, 503)

			err := coreClient.Restart()
			Expect(err).To(BeNil())
		})

		It("should error, 400", func() {
			PrepareRestart(roundTripper, coreClient.URL, username, password, 400)

			err := coreClient.Restart()
			Expect(err).To(HaveOccurred())
		})

		It("should success", func() {
			PrepareRestartDirectly(roundTripper, coreClient.URL, username, password, 503)

			err := coreClient.RestartDirectly()
			Expect(err).To(BeNil())
		})

		It("GetIdentity", func() {
			PrepareForGetIdentity(roundTripper, coreClient.URL, username, password)

			identity, err := coreClient.GetIdentity()
			Expect(err).NotTo(HaveOccurred())
			Expect(identity).To(Equal(JenkinsIdentity{
				Fingerprint:   "fingerprint",
				PublicKey:     "publicKey",
				SystemMessage: "systemMessage",
			}))
		})
	})
})
